package config

import (
	"testing"
)

func TestLoadINI(t *testing.T) {

	t.Log(CfgMaster)
	t.Log(CfgWorker)
	t.Log(CfgWorker.LoadReport)
}
