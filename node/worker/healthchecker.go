package worker

import (
	"net"
	"strings"
	"time"

	"github.com/harlanc/gobalan/logger"
	"github.com/harlanc/gobalan/utils"
)

const (
	defaultTCPDailTimeout  = time.Duration(1) * time.Second
	defaultTCPReadTimeout  = time.Duration(2) * time.Second
	defaultTCPWriteTimeout = time.Duration(2) * time.Second
)

//HealthChecker Check worker service status
type HealthChecker struct {
	ip   string
	port int

	Receive string
	Send    string
}

//NewHealthChecker new a health checker.
func NewHealthChecker(ip string, port int, s string, r string) *HealthChecker {
	return &HealthChecker{port: port, Receive: r, Send: s}
}

//Check servcie is alive or not
func (hc *HealthChecker) Check() bool {
	conn, err := net.DialTimeout("tcp", hc.ip+":"+utils.Int2String(hc.port), defaultTCPDailTimeout)

	if err != nil {
		logger.LogWarnf("%v.\n", err)
		return false
	}

	defer conn.Close()
	if hc.Send == "" && hc.Receive == "" {
		return true
	}

	conn.SetReadDeadline(time.Now().Add(defaultTCPReadTimeout))
	conn.SetWriteDeadline(time.Now().Add(defaultTCPWriteTimeout))

	if _, err := conn.Write([]byte(hc.Send)); err != nil {
		logger.LogWarnf("%v\n", err)
		return false
	}

	out := make([]byte, 1024)
	if num, err := conn.Read(out); err == nil {
		receivestr := string(out[0:num])
		logger.LogDebugf("Client Receive: %v\n", receivestr)
		if strings.Compare(receivestr, hc.Receive) != 0 {
			logger.LogWarn("response did not match expected output")
			return false
		}
	} else {
		logger.LogWarnf("Receive:%v\n", err)
		return false
	}

	return true

}
