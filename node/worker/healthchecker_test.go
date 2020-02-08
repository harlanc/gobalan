package worker

import (
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/utils"
)

var (
	ip   = "localhost"
	port = 6689
)

func newServer() {

	t := time.NewTimer(time.Second * time.Duration(10))

	l, err := net.Listen("tcp", ip+":"+utils.Int2String(port))
	if err != nil {
		logger.LogErrf("Error listening:%v\n", err.Error())
		os.Exit(1)
	}

	logger.LogDebug("Listening on " + ip + ":" + utils.Int2String(port))

	go func() {
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				logger.LogErrf("Error accepting: %v\n", err.Error())
				return
			}
			// Handle connections in a new goroutine.
			go func() {
				defer conn.Close()
				out := make([]byte, 1024)
				var receivestr string
				if num, err := conn.Read(out); err == nil {
					receivestr = string(out[0:num])
					logger.LogDebugf("Server Receive: %s\n", receivestr)
				} else {
					logger.LogErr("Server Receive Error")
					return
				}
				if receivestr == "Send" {
					logger.LogDebug("Server send")
					if _, err := conn.Write([]byte("Receive")); err != nil {
					}
				} else {
					logger.LogDebug("Server receive is not Send.")
				}

			}()
		}
	}()

	go func() {
		select {
		case <-t.C:
			l.Close()
		}
	}()

}

func TestHealthCheck(t *testing.T) {

	logger.SetLogLevel(logger.Debug)

	newServer()
	hc := NewHealthChecker(ip, port, "", "")
	if !reflect.DeepEqual(hc.Check(), true) {
		t.Error("the check value is not right")
	} else {
		logger.LogInfo("The check value is right.")
	}

	ti := time.NewTimer(time.Second * time.Duration(10))

	select {
	case <-ti.C:
	}
	if !reflect.DeepEqual(hc.Check(), false) {
		t.Error("the check value is not right")
	}
}

func TestHCSendReceive(t *testing.T) {

	logger.SetLogLevel(logger.Debug)

	newServer()
	hc := NewHealthChecker(ip, port, "Send", "Receive")
	if !reflect.DeepEqual(hc.Check(), true) {
		t.Error("the check value is not right")
	} else {
		logger.LogInfo("The check value is right.")
	}

	hc1 := NewHealthChecker(ip, port, "Send", "Error")

	if !reflect.DeepEqual(hc1.Check(), false) {
		t.Error("the check value is not right")
	}
}
