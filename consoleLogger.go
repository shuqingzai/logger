package logger

import (
	"fmt"
	"os"
)

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

func NewConsoleLogger(config map[string]string) (c LogInterface, err error) {
	c = &ConsoleLogger{}
	return
}

func (c *ConsoleLogger) writeLog(level int, format string, args ...interface{}) {
	logData := writeFile(4, level, format, args...)
	_, err := fmt.Fprintf(os.Stdout, "%s %s [%s::%d %s()]\n%s\n",
		logData.LogTime,
		getLogLevelText(logData.Level),
		logData.FileName,
		logData.LineNo,
		logData.FuncName,
		logData.Message)
	if err != nil {
		panic(fmt.Sprintf("Write Log error: %s", err))
	}
}
