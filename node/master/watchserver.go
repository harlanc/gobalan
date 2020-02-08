package master

import (
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/harlanc/gobalan/balancer"
	"github.com/harlanc/gobalan/logger"
	pb "github.com/harlanc/gobalan/proto"
	"github.com/harlanc/gobalan/utils"

	"github.com/harlanc/gobalan/node"
)

var (
	streamRecvTimeout uint32 = 5 * 60
)

type watchServer struct {
}

// NewWatchServer returns a new watch server.
func NewWatchServer() pb.WatchServer {
	return &watchServer{}
}

type serverWatchStream struct {
	ctx              context.Context
	owner            *watchServer
	gRPCServerStream pb.Watch_WatchServer
	resp             chan *pb.WatchResponse
	workerID         uint32
	cancel           context.CancelFunc

	wg sync.WaitGroup
}

var ctrlStreamBufLen uint = 16

func (ws *watchServer) Watch(stream pb.Watch_WatchServer) (err error) {

	ctx, cancel := context.WithCancel(stream.Context())
	sws := serverWatchStream{
		ctx:              ctx,
		gRPCServerStream: stream,
		resp:             make(chan *pb.WatchResponse, ctrlStreamBufLen),
		cancel:           cancel,
		owner:            ws,
	}

	sws.wg.Add(1)

	go func() {
		sws.sendLoop()
		sws.wg.Done()
	}()
	go func() {
		sws.recvLoop()
	}()

	sws.wg.Wait()
	return err
}

func (sws *serverWatchStream) recvLoop() {

	t := time.NewTimer(time.Duration(streamRecvTimeout) * time.Second)
	c := make(chan struct{}, 1)

	for {
		go func() {

			select {
			case <-t.C:
				logger.LogWarnf("The recv stream is timeout in %d seconds.\n", streamRecvTimeout)
				sws.close()
				return
			case <-c:

			case <-sws.ctx.Done():

				return
			}

		}()

		req, err := sws.gRPCServerStream.Recv()

		c <- struct{}{}

		t.Reset(time.Duration(streamRecvTimeout) * time.Second)

		if err != nil {
			logger.LogErr(err)
			sws.close()
			return
		}

		switch uv := req.RequestUnion.(type) {

		case *pb.WatchRequest_CreateRequest:
			if uv.CreateRequest == nil {
				break
			}
			ip, err := utils.GetContextIP(sws.gRPCServerStream.Context())
			logger.LogDebugf("The client IP is %s\n", ip)

			workerid := utils.GenerateWorkerID()
			nd := node.Node{IP: ip, Port: uv.CreateRequest.ServicePort, WorkerID: workerid, ServiceStatus: pb.ServiceStatus_Up}
			node.NodeContainer.InsertNode(workerid, &nd)
			sws.workerID = workerid

			streamRecvTimeout = uv.CreateRequest.HeartbeatInterval * 2
			t.Reset(time.Duration(streamRecvTimeout) * time.Second)

			wr := &pb.WatchResponse{
				WorkerId:    workerid,
				BalanceType: balancer.CurrentBalanceType,
			}
			if err != nil {
				logger.LogErr(err)
			}
			sws.resp <- wr

		case *pb.WatchRequest_HeartbeatRequest:

			ss := uv.HeartbeatRequest.ServiceStatus
			node.NodeContainer.UpdateNodeServiceStatus(sws.workerID, ss)
			logger.LogDebugf("The Woker %d receive Heartbeat\n", sws.workerID)

		case *pb.WatchRequest_LoadReportRequest:

			data := uv.LoadReportRequest.LoadReportData
			stat := &pb.Stat{}
			if err := ptypes.UnmarshalAny(data, stat); err != nil {
				logger.LogErr("The load report data cannot be unmarshaled ")
				continue
			}
			node.NodeContainer.UpdateNodeStat(sws.workerID, stat)
		}
		select {
		case <-sws.ctx.Done():
			logger.LogDebug("recvLoop ctx is done")
			return
		default:

		}
	}
}

func (sws *serverWatchStream) sendLoop() {

	for {
		select {
		case c, ok := <-sws.resp:
			if !ok {
				sws.close()
				return
			}
			if err := sws.gRPCServerStream.Send(c); err != nil {
				logger.LogErr(err)
				sws.close()
				return
			}
		case <-sws.ctx.Done():
			logger.LogDebug("sendLoop ctx is done")
			return

		}

	}

}

func (sws *serverWatchStream) close() {

	logger.LogErr("serverWatchStream is closed")

	node.NodeContainer.DeleteNode(sws.workerID)
	// sws.closec <- struct{}{}
	sws.cancel()

}
