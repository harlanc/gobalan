package balancer

import (
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/proto"
	pb "github.com/harlanc/gobalan/proto"
)

var (
	//MaxPerformanceScore which can not be reached.
	MaxPerformanceScore float32 = 100
)

func init() {
	Register(NewOptimalPerformance())
}

//NewOptimalPerformance new RoundRobin balancer
func NewOptimalPerformance() Balancer {
	rr := &OptimalPerformance{}
	return rr
}

//OptimalPerformance balancer
type OptimalPerformance struct {
	next int
}

//Name get balancer name
func (op *OptimalPerformance) Name() string {
	return proto.BalanceType_OptimalPerformance.String()
}

//Pick pick a node
func (op *OptimalPerformance) Pick() *node.Node {

	nl := node.NodeContainer.GetNodeList()

	var idx int = -1
	var score float32 = MaxPerformanceScore

	for i, v := range nl {

		curscore := op.Score(v.Stat)
		if curscore < score {
			score = curscore
			idx = i
		}
	}
	if idx == -1 {
		return nil
	}

	return nl[idx]
}

//Score Score
func (op *OptimalPerformance) Score(s *pb.Stat) float32 {

	if s.GetCpuUsageRate() > 0.8 || s.GetMemoryUsageRate() > 0.8 || s.GetReadNetworkIOUsageRate() > 0.8 || s.GetWriteNetworkIOUsageRate() > 0.8 {
		return MaxPerformanceScore
	}
	return s.GetCpuUsageRate() + s.GetMemoryUsageRate() + s.GetReadNetworkIOUsageRate() + s.GetWriteNetworkIOUsageRate()
}
