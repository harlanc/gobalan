package balancer

import (
	"context"
	"errors"

	"github.com/harlanc/gobalan/logger"
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
	if nd == nil {
		return nil, errors.New("There is no proper node to return")
	}
	logger.LogDebugf("###########The worker id %d is picked\n", nd.WorkerID)
	return &pb.PickResponse{Ip: nd.IP, Port: nd.Port, WorkerId: nd.WorkerID}, nil
}
