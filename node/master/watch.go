package master

import (
	"context"
	"fmt"
	"io"
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

	cancel context.CancelFunc
	wg     sync.WaitGroup
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
		go func(cancel context.CancelFunc) {

			select {
			case <-t.C:
				fmt.Printf("The recv stream is timeout in %d seconds.", streamRecvTimeout)
				cancel()
			case <-c:
			}

		}(sws.cancel)

		req, err := sws.gRPCServerStream.Recv()

		c <- struct{}{}
		t.Reset(time.Duration(streamRecvTimeout) * time.Second)

		if err == io.EOF {
			return
		}
		if err != nil {
			sws.cancel()
			return
		}

		switch uv := req.RequestUnion.(type) {

		case *pb.WatchRequest_CreateRequest:
			if uv.CreateRequest == nil {
				break
			}
			ip, err := utils.GetContextIP(sws.gRPCServerStream.Context())
			nd := node.Node{IP: ip, Port: uv.CreateRequest.ServicePort}
			workerid := utils.GenerateWorkerID()
			node.NodeContainer.InsertNode(workerid, &nd)

			streamRecvTimeout = uv.CreateRequest.HeartbeatInterval * 3

			wr := &pb.WatchResponse{
				WorkerId:    utils.GenerateWorkerID(),
				BalanceType: balancer.CurrentBalanceType,
				Created:     true,
				Canceled:    err != nil,
			}
			if err != nil {
				logger.LogErr(err)
			}
			sws.resp <- wr

		case *pb.WatchRequest_HeartbeatRequest:
			logger.LogDebug("receive Heartbeat")

		case *pb.WatchRequest_LoadReportRequest:

			data := uv.LoadReportRequest.LoadReportData
			stat := &pb.Stat{}
			if err := ptypes.UnmarshalAny(data, stat); err != nil {
				logger.LogErr("The load report data cannot be unmarshaled ")
				continue
			}

			workerid := uv.LoadReportRequest.WorkerId
			node.NodeContainer.UpdateNode(workerid, stat)

		case *pb.WatchRequest_CancelRequest:
			if uv.CancelRequest != nil {

				if err == nil {
					sws.resp <- &pb.WatchResponse{
						Canceled: true,
					}
				}
			}
		}
		select {
		case <-sws.ctx.Done():
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
				sws.cancel()
				return
			}
			if err := sws.gRPCServerStream.Send(c); err != nil {
				logger.LogErr(err)
				sws.cancel()
				return
			}
		case <-sws.ctx.Done():
			return

		}

	}

}
