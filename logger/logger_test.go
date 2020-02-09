package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	curLogLevel = Debug

	LogDebug("Log Debug", "abcd", "123q3e")
	LogDebugf("+++%s+++\n", "Log Debugf")

	LogInfo("Log Info")
	LogInfof("--%s--\n", "Log Infof")

	LogWarn("Log Warn")
	LogWarnf("--%s--\n", "Log Warnf")

	LogErr("Log Err")
	LogErrf("--%s--\n", "Log Errfdad")
}
