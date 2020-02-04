package main

import (
	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/node/master"
	"github.com/harlanc/gobalan/node/worker"
)

func main() {

	config.SetCfgPath("/Users/zexu/go/src/github.com/harlanc/gobalan/config/config.ini")
	config.LoadCfg()

	if config.CfgMaster.IsMaster {

		s := master.NewMasterServer()
		s.Run()
	}

	if config.CfgWorker.IsWorker {

		c := worker.NewWorkerClient()
		c.Run()
	}

}
