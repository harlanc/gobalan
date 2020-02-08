package config

import (
	"os"
	"path"
	"runtime"

	"github.com/harlanc/gobalan/logger"
	"gopkg.in/ini.v1"
)

var (
	//CfgMaster exported Master configurations
	CfgMaster *cfgMaster
	//CfgWorker exported Worker configurations
	CfgWorker *cfgWorker
)

func init() {
	CfgMaster = new(cfgMaster)
	CfgWorker = new(cfgWorker)
	cfgLoadReport := new(loadReport)
	CfgWorker.LoadReport = cfgLoadReport

	LoadCfg()
}

//cfgMaster configurations
type cfgMaster struct {
	IsMaster    bool
	Port        int32
	LBAlgorithm string
}

//Worker configurations
type cfgWorker struct {
	IsWorker          bool
	MasterIP          string
	MasterPort        int32
	ServicePort       int32
	HeartbeatInterval uint32
	LoadReport        *loadReport
}

//LoadReport of worker
type loadReport struct {
	LoadReportInterval  uint32
	MaxNetworkBandwidth float32
	NetworkAdapterName  string
}

//LoadCfg load the configuration file config.ini into revelant structs.
func LoadCfg() {

	_, filename, _, ok := runtime.Caller(1)
	var cfgpath string
	if ok {
		cfgpath = path.Join(path.Dir(filename), "config.ini")
	} else {
		logger.LogErrf("Cannot open configuration file :%s\n", cfgpath)
		os.Exit(3)
	}

	cfg, err := ini.Load(cfgpath)
	if err != nil {
		logger.LogErr(err)
		os.Exit(3)
	}

	err = cfg.Section("Master").MapTo(CfgMaster)
	if err != nil {
		logger.LogErr("The section Master's data cannot be load")
		os.Exit(3)
	}

	err = cfg.Section("Worker").MapTo(CfgWorker)
	if err != nil {
		logger.LogErr("The section Worker's data cannot be load")
		os.Exit(3)
	}

	err = cfg.ChildSections("Worker")[0].MapTo(CfgWorker.LoadReport)
	if err != nil {
		logger.LogErr("The section Worker.LoadReport's data cannot be load")
		os.Exit(3)
	}
}
