package logger

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logFilePath   string
	logFilename   string
	splitType     int
	splitSize     int64
	lastSplitHour int
	fileMap       map[string]*os.File
	logDataChan   chan *LogData
}

func (f *FileLogger) Close() {
	for k, file := range f.fileMap {
		if file != nil {
			_ = file.Close()
		}
		delete(f.fileMap, k)
	}
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.writeLog(LogLevelDebug, format, args...)
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	f.writeLog(LogLevelTrace, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	f.writeLog(LogLevelInfo, format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	f.writeLog(LogLevelWarn, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	f.writeLog(LogLevelError, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.writeLog(LogLevelFatal, format, args...)
}

func NewFileLogger(config map[string]string) (LogInterface, error) {
	logPath, ok := config["log_path"]
	if !ok {
		logPath = "logs"
	}

	err := checkDir(logPath)
	if err != nil {
		return nil, err
	}

	logFilename, ok := config["log_name"]
	if !ok {
		logFilename = "logger"
	}

	logChanSizeStr, ok := config["log_chan_size"]
	if !ok {
		logChanSizeStr = "50000"
	}
	logChanSize, err := strconv.Atoi(logChanSizeStr)
	if err != nil || logChanSize < 1 || logChanSize > 100000 {
		logChanSize = 50000
	}

	f := &FileLogger{
		level:         LogLevelDebug,
		logFilePath:   logPath,
		logFilename:   checkFileName(logFilename, ".log"),
		fileMap:       make(map[string]*os.File, 20),
		logDataChan:   make(chan *LogData, logChanSize),
		lastSplitHour: time.Now().Hour(),
	}
	// 日志切割方式
	logSplitTypeStr, ok := config["log_split_type"]
	if !ok {
		logSplitTypeStr = "size"
	}
	if logSplitTypeStr == "size" {
		f.splitType = LogSplitTypeSize
		logSplitSizeStr, ok := config["log_split_size"]
		if !ok {
			logSplitSizeStr = "20"
		}
		logSplitSize, err := strconv.ParseInt(logSplitSizeStr, 10, 64)
		if err != nil {
			logSplitSize = 20
		}
		f.splitSize = logSplitSize * 1024 * 1024 // 单位是Mb
	} else {
		f.splitType = LogSplitTypeHour
	}

	// 启动协程写入日志文件
	go f.saveLogFile()

	return f, err
}

func (f *FileLogger) backupFile(logFilePath string, fileHandle *os.File, fileMapKey string, now *time.Time) (err error) {
	err = fileHandle.Close()
	if err != nil {
		return
	}
	delete(f.fileMap, fileMapKey) // 删除Map元素

	backupFilename := fmt.Sprintf("%s_%04d_%02d_%02d_%02d_%02d_%02d",
		logFilePath,
		now.Year(),
		now.Month(),
		now.Day(),
		f.lastSplitHour,
		now.Minute(),
		now.Second(),
	)
	f.lastSplitHour = now.Hour()
	err = os.Rename(logFilePath, backupFilename)

	return
}

func (f FileLogger) splitFileOfHour(logFilePath string) {
	now := time.Now()
	if now.Hour() == f.lastSplitHour {
		return
	}
	fileHandle, ok, fileMapKey := f.checkFileHandle(logFilePath)
	if !ok {
		return
	}
	_ = f.backupFile(logFilePath, fileHandle, fileMapKey, &now)
}

func (f *FileLogger) splitFileOfSize(logFilePath string) {
	fileHandle, ok, fileMapKey := f.checkFileHandle(logFilePath)
	if !ok {
		return
	}
	statInfo, err := fileHandle.Stat()
	if err != nil {
		return
	}
	if statInfo.Size() < f.splitSize {
		return
	}
	now := time.Now()
	_ = f.backupFile(logFilePath, fileHandle, fileMapKey, &now)
}

func (f *FileLogger) checkSplitFile(logFilePath string) {
	if f.splitType == LogSplitTypeHour {
		// 按小时切割文件
		f.splitFileOfHour(logFilePath)
		return

	}

	// 按文件大小切割文件
	f.splitFileOfSize(logFilePath)
}

func (f *FileLogger) saveLogFile() {
	for logData := range f.logDataChan {
		// 检查文件是否需要切割
		f.checkSplitFile(logData.LogFilePath)
		fileHandle := f.getFileHandle(logData.LogFilePath)
		_, _ = fmt.Fprintf(fileHandle, "%s %s [%s::%d %s()]\n%s\n",
			logData.LogTime,
			logData.LevelStr,
			logData.FileName,
			logData.LineNo,
			logData.FuncName,
			logData.Message)
	}
}

func (f *FileLogger) openFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Open file %s failed, err: %s ", filePath, err))
	}

	return file
}

func (f *FileLogger) getFileHandle(filePath string) *os.File {
	fileHandle, ok, fileMapKey := f.checkFileHandle(filePath)
	if !ok {
		// 打开文件并存储在fileMap
		fileHandle = f.openFile(filePath)
		f.fileMap[fileMapKey] = fileHandle
	}

	return fileHandle
}

func (f *FileLogger) checkFileHandle(filePath string) (fileHandle *os.File, ok bool, fileMapKey string) {
	fileMapKey = path.Base(filePath)
	fileHandle, ok = f.fileMap[fileMapKey]
	return
}

func (f *FileLogger) getFilePath(fileName string) (filePath string, levelStr string) {
	dateTime := time.Now().Format("2006-01-02")
	levelStr = getLogLevelText(f.level)
	filePath = fmt.Sprintf("%s/%s_%s_%s", f.logFilePath, dateTime, levelStr, fileName)
	return
}

func (f *FileLogger) setLogLevel(level int) {
	if level >= LogLevelDebug && level <= LogLevelFatal {
		f.level = level
	}
}

func (f *FileLogger) writeLog(level int, format string, args ...interface{}) {
	f.setLogLevel(level)
	// 解析获取文件路径与错误级别描述
	filePath, levelStr := f.getFilePath(f.logFilename)
	// 获取文件句柄并存储在fileMap中
	//fileHandle := f.getFileHandle(filePath)
	// 组装信息放入chan队列
	logData := writeFile(4, level, format, args...)
	logData.LogFilePath = filePath
	logData.LevelStr = levelStr
	// 放入chan 如果chan已满，直接跳过
	select {
	case f.logDataChan <- logData:
	default:
	}
}
