package balancer

import (
	"context"

	pb "github.com/harlanc/gobalan/proto"
	"google.golang.org/grpc"
)

//RPCPickClient rpc pick client
type RPCPickClient struct {
	ctx    context.Context
	client pb.RPCPickClient
}

//NewRPCPickClient NewWatcherClient
func NewRPCPickClient(c *grpc.ClientConn) *RPCPickClient {
	return &RPCPickClient{client: pb.NewRPCPickClient(c), ctx: context.Background()}
}

// RPCPick call client RPC Pick
func (c *RPCPickClient) RPCPick() (*pb.PickResponse, error) {
	return c.client.RPCPick(c.ctx, &pb.PickRequest{})
}
