package worker

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/node/master"
)

func Load() {

	config.CfgMaster.IsMaster = true
	config.CfgMaster.Port = 6660
	config.CfgMaster.LBAlgorithm = "OP"

	config.CfgWorker.IsWorker = true
	config.CfgWorker.MasterIP = "localhost"
	config.CfgWorker.MasterPort = 6660
	config.CfgWorker.HeartbeatInterval = 2

	config.CfgWorker.LoadReport.LoadReportInterval = 5
	config.CfgWorker.LoadReport.MaxNetworkBandwidth = 200
	config.CfgWorker.LoadReport.NetworkAdapterName = "en0"

	config.CfgWorker.ServicePort = -1

	//logger.SetLogLevel(logger.Debug)

}

func RunWorker() *WorkerClient {

	c := NewWorkerClient()
	c.Run()

	return c

}

func TestWorkerNumber(t *testing.T) {

	Load()

	config.CfgMaster.Port = 6660
	config.CfgWorker.MasterPort = 6660

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

	logger.SetLogLevel(logger.Debug)
	config.CfgMaster.Port = 6661
	config.CfgWorker.MasterPort = 6661

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
		logger.LogInfof("The worker node number is correct.\n")
	}

	<-ch

	time.Sleep(time.Second * time.Duration(config.CfgWorker.HeartbeatInterval*4))
	num := node.NodeContainer.GetNodeListLen()
	if !reflect.DeepEqual(num, 3) {
		t.Error("The worker node number is not correct!!!")
		logger.LogErrf("The worker node number is %d\n", num)
	}

	wg.Wait()

}

func TestWorkerRetry(t *testing.T) {

	Load()
	config.CfgMaster.Port = 6662
	config.CfgWorker.MasterPort = 6662

	var wc *WorkerClient
	go func() {
		wc = RunWorker()
	}()

	time.Sleep(time.Second * time.Duration(10))

	ticker := time.NewTicker(time.Duration(20) * time.Second)
	s := master.NewMasterServer()

	go func() {
		s.Run()
	}()

	select {
	case <-ticker.C:
		wc.Stop()
		s.Stop()
	}

}
