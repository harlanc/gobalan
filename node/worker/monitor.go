package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/harlanc/gobalan/config"
	pb "github.com/harlanc/gobalan/proto"
	"github.com/harlanc/gobalan/stat"
)

//Monitor monitor the server stat
type Monitor struct {
	Ctx        context.Context
	Stat       chan pb.Stat
	ticker     *time.Ticker
	adapterIdx int
}

//NewMonitor new a monitor
func NewMonitor(c context.Context) *Monitor {
	return &Monitor{Stat: make(chan pb.Stat, 1),
		ticker:     time.NewTicker(time.Duration(config.CfgWorker.LoadReportInterval) * time.Second),
		adapterIdx: -1,
		Ctx:        c}
}

//Start start the monitor
func (m *Monitor) Start() {

	for {
		select {
		case <-m.ticker.C:
			m.Stat <- m.ReadStat()
		case <-m.Ctx.Done():
			return
		}
	}
}

//GetAdapterIndex get the network adapter index
func (m *Monitor) GetAdapterIndex(nis []*stat.NetIOStats) int {

	for i, v := range nis {

		if v.IntName == config.CfgWorker.NetworkAdapterName {
			return i
		}

	}
	return -1
}

//ReadStat read the server load
func (m *Monitor) ReadStat() pb.Stat {

	st := stat.NewStat()

	time.Sleep(time.Millisecond * 100)

	var cpu chan float32
	var memory chan float32
	var bandwidthR chan float32
	var bandwidthW chan float32

	go func(c chan float32, s *stat.Stat) {
		cpus := s.CPUStats()
		c <- (100 - float32(cpus.Idle)) / 100 // range from 0 ~ 1

	}(cpu, st)

	go func(m chan float32, s *stat.Stat) {
		memorys := s.MemStats()
		m <- (float32(memorys.Used) / float32(memorys.Total)) //range from 0 ~ 1

	}(memory, st)

	go func(br chan float32, bw chan float32, s *stat.Stat) {

		io := s.NetIOStats()

		if m.adapterIdx == -1 {
			m.adapterIdx = m.GetAdapterIndex(io)
			if m.adapterIdx == -1 {
				fmt.Println("Cannot find the adpater name")
				return
			}
		}

		totalbandwidth := config.CfgWorker.MaxNetworkBandwidth

		br <- (float32(io[m.adapterIdx].RX) / 1024 / 1024 * 8) / totalbandwidth // range from 0 ~ 1
		bw <- (float32(io[m.adapterIdx].TX) / 1024 / 1024 * 8) / totalbandwidth // range from 0 ~ 1

	}(bandwidthR, bandwidthW, st)

	return pb.Stat{CpuUsageRate: <-cpu, MemoryUsageRate: <-memory, ReadNetworkIOUsageRate: <-bandwidthR, WriteNetworkIOUsageRate: <-bandwidthW}

}
