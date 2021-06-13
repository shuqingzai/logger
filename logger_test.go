package logger

import (
	"testing"
	"time"
)

func TestFileLogger(t *testing.T) {
	logger := NewFileLogger("logs", "test123")
	defer logger.Close()
	logger.Debug("testDebug %d", time.Now().Unix())
	logger.Trace("testTrace %d", time.Now().Unix())
	logger.Info("testInfo %d", time.Now().Unix())
	logger.Warn("testWarn %d", time.Now().Unix())
	logger.Error("testError %d", time.Now().Unix())
	logger.Fatal("testFatal %d", time.Now().Unix())
}

func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger()
	defer logger.Close()
	logger.Debug("testDebug %d", time.Now().Unix())
	logger.Trace("testTrace %d", time.Now().Unix())
	logger.Info("testInfo %d", time.Now().Unix())
	logger.Warn("testWarn %d", time.Now().Unix())
	logger.Error("testError %d", time.Now().Unix())
	logger.Fatal("testFatal %d", time.Now().Unix())
}
