package balancer

import (
	"os"

	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/node"
	"github.com/harlanc/gobalan/proto"
)

var (
	// m is a map from name to balancer builder.
	m = make(map[string]Balancer)
	//CurrentBalanceType current balancer name
	CurrentBalanceType = proto.BalanceType_RoundRobin
)

//LoadBalanceType set balance type
func LoadBalanceType() {

	switch config.CfgMaster.LBAlgorithm {
	case "RR":
		CurrentBalanceType = proto.BalanceType_RoundRobin
	case "OP":
		CurrentBalanceType = proto.BalanceType_OptimalPerformance
	default:
		logger.LogErr("The Algorithm configured is not supported.")
		os.Exit(3)
	}
}

//Register a Balancer
func Register(b Balancer) {
	m[b.Name()] = b
}

//Balancer interface
type Balancer interface {
	Pick() *node.Node
	Name() string
}
