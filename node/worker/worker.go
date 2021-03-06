package worker

import (
	"fmt"

	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	"google.golang.org/grpc"
)

//WorkerClient Load Balance Worker Client
type WorkerClient struct {
	watchClient *WatchClient
}

//NewWorkerClient new a worker client
func NewWorkerClient() *WorkerClient {

	dialAddr := fmt.Sprintf("%s:%d", config.CfgWorker.MasterIP, config.CfgWorker.MasterPort)
	conn, err := grpc.Dial(dialAddr, grpc.WithInsecure())
	//defer conn.Close()

	if err != nil {
		logger.LogErrf("faild to connect: %v", err)
	}

	worker := &WorkerClient{watchClient: NewWatcherClient(conn)}
	return worker

}

//Run the worker
func (wc *WorkerClient) Run() {
	wc.watchClient.run()
}

//Stop the worker
func (wc *WorkerClient) Stop() {

	logger.LogInfof("Work client %d is stoped.\n", wc.watchClient.workerID)
	wc.watchClient.stop()

}
