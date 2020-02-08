package config

import (
	"github.com/harlanc/gobalan/logger"
	"gopkg.in/ini.v1"
)

var (
	//CfgMaster exported Master configurations
	CfgMaster *cfgMaster
	//CfgWorker exported Worker configurations
	CfgWorker *cfgWorker
	//cfgPath the configuration file full path
	cfgPath *string
)

func init() {

	CfgMaster = new(cfgMaster)

	CfgWorker = new(cfgWorker)
	cfgLoadReport := new(loadReport)
	CfgWorker.LoadReport = cfgLoadReport

	cfgPath = new(string)

}

//cfgMaster configurations
type cfgMaster struct {
	IsMaster    bool
	Port        string
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

//LoadCfg load configuration data
func LoadCfg() {

	cfg, err := ini.Load(*cfgPath)
	if err != nil {
		logger.LogErr(err)
		return
	}

	err = cfg.Section("Master").MapTo(CfgMaster)
	if err != nil {
		logger.LogErr("The section Master's data cannot be load")
		return
	}

	err = cfg.Section("Worker").MapTo(CfgWorker)
	if err != nil {
		logger.LogErr("The section Worker's data cannot be load")
		return
	}

	err = cfg.ChildSections("Worker")[0].MapTo(CfgWorker.LoadReport)
	if err != nil {
		logger.LogErr("The section Worker.LoadReport's data cannot be load")
		return
	}
}

//SetCfgPath set the cfg Path
func SetCfgPath(path string) {

	cfgPath = &path

}
