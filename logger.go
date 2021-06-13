package logger

import "errors"

// InitLogger
// 初始化日志驱动
// @param driver file 文件驱动，console 控制台驱动
// @param config 配置
// @return log
// @return err
// @Author ShuQingZai<overbeck.jack@qq.com>
func InitLogger(driver string, config map[string]string) (log LogInterface, err error) {
	switch driver {
	case "file":
		log, err = NewFileLogger(config)
	case "console":
		log, err = NewConsoleLogger(config)
	default:
		log = nil
		err = errors.New("this driver is not currently supported")
	}
	return
}
