package balancer

import (
	"github.com/harlanc/gobalan/node"
	pb "github.com/harlanc/gobalan/proto"
)

func init() {
	Register(NewRoundRobin())
}

//NewRoundRobin new RoundRobin balancer
func NewRoundRobin() Balancer {
	rr := &RoundRobin{next: 0}
	return rr
}

//RoundRobin balancer
type RoundRobin struct {
	next int
}

//Name get balancer name
func (rr *RoundRobin) Name() string {
	return pb.BalanceType_RoundRobin.String()
}

//Pick pick a node
func (rr *RoundRobin) Pick() *node.Node {

	nl := node.NodeContainer.GetNodeList()
	l := len(nl)

	var rv *node.Node
	for i := 0; i < l; i++ {
		rv = nl[rr.next%l]
		rr.next = (rr.next + 1) % l
		if rv.ServiceStatus == pb.ServiceStatus_Up {
			return rv
		}
	}

	return nil
}
