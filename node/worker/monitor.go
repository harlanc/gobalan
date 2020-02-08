package worker

import (
	"context"
	"time"

	stat "github.com/akhenakh/statgo"
	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	pb "github.com/harlanc/gobalan/proto"
)

//Monitor monitor the server stat//https://zhuanlan.zhihu.com/p/35914450
type Monitor struct {
	Ctx        context.Context
	StatPB     chan pb.Stat
	Stat       *stat.Stat
	ticker     *time.Ticker
	adapterIdx int
}

//NewMonitor new a monitor
func NewMonitor(c context.Context) *Monitor {
	return &Monitor{StatPB: make(chan pb.Stat, 1),
		ticker:     time.NewTicker(time.Duration(config.CfgWorker.LoadReport.LoadReportInterval) * time.Second),
		adapterIdx: -1,
		Ctx:        c,
		Stat:       stat.NewStat()}
}

//Start start the monitor
func (m *Monitor) Start() {

	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.StatPB <- m.ReadStat()

			case <-m.Ctx.Done():
				return
			}
		}
	}()
}

//GetAdapterIndex get the network adapter index
func (m *Monitor) GetAdapterIndex(nis []*stat.NetIOStats) int {

	for i, v := range nis {

		if v.IntName == config.CfgWorker.LoadReport.NetworkAdapterName {
			return i
		}
	}
	return -1
}

//ReadStat read the server load
func (m *Monitor) ReadStat() pb.Stat {

	cpu := make(chan float32, 1)
	memory := make(chan float32, 1)
	bandwidthR := make(chan float32, 1)
	bandwidthW := make(chan float32, 1)

	go func() {
		m.Stat.CPUStats()
		time.Sleep(time.Second)
		cpus := m.Stat.CPUStats()

		//logger.LogDebugf("CPU Idle %f %f %f %f %f %f\n", cpus.User, cpus.Kernel, cpus.Idle, cpus.IOWait, cpus.Swap, cpus.Nice)
		cpu <- (100 - float32(cpus.Idle)) / 100 // range from 0 ~ 1

	}()

	go func() {

		m.Stat.MemStats()
		time.Sleep(time.Second)
		memorys := m.Stat.MemStats()

		memory <- (float32(memorys.Used) / float32(memorys.Total)) //range from 0 ~ 1

	}()

	go func() {

		m.Stat.NetIOStats()
		time.Sleep(time.Second)
		io := m.Stat.NetIOStats()

		if m.adapterIdx == -1 {
			m.adapterIdx = m.GetAdapterIndex(io)
			if m.adapterIdx == -1 {
				logger.LogErr("Cannot find the adpater name")
				return
			}
		}

		totalbandwidth := config.CfgWorker.LoadReport.MaxNetworkBandwidth

		bandwidthR <- (float32(io[m.adapterIdx].RX) / 1024 / 1024 * 8) / totalbandwidth // range from 0 ~ 1
		bandwidthW <- (float32(io[m.adapterIdx].TX) / 1024 / 1024 * 8) / totalbandwidth // range from 0 ~ 1

	}()

	return pb.Stat{CpuUsageRate: <-cpu, MemoryUsageRate: <-memory, ReadNetworkIOUsageRate: <-bandwidthR, WriteNetworkIOUsageRate: <-bandwidthW}

}
