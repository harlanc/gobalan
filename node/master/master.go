package master

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/harlanc/gobalan/balancer"
	"github.com/harlanc/gobalan/config"
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

	masterAddr := fmt.Sprintf("0.0.0.0:%s", config.CfgMaster.MasterPort)

	listener, err := net.Listen("tcp", masterAddr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	log.Printf("rpc listening on:%s", masterAddr)

	pb.RegisterRPCPickServer(ms.rpcs, ms.rpcps)
	pb.RegisterWatchServer(ms.rpcs, ms.rpcws)

	ms.rpcs.Serve(listener)
}
