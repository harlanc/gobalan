package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/harlanc/gobalan/balancer"
	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
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
	config.CfgWorker.MasterPort = "6666"
	config.CfgWorker.HeartbeatInterval = 200

	config.CfgWorker.LoadReport.LoadReportInterval = 1
	config.CfgWorker.LoadReport.MaxNetworkBandwidth = 200
	config.CfgWorker.LoadReport.NetworkAdapterName = "en0"

	balancer.SetBalanceType()

	logger.SetLogLevel(logger.Debug)

}

func RunWorker() *WorkerClient {

	c := NewWorkerClient()
	c.Run()

	return c

}

func TestWorker(t *testing.T) {

	Load()

	ticker := time.NewTicker(time.Duration(60) * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)

	s := master.NewMasterServer()

	go func() {
		s.Run()
	}()

	c1 := RunWorker()
	//c2 := RunWorker()

	go func() {

		select {
		case <-ticker.C:
			c1.Stop()
			//c2.Stop()
			s.Stop()
			wg.Done()
		}

	}()

	wg.Wait()

}
