package worker

import (
	"fmt"
	"log"

	"github.com/harlanc/gobalan/config"
	"google.golang.org/grpc"
)

//WorkerClient Load Balance Worker Client
type WorkerClient struct {
	watchClient *WatchClient
}

//NewWorkerClient new a worker client
func NewWorkerClient() *WorkerClient {

	dialAddr := fmt.Sprintf("%s:%s", config.CfgWorker.MasterIP, config.CfgWorker.MasterPort)
	conn, err := grpc.Dial(dialAddr, grpc.WithInsecure())
	//defer conn.Close()

	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}

	worker := &WorkerClient{watchClient: NewWatcher(conn)}

	return worker

}

//Run the worker
func (wc *WorkerClient) Run() {
	wc.watchClient.Run()
}
