package worker

import (
	"context"
	"fmt"
	"time"

	stat "github.com/akhenakh/statgo"
	"github.com/harlanc/gobalan/config"
	pb "github.com/harlanc/gobalan/proto"
)

//Monitor monitor the server stat
type Monitor struct {
	Ctx        context.Context
	StatPB     chan pb.Stat
	stat       *stat.Stat
	ticker     *time.Ticker
	adapterIdx int
}

//NewMonitor new a monitor
func NewMonitor(c context.Context) *Monitor {
	return &Monitor{StatPB: make(chan pb.Stat, 1),
		ticker:     time.NewTicker(time.Duration(config.CfgWorker.LoadReportInterval) * time.Second),
		adapterIdx: -1,
		Ctx:        c,
		stat:       stat.NewStat()}
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

		if v.IntName == config.CfgWorker.NetworkAdapterName {
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

	go func(c chan float32, s *stat.Stat) {

		s.CPUStats()
		time.Sleep(time.Second)
		cpus := s.CPUStats()

		c <- (100 - float32(cpus.Idle)) / 100 // range from 0 ~ 1

	}(cpu, m.stat)

	go func(m chan float32, s *stat.Stat) {

		s.MemStats()
		time.Sleep(time.Second)
		memorys := s.MemStats()

		m <- (float32(memorys.Used) / float32(memorys.Total)) //range from 0 ~ 1

	}(memory, m.stat)

	go func(br chan float32, bw chan float32, s *stat.Stat) {

		s.NetIOStats()
		time.Sleep(time.Second)
		io := s.NetIOStats()

		if m.adapterIdx == -1 {
			m.adapterIdx = m.GetAdapterIndex(io)
			if m.adapterIdx == -1 {
				fmt.Println("Cannot find the adpater name")
				return
			}
		}

		totalbandwidth := config.CfgWorker.MaxNetworkBandwidth
		fmt.Printf("R %d KB\n", io[m.adapterIdx].RX/1024)

		br <- (float32(io[m.adapterIdx].RX) / 1024 / 1024 * 8) / totalbandwidth // range from 0 ~ 1
		bw <- (float32(io[m.adapterIdx].TX) / 1024 / 1024 * 8) / totalbandwidth // range from 0 ~ 1

	}(bandwidthR, bandwidthW, m.stat)

	return pb.Stat{CpuUsageRate: <-cpu, MemoryUsageRate: <-memory, ReadNetworkIOUsageRate: <-bandwidthR, WriteNetworkIOUsageRate: <-bandwidthW}

}
