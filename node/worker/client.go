package worker

import (
	"log"

	"google.golang.org/grpc"
)

//LBWorkerClient Load Balance Worker Client
type LBWorkerClient struct {
	conn *grpc.ClientConn
}

func newLBWorkerClient(addr *string) *LBWorkerClient {

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	worker := &LBWorkerClient{conn: conn}
	return worker
	// stream, err := c.SendHeartbeat(context.Background())
	// stream.Send(&pb.HeartbeatRequest{CpuUsageRate: 0.5, MemoryUsageRate: 0.5, BandwidthUsageRate: 0.5})

}

// func watch(worker *WorkerClient ){
// 	watcher  := NewWatcher(worker)
// 	//watcher.Watch()
// }
