package logger

import "os"

type ConsoleLogger struct {
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	c.writeLog(LogLevelDebug, format, args...)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	c.writeLog(LogLevelTrace, format, args...)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	c.writeLog(LogLevelInfo, format, args...)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	c.writeLog(LogLevelWarn, format, args...)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	c.writeLog(LogLevelError, format, args...)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	c.writeLog(LogLevelFatal, format, args...)
}

func (c *ConsoleLogger) Close() {

}

func NewConsoleLogger() LogInterface {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) writeLog(level int, format string, args ...interface{}) {
	// 文件写入
	writeFile(os.Stdout, 4, getLogLevelText(level), format, args...)
}
