package worker

import (
	"context"
	"testing"
	"time"

	"github.com/harlanc/gobalan/config"
)

func TestMonitor(t *testing.T) {

	config.SetCfgPath("/Users/zexu/go/src/github.com/harlanc/gobalan/config/config.ini")
	config.LoadCfg()

	config.CfgWorker.LoadReportInterval = 2
	config.CfgWorker.NetworkAdapterName = "en0"
	config.CfgWorker.MaxNetworkBandwidth = 200

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()

	monitor := NewMonitor(ctx)
	monitor.Start()

	for {
		select {
		case val := <-monitor.StatPB:
			t.Log(val)
		case <-ctx.Done():
			t.Log("cancel is issued.")
			return
		}
	}

}
