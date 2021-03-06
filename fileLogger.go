package logger

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileLogger struct {
	level         int
	logFilePath   string
	logFilename   string
	logFileExt    string
	splitType     int
	splitSize     int64
	lastSplitHour int
	fileMapSize   int
	logChanSize   int
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
	logFileExt := config["log_file_ext"]
	if logFileExt != ".log" && logFileExt != ".txt" {
		logFileExt = ".log"
	}

	fileMapSizeStr, ok := config["file_map_size"]
	if !ok {
		fileMapSizeStr = "30"
	}
	fileMapSize, err := strconv.Atoi(fileMapSizeStr)
	if err != nil || fileMapSize < 30 || fileMapSize > 300 {
		fileMapSize = 30
	}

	f := &FileLogger{
		level:         LogLevelDebug,
		logFilePath:   logPath,
		logFilename:   checkFileName(logFilename, logFileExt),
		logFileExt:    logFileExt,
		fileMapSize:   fileMapSize,
		fileMap:       make(map[string]*os.File, fileMapSize),
		logChanSize:   logChanSize,
		logDataChan:   make(chan *LogData, logChanSize),
		lastSplitHour: time.Now().Hour(),
	}
	// ??????????????????
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
		f.splitSize = logSplitSize * 1024 * 1024 // ?????????Mb
	} else {
		f.splitType = LogSplitTypeHour
	}

	// ??????????????????????????????
	go f.saveLogFile()

	return f, err
}

func (f *FileLogger) backupFile(logFilePath string, fileHandle *os.File, fileMapKey string, now *time.Time) (err error) {
	err = fileHandle.Close()
	if err != nil {
		return
	}
	delete(f.fileMap, fileMapKey) // ??????Map??????

	backupFilename := fmt.Sprintf("%s_%04d_%02d_%02d_%02d_%02d_%02d%s",
		strings.Replace(logFilePath, f.logFileExt, "", -1),
		now.Year(),
		now.Month(),
		now.Day(),
		f.lastSplitHour,
		now.Minute(),
		now.Second(),
		f.logFileExt,
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
		// ?????????????????????
		f.splitFileOfHour(logFilePath)
		return

	}

	// ???????????????????????????
	f.splitFileOfSize(logFilePath)
}

func (f *FileLogger) saveLogFile() {
	for logData := range f.logDataChan {
		// ??????????????????????????????
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
		// ????????????????????????fileMap
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
	// ?????????????????????????????????????????????
	filePath, levelStr := f.getFilePath(f.logFilename)
	// ??????????????????????????????fileMap???
	//fileHandle := f.getFileHandle(filePath)
	// ??????????????????chan??????
	logData := writeFile(4, level, format, args...)
	logData.LogFilePath = filePath
	logData.LevelStr = levelStr
	// ??????chan ??????chan?????????????????????
	select {
	case f.logDataChan <- logData:
	default:
	}
}
