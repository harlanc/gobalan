package master

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/harlanc/gobalan/balancer"
	"github.com/harlanc/gobalan/config"
	"github.com/harlanc/gobalan/logger"
	pb "github.com/harlanc/gobalan/proto"
)

//MasterServer the master server struct
type MasterServer struct {
	rpcs  *grpc.Server
	rpcps pb.RPCPickServer
	rpcws pb.WatchServer
}

//NewMasterServer new a master server
func NewMasterServer() *MasterServer {

	s := grpc.NewServer()
	ms := &MasterServer{

		rpcs:  s,
		rpcws: NewWatchServer(),
		rpcps: balancer.NewRPCPickServer(),
	}
	return ms
}

//Run Run the master server
func (ms *MasterServer) Run() {

	masterAddr := fmt.Sprintf("192.168.0.104:%s", config.CfgMaster.Port)

	listener, err := net.Listen("tcp", masterAddr)
	if err != nil {
		logger.LogErrf("failed to listen: %v", err)
		return
	}
	logger.LogInfof("rpc listening on:%s\n", masterAddr)

	pb.RegisterRPCPickServer(ms.rpcs, ms.rpcps)
	pb.RegisterWatchServer(ms.rpcs, ms.rpcws)

	ms.rpcs.Serve(listener)
}

//Stop the master server
func (ms *MasterServer) Stop() {
	logger.LogInfo("Master server is stoped.")
	ms.rpcs.Stop()
}
