package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	//CfgMaster exported Master configurations
	CfgMaster *cfgMaster = new(cfgMaster)
	//CfgWorker exported Worker configurations
	CfgWorker *cfgWorker = new(cfgWorker)
	//CfgPath the configuration file full path
	cfgPath *string = new(string)
)

//cfgMaster configurations
type cfgMaster struct {
	IsMaster    bool
	MasterPort  string
	ServicePort string
}

//Worker configurations
type cfgWorker struct {
	IsWorker            bool
	MasterIP            string
	MasterPort          string
	ServicePort         string
	HeartbeatInterval   uint
	LoadReportInterval  uint
	MaxNetworkBandwidth float32
	NetworkAdapterName  string
}

// func init() {

// 	LoadCfg()

// }

//LoadCfg load configuration data
func LoadCfg() {

	cfg, err := ini.Load(*cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cfg.Section("Master").MapTo(CfgMaster)
	if err != nil {
		fmt.Println("The section Master's data cannot be load")
		return
	}

	err = cfg.Section("Worker").MapTo(CfgWorker)
	if err != nil {
		fmt.Println("The section Worker's data cannot be load")
		return
	}
}

//SetCfgPath set the cfg Path
func SetCfgPath(path string) {

	cfgPath = &path

}
