package logger

import "fmt"

const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// getLogLevelText
// 获取日志级别描述 
// @param level
// @return string
// @Author ShuQingZai<overbeck.jack@qq.com>
func getLogLevelText(level int) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		panic(fmt.Sprintf("error Log Level: %d", level))
	}
}
