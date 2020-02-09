package logger

import (
	"fmt"
	"runtime"
	"time"
)

var (
	//curLogLevel
	curLogLevel LogLevel = Warn
)

//LogLevel log level
type LogLevel uint32

const (
	_ LogLevel = iota
	//Debug level 1
	Debug
	//Info level 2
	Info
	//Warn level 3
	Warn
	//Error level 4
	Error
)

const (
	infoColor    = "\033[1;34m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	debugColor   = "\033[0;36m%s\033[0m"
)

func generareLogTag(l LogLevel) interface{} {

	var tag string

	time := fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))

	switch l {
	case Debug:
		tag = fmt.Sprintf(debugColor, "[DEBUG] "+time)
	case Info:
		tag = fmt.Sprintf(infoColor, "[INFO ] "+time)
	case Warn:
		tag = fmt.Sprintf(warningColor, "[WARN ] "+time)
	case Error:

		tag = fmt.Sprintf(errorColor, "[ERROR] "+time)
	}

	return tag

}

func printLog(l LogLevel, s ...interface{}) {

	result := []interface{}{generareLogTag(l)}
	result = append(result, s...)

	fmt.Println(result...)
}

func printfLog(l LogLevel, f string, s ...interface{}) {

	result := []interface{}{generareLogTag(l)}
	result = append(result, " ")
	sf := fmt.Sprintf(f, s...)
	result = append(result, sf)

	fmt.Print(result...)
}

//SetLogLevel set log level
func SetLogLevel(l LogLevel) {

	curLogLevel = l
}

//LogDebug log debug
func LogDebug(s ...interface{}) {

	if curLogLevel <= Debug {
		printLog(Debug, s...)
	}
}

//LogDebugf log debug format
func LogDebugf(f string, s ...interface{}) {

	if curLogLevel <= Debug {
		printfLog(Debug, f, s...)
	}
}

//LogInfo log info
func LogInfo(s ...interface{}) {

	if curLogLevel <= Info {
		printLog(Info, s...)
	}
}

//LogInfof log info
func LogInfof(f string, s ...interface{}) {

	if curLogLevel <= Info {
		printfLog(Info, f, s...)
	}
}

//LogWarn log warn
func LogWarn(s ...interface{}) {

	if curLogLevel <= Warn {
		printLog(Warn, s...)
	}
}

//LogWarnf log info
func LogWarnf(f string, s ...interface{}) {

	if curLogLevel <= Warn {
		printfLog(Warn, f, s...)
	}
}

//LogErr log err
func LogErr(s ...interface{}) {

	if curLogLevel <= Error {
		_, fn, line, _ := runtime.Caller(1)
		errLineInfo := fmt.Sprintf("[%s:%d]", fn, line)
		allinfo := append([]interface{}{errLineInfo}, s...)
		printLog(Error, allinfo...)
	}
}

//LogErrf log info
func LogErrf(f string, s ...interface{}) {

	if curLogLevel <= Error {
		_, fn, line, _ := runtime.Caller(1)
		errLineInfo := fmt.Sprintf("[%s:%d]", fn, line)
		printfLog(Error, errLineInfo+f, s...)
	}
}
