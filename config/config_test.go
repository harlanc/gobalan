package config

import (
	"testing"
)

func TestLoadINI(t *testing.T) {

	// LoadCfg()

	SetCfgPath("/Users/zexu/go/src/github.com/harlanc/gobalan/config/config.ini")
	LoadCfg()

	t.Log(CfgMaster)
	t.Log(CfgWorker)
}
