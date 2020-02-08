package worker

import (
	"context"
	"testing"
	"time"

	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
)

func TestMonitor(t *testing.T) {

	//logger.SetLogLevel(logger.Debug)

	config.CfgWorker.LoadReport.LoadReportInterval = 2
	config.CfgWorker.LoadReport.NetworkAdapterName = "en0"
	config.CfgWorker.LoadReport.MaxNetworkBandwidth = 200

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	monitor := NewMonitor(ctx)
	monitor.Start()

	for {
		select {
		case val := <-monitor.StatPB:
			logger.LogDebug(val)
		case <-ctx.Done():
			logger.LogDebug("cancel is issued.")
			return
		}
	}

}
