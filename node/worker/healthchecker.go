package worker

import (
	"net"
	"strings"
	"time"

	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/utils"
)

const (
	defaultTCPDailTimeout      = time.Duration(1) * time.Second
	defaultTCPReadTimeout      = time.Duration(2) * time.Second
	defaultTCPWriteTimeout     = time.Duration(2) * time.Second
	defaultHealthCheckInterval = time.Duration(2) * time.Second
)

//HealthChecker Check worker service status
type HealthChecker struct {
	ip   string
	port int

	receive string
	send    string

	serviceUp chan bool
	closec    chan struct{}
}

//NewHealthChecker new a health checker.
func NewHealthChecker(ip string, port int, s string, r string) *HealthChecker {
	if port == -1 {
		return nil
	}
	return &HealthChecker{port: port, receive: r, send: s, serviceUp: make(chan bool, 1), closec: make(chan struct{}, 1)}
}

//run return the checkout value
func (hc *HealthChecker) run() {
	ticker := time.NewTicker(defaultHealthCheckInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				hc.serviceUp <- hc.check()
			case <-hc.closec:
				return
			}
		}
	}()
}

//stop the health checker
func (hc *HealthChecker) stop() {
	hc.closec <- struct{}{}
}

//check servcie is alive or not
func (hc *HealthChecker) check() bool {
	conn, err := net.DialTimeout("tcp", hc.ip+":"+utils.Int2String(hc.port), defaultTCPDailTimeout)

	if err != nil {
		logger.LogWarnf("%v.\n", err)
		return false
	}

	defer conn.Close()
	if hc.send == "" && hc.receive == "" {
		return true
	}

	conn.SetReadDeadline(time.Now().Add(defaultTCPReadTimeout))
	conn.SetWriteDeadline(time.Now().Add(defaultTCPWriteTimeout))

	if _, err := conn.Write([]byte(hc.send)); err != nil {
		logger.LogWarnf("%v\n", err)
		return false
	}

	out := make([]byte, 1024)
	if num, err := conn.Read(out); err == nil {
		receivestr := string(out[0:num])
		logger.LogDebugf("Client Receive: %v\n", receivestr)
		if strings.Compare(receivestr, hc.receive) != 0 {
			logger.LogWarn("response did not match expected output")
			return false
		}
	} else {
		logger.LogWarnf("Receive:%v\n", err)
		return false
	}

	return true

}
