package master

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/harlanc/gobalan/balancer"
	pb "github.com/harlanc/gobalan/proto"
)

//MasterServer the master server struct
type MasterServer struct {
	addr  string
	rpcs  *grpc.Server
	rpcps pb.RPCPickServer
	rpcws pb.WatchServer
}

//NewMasterServer new a master server
func NewMasterServer(addr string) *MasterServer {

	s := grpc.NewServer()
	ms := &MasterServer{
		addr:  addr,
		rpcs:  s,
		rpcws: NewWatchServer(),
		rpcps: balancer.NewRPCPickServer(),
	}
	return ms
}

//Run Run the master server
func (ms *MasterServer) Run() {

	listener, err := net.Listen("tcp", ms.addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	log.Printf("rpc listening on:%s", ms.addr)

	pb.RegisterRPCPickServer(ms.rpcs, ms.rpcps)
	pb.RegisterWatchServer(ms.rpcs, ms.rpcws)

	ms.rpcs.Serve(listener)
}
