package examples

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/harlanc/gobalan/balancer"
	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/node/master"
	"github.com/harlanc/gobalan/node/worker"
	"github.com/harlanc/gobalan/proto"
	"github.com/harlanc/gobalan/utils"
	"google.golang.org/grpc"
)

func Load() {

	config.SetCfgPath("/Users/zexu/go/src/github.com/harlanc/gobalan/config/config.ini")
	config.LoadCfg()

	config.CfgMaster.IsMaster = true
	config.CfgMaster.Port = 6666
	config.CfgMaster.LBAlgorithm = "OP"

	config.CfgWorker.IsWorker = true
	config.CfgWorker.MasterIP = "192.168.0.104"
	config.CfgWorker.MasterPort = 6666
	config.CfgWorker.HeartbeatInterval = 2

	config.CfgWorker.LoadReport.LoadReportInterval = 5
	config.CfgWorker.LoadReport.MaxNetworkBandwidth = 200
	config.CfgWorker.LoadReport.NetworkAdapterName = "en0"

	config.CfgWorker.ServicePort = -1

	balancer.LoadBalanceType()

	logger.SetLogLevel(logger.Debug)

}

func RunWorker() *worker.WorkerClient {

	c := worker.NewWorkerClient()
	c.Run()

	return c

}

func newRPCPickClient() *balancer.RPCPickClient {

	dialAddr := fmt.Sprintf("%s:%s", config.CfgWorker.MasterIP, utils.Int322String(config.CfgWorker.MasterPort))
	conn, err := grpc.Dial(dialAddr, grpc.WithInsecure())

	if err != nil {
		logger.LogErrf("faild to connect: %v", err)
	}

	return balancer.NewRPCPickClient(conn)

}

func TestBalancerRoundRobin(t *testing.T) {

	Load()

	balancer.CurrentBalanceType = proto.BalanceType_RoundRobin

	ticker := time.NewTicker(time.Duration(10) * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(2)

	s := master.NewMasterServer()
	go func() {
		s.Run()
		wg.Done()
	}()

	c1 := RunWorker()
	c2 := RunWorker()
	c3 := RunWorker()
	c4 := RunWorker()

	go func() {
		select {
		case <-ticker.C:
			s.Stop()
			c1.Stop()
			c2.Stop()
			c3.Stop()
			c4.Stop()
			wg.Done()
		}

	}()

	time.Sleep(time.Second * time.Duration(4))

	client := newRPCPickClient()
	n1, _ := client.RPCPick()
	nodelist := node.NodeContainer.GetNodeList()
	l := node.NodeContainer.GetNodeListLen()
	var idx int
	var v *node.Node

	for idx, v = range nodelist {
		if n1.WorkerId == v.WorkerID {
			break
		}
	}

	n1, _ = client.RPCPick()
	if !reflect.DeepEqual(n1.WorkerId, nodelist[(idx+1)%l].WorkerID) {
		t.Error("The roundrobin is not correct")

	}

	n1, _ = client.RPCPick()
	if !reflect.DeepEqual(n1.WorkerId, nodelist[(idx+2)%l].WorkerID) {
		t.Error("The roundrobin is not correct")

	}
	n1, _ = client.RPCPick()
	if !reflect.DeepEqual(n1.WorkerId, nodelist[(idx+3)%l].WorkerID) {
		t.Error("The roundrobin is not correct")

	}
	n1, _ = client.RPCPick()
	if !reflect.DeepEqual(n1.WorkerId, nodelist[(idx+4)%l].WorkerID) {
		t.Error("The roundrobin is not correct")

	}
	n1, _ = client.RPCPick()
	if !reflect.DeepEqual(n1.WorkerId, nodelist[(idx+5)%l].WorkerID) {
		t.Error("The roundrobin is not correct")

	}

	wg.Wait()

}
