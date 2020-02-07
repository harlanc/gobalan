package worker

import (
	"net"
	"time"

	"github.com/google/seesaw/healthcheck"
)

var (
	tcpCheckTimeout = time.Duration(2) * time.Second
)

//HealthChecker Check worker service status
type HealthChecker struct {
	tcpChecker *healthcheck.TCPChecker
}

//NewHealthChecker new a health checker.
func NewHealthChecker(ip net.IP, port int) *HealthChecker {
	return &HealthChecker{tcpChecker: healthcheck.NewTCPChecker(ip, port)}
}

//IsServiceAlive judge servcie is alive or not
func (hc *HealthChecker) IsServiceAlive() bool {

	result := hc.tcpChecker.Check(tcpCheckTimeout)
	return result.Success
}
