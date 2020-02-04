package worker

import (
	"context"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	pb "github.com/harlanc/gobalan/proto"
)

// WatchClient implements the WatchClient interface
type WatchClient struct {
	ctx      context.Context
	client   pb.WatchClient
	callOpts []grpc.CallOption

	stream      *clientWatchStream
	BalanceType pb.BalanceType
	workerID    uint32
}

// clientWatchStream tracks all watch resources attached to a single grpc stream.
type clientWatchStream struct {
	ctx              context.Context
	owner            WatchClient
	gRPCClientStream pb.Watch_WatchClient
	cancel           context.CancelFunc
	wgWorkerID       sync.WaitGroup
	heartbeatTicker  *time.Ticker
}

// watchStreamRequest is a union of the supported watch request operation types
type watchStreamRequest interface {
	toWatchCreateRequestPB() *pb.WatchRequest
	toWatchHeartbeatRequestPB() *pb.WatchRequest
	toWatchLoadReportRequestPB() *pb.WatchRequest
}

// watchRequest is issued by the subscriber to start a new watcher
type watchRequest struct {
	workerID          uint32
	servicePort       string
	heartbeatInterval uint32
	loadReportData    *any.Any
}

//NewWatcher NewWatcher
func NewWatcher(c *grpc.ClientConn) *WatchClient {
	return &WatchClient{client: pb.NewWatchClient(c), ctx: context.Background()}
}

func (w *WatchClient) newClientWatchStream(inctx context.Context) *clientWatchStream {

	ctx, cancel := context.WithCancel(inctx)
	var wc pb.Watch_WatchClient
	var err error
	//https://github.com/grpc/grpc-go/blob/master/examples/features/cancellation/client/main.go
	if wc, err = w.client.Watch(w.ctx, w.callOpts...); err != nil {
		cancel()
		return nil
	}
	wws := &clientWatchStream{

		gRPCClientStream: wc,
		ctx:              ctx,
		cancel:           cancel,

		heartbeatTicker: time.NewTicker(time.Duration(config.CfgWorker.HeartbeatInterval) * time.Second),
	}

	return wws
}

//Run run the client
func (w *WatchClient) Run() {

	cws := w.newClientWatchStream(w.ctx)
	w.stream = cws

	go func() {
		cws.sendLoop()
	}()

	go func() {
		cws.receiveLoop()
	}()
}

func (w *clientWatchStream) sendLoop() {

	wr := &watchRequest{}

	//1.first send a create request to get worker id
	w.wgWorkerID.Add(1)
	err := w.gRPCClientStream.Send(wr.toWatchCreateRequestPB())
	if err != nil {
		logger.LogErr(err)
		w.cancel()
		return
	}
	w.wgWorkerID.Wait()
	wr.workerID = w.owner.workerID

	go func() {
		for {
			select {

			case <-w.heartbeatTicker.C:

				err := w.gRPCClientStream.Send(wr.toWatchHeartbeatRequestPB())
				logger.LogDebug("Send Heartbeat")
				if err != nil {
					w.cancel()
				}

			case <-w.ctx.Done():
				return
			}
		}

	}()

	if w.owner.BalanceType == pb.BalanceType_OptimalPerformance {

		go func() {
			m := NewMonitor(w.ctx)
			m.Start()
			select {
			case s := <-m.StatPB:
				any, err := ptypes.MarshalAny(&s)
				if err != nil {
					logger.LogErr(err)
					w.cancel()
				}
				wr.loadReportData = any
				err = w.gRPCClientStream.Send(wr.toWatchLoadReportRequestPB())
				logger.LogDebug("Send Load Report")
				if err != nil {
					logger.LogErr(err)
					w.cancel()
				}
			case <-w.ctx.Done():
				return
			}

		}()
	}

}
func (w *clientWatchStream) receiveLoop() {

	for {

		resp, err := w.gRPCClientStream.Recv()

		if err != nil {
			if err != io.EOF {
				w.cancel()
			}
			return
		}
		w.owner.workerID = resp.WorkerId
		w.owner.BalanceType = resp.BalanceType
		w.wgWorkerID.Done()

	}

}

func (wr *watchRequest) toWatchCreateRequestPB() *pb.WatchRequest {
	req := &pb.WatchCreateRequest{
		ServicePort:       config.CfgWorker.ServicePort,
		HeartbeatInterval: config.CfgWorker.HeartbeatInterval,
	}
	cr := &pb.WatchRequest_CreateRequest{CreateRequest: req}
	return &pb.WatchRequest{RequestUnion: cr}
}

func (wr *watchRequest) toWatchHeartbeatRequestPB() *pb.WatchRequest {
	req := &pb.WatchHeartbeatRequest{
		WorkerId: wr.workerID,
	}
	cr := &pb.WatchRequest_HeartbeatRequest{HeartbeatRequest: req}
	return &pb.WatchRequest{RequestUnion: cr}
}

func (wr *watchRequest) toWatchLoadReportRequestPB() *pb.WatchRequest {
	req := &pb.WatchLoadReportRequest{
		WorkerId:       wr.workerID,
		LoadReportData: wr.loadReportData,
	}
	cr := &pb.WatchRequest_LoadReportRequest{LoadReportRequest: req}
	return &pb.WatchRequest{RequestUnion: cr}
}
