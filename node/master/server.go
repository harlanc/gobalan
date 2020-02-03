package master

import (
	"log"
	"net"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	pb "github.com/harlanc/gobalan/proto"
)

//LBMasterServer the master server struct
type LBMasterServer struct {
	addr string
	rpcs *grpc.Server
	hbs  pb.HeartbeatServer
}

//NewLBMasterServer new a master server
func NewLBMasterServer(addr string) *LBMasterServer {

	s := grpc.NewServer()
	ms := &LBMasterServer{
		addr: addr,
		rpcs: s,
	}
	return ms
}

//Run Run the master server
func (ms *LBMasterServer) Run() {

	listener, err := net.Listen("tcp", ms.addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	log.Printf("rpc listening on:%s", ms.addr)

	pb.RegisterHeartbeatServer(ms.rpcs, ms.hbs)
	ms.rpcs.Serve(listener)

}

//SendHeartbeat send heartbeat
func (ms *LBMasterServer) SendHeartbeat(srv pb.Heartbeat_SendHeartbeatServer) error {

	return status.Errorf(codes.Unimplemented, "method SendHeartbeat not implemented")
}
