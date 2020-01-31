package config

import (
	"testing"
)

func TestLoadINI(t *testing.T) {

	LoadCfg()

	t.Log(CfgMaster)
	t.Log(CfgWorker)
}
