package logger

import (
	"fmt"
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
	LogErrf("--%s--\n", "Log Errf")
}

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func TestColor(t *testing.T) {
	fmt.Printf(InfoColor, "Info")
	fmt.Println("")
	fmt.Printf(NoticeColor, "Notice")
	fmt.Println("")
	fmt.Printf(WarningColor, "Warning")
	fmt.Println("")
	fmt.Printf(ErrorColor, "Error")
	fmt.Println("")
	fmt.Printf(DebugColor, "Debug")
	fmt.Println("")
}
