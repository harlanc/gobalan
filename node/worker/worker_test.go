package worker

import (
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/harlanc/gobalan/balancer"

	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/node/master"
)

func Load() {

	config.SetCfgPath("/Users/zexu/go/src/github.com/harlanc/gobalan/config/config.ini")
	config.LoadCfg()

	config.CfgMaster.IsMaster = true
	config.CfgMaster.Port = "6666"
	config.CfgMaster.LBAlgorithm = "OP"

	config.CfgWorker.IsWorker = true
	config.CfgWorker.MasterIP = "192.168.0.104"
	config.CfgWorker.MasterPort = 6666
	config.CfgWorker.HeartbeatInterval = 2

	config.CfgWorker.LoadReport.LoadReportInterval = 5
	config.CfgWorker.LoadReport.MaxNetworkBandwidth = 200
	config.CfgWorker.LoadReport.NetworkAdapterName = "en0"

	config.CfgWorker.ServicePort = -1

	balancer.SetBalanceType()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger.SetLogLevel(logger.Debug)

}

func RunWorker() *WorkerClient {

	c := NewWorkerClient()
	c.Run()

	return c

}

func TestWorkerNumber(t *testing.T) {

	Load()

	ticker := time.NewTicker(time.Duration(10) * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)

	s := master.NewMasterServer()

	go func() {
		s.Run()
	}()

	c1 := RunWorker()
	c2 := RunWorker()

	go func() {

		select {
		case <-ticker.C:
			c1.Stop()
			c2.Stop()
			s.Stop()
			wg.Done()
		}

	}()

	time.Sleep(time.Second * time.Duration(config.CfgWorker.HeartbeatInterval))

	if !reflect.DeepEqual(node.NodeContainer.GetNodeListLen(), 2) {
		t.Error("The worker node number is not correct!!!")

	}

	wg.Wait()

}

func TestWorkerTimeout(t *testing.T) {

	Load()

	serverticker := time.NewTicker(time.Duration(20) * time.Second)
	clientticker := time.NewTicker(time.Duration(10) * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(2)

	s := master.NewMasterServer()

	go func() {
		s.Run()
	}()

	c1 := RunWorker()
	c2 := RunWorker()
	c3 := RunWorker()
	c4 := RunWorker()

	ch := make(chan int)

	go func() {

		select {
		case <-clientticker.C:
			c1.Stop()
			wg.Done()
			ch <- 1
		}
	}()

	go func() {
		select {
		case <-serverticker.C:
			s.Stop()
			c2.Stop()
			c3.Stop()
			c4.Stop()
			wg.Done()
		}

	}()

	time.Sleep(time.Second * time.Duration(config.CfgWorker.HeartbeatInterval))

	if !reflect.DeepEqual(node.NodeContainer.GetNodeListLen(), 4) {
		t.Error("The worker node number is not correct!!!")

	} else {
		logger.LogInfof("The worker node number is correct.")
	}

	<-ch

	time.Sleep(time.Second * time.Duration(config.CfgWorker.HeartbeatInterval*2))
	num := node.NodeContainer.GetNodeListLen()
	if !reflect.DeepEqual(num, 3) {
		t.Error("The worker node number is not correct!!!")
		logger.LogErrf("The worker node number is %d\n", num)
	}

	wg.Wait()

}
