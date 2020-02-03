package master

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/harlanc/gobalan/balancer"
	pb "github.com/harlanc/gobalan/proto"
	"github.com/harlanc/gobalan/utils"

	"github.com/harlanc/gobalan/node"
)

//

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
	mu               sync.RWMutex
	cancel           context.CancelFunc
	wg               sync.WaitGroup
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
		sws.sendAndClosee()
		sws.wg.Done()
	}()

	errc := make(chan error, 1)

	go func() {
		if rerr := sws.recvLoop(); rerr != nil {
			errc <- rerr
		}
	}()

	sws.wg.Wait()
	return err
}

func (sws *serverWatchStream) recvLoop() error {

	for {

		req, err := sws.gRPCServerStream.Recv()
		if err == io.EOF {

			return nil
		}
		if err != nil {
			return err
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

			wr := &pb.WatchResponse{
				WorkerId:    utils.GenerateWorkerID(),
				BalanceType: balancer.CurrentBalanceType,
				Created:     true,
				Canceled:    err != nil,
			}
			if err != nil {
				fmt.Println(err)
			}
			sws.resp <- wr

		case *pb.WatchRequest_HeartbeatRequest:

		case *pb.WatchRequest_LoadReportRequest:

			data := uv.LoadReportRequest.LoadReportData
			stat := &pb.Stat{}
			if err := ptypes.UnmarshalAny(data, stat); err != nil {
				fmt.Println("The load report data cannot be unmarshaled ")
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
			return sws.ctx.Err()
		default:

		}
	}
}

func (sws *serverWatchStream) sendAndClosee() {

	select {
	case c, ok := <-sws.resp:
		if !ok {
			sws.cancel()
			return
		}
		if err := sws.gRPCServerStream.SendAndClose(c); err != nil {
			fmt.Println(err)
			sws.cancel()
			return
		}
	case <-sws.ctx.Done():
		return

	}

}
