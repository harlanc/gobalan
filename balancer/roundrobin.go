package balancer

import (
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/proto"
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
	return proto.BalanceType_RoundRobin.String()
}

//Pick pick a node
func (rr *RoundRobin) Pick() *node.Node {

	nl := node.NodeContainer.GetNodeList()
	n := nl[rr.next]
	rr.next = (rr.next + 1) % len(nl)
	return n
}
