package utils

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc/peer"
)

var (
	//WorkerID provide a unique id for Load Balance Worker.
	workerID uint32
)

//GenerateWorkerID generate a new work id for LB worker
func GenerateWorkerID() uint32 {

	atomic.AddUint32(&workerID, 1)
	return workerID
}

//GetContextIP get the IP of context
func GetContextIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("getClinetIP, invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("getClientIP, peer.Addr is nil")
	}

	switch addr := pr.Addr.(type) {
	case *net.UDPAddr:
		return addr.IP.String(), nil

	case *net.TCPAddr:
		return addr.IP.String(), nil
	default:
		return "", fmt.Errorf("getClientIP, addr type cannot be processed")

	}

}
