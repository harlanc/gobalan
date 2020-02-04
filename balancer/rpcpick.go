package balancer

import (
	"context"

	pb "github.com/harlanc/gobalan/proto"
)

// server is used to implement helloworld.GreeterServer.
type rpcPickServer struct{}

//NewRPCPickServer new a rpc pick server
func NewRPCPickServer() pb.RPCPickServer {
	return &rpcPickServer{}
}

// RPCPick implements proto.RPCPick
func (s *rpcPickServer) RPCPick(ctx context.Context, req *pb.PickRequest) (*pb.PickResponse, error) {

	nd := m[CurrentBalanceType.String()].Pick()
	return &pb.PickResponse{Ip: nd.IP, Port: nd.Port}, nil
}
