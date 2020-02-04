package balancer

import (
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/proto"
)

var (
	// m is a map from name to balancer builder.
	m = make(map[string]Balancer)
	//CurrentBalanceType current balancer name
	CurrentBalanceType = proto.BalanceType_RoundRobin
)

//Register a Balancer
func Register(b Balancer) {
	m[b.Name()] = b
}

//Balancer interface
type Balancer interface {
	Pick() *node.Node
	Name() string
}
