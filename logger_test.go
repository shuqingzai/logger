package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestInitLogger(t *testing.T) {
	config := make(map[string]string, 8)
	//config["log_path"] = "logs"
	//config["log_name"] = "test1"
	//config["file_map_size"] = "30" // 30 -300
	//config["log_chan_size"] = "50000" // 1 - 10000
	//config["log_split_type"] = "size" // size || hour
	//config["log_split_size"] = "1" // 单位: MB
	//config["log_file_ext"] = ".log" // .txt || .log
	logger, err := InitLogger("file", config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Close()
	for {
		logger.Debug("testDebug %d", time.Now().Unix())
		logger.Trace("testTrace %d", time.Now().Unix())
		logger.Info("testInfo %d", time.Now().Unix())
		logger.Warn("testWarn %d", time.Now().Unix())
		logger.Error("testError %d", time.Now().Unix())
		logger.Fatal("testFatal %d", time.Now().Unix())
		time.Sleep(time.Second)
	}
}

func TestFileLogger(t *testing.T) {
	config := make(map[string]string, 8)
	logger, err := NewFileLogger(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Close()
	logger.Debug("testDebug %d", time.Now().Unix())
	logger.Trace("testTrace %d", time.Now().Unix())
	logger.Info("testInfo %d", time.Now().Unix())
	logger.Warn("testWarn %d", time.Now().Unix())
	logger.Error("testError %d", time.Now().Unix())
	logger.Fatal("testFatal %d", time.Now().Unix())
}

func TestConsoleLogger(t *testing.T) {

	config := map[string]string{}
	logger, err := NewConsoleLogger(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Close()
	logger.Debug("testDebug %d", time.Now().Unix())
	logger.Trace("testTrace %d", time.Now().Unix())
	logger.Info("testInfo %d", time.Now().Unix())
	logger.Warn("testWarn %d", time.Now().Unix())
	logger.Error("testError %d", time.Now().Unix())
	logger.Fatal("testFatal %d", time.Now().Unix())
}
