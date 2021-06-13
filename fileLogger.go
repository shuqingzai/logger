package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	level   int
	logPath string
	logName string
	fileMap map[string]*os.File
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

func NewFileLogger(logPath, logName string) LogInterface {
	checkDir(logPath)
	f := &FileLogger{
		level:   LogLevelDebug,
		logPath: logPath,
		logName: logName,
		fileMap: make(map[string]*os.File, 20),
	}

	return f
}

func checkDir(path string) {
	dirPath, err := os.Stat(path)
	if err == nil {
		if dirPath.IsDir() {
			return
		}
		panic(fmt.Sprintf("%s is Exist, but it is not a dir", path))
	}
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err == nil {
			return
		}
		panic(fmt.Sprintf("mkdir error: %s", err))
	}
	panic(fmt.Sprintf("checkDir error: %s", err))
}

func (f *FileLogger) openFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Open file %s failed, err: %s ", filePath, err))
	}

	return file
}

func (f *FileLogger) getFileHandle(filePath string) *os.File {
	fileMapKey := path.Base(filePath)
	fileHandle, ok := f.fileMap[fileMapKey]
	if !ok {
		// 打开文件并存储在fileMap
		fileHandle = f.openFile(filePath)
		f.fileMap[fileMapKey] = fileHandle
	}

	return fileHandle
}

func (f *FileLogger) getFilePath(fileName string) (filePath string, levelStr string) {
	dateTime := time.Now().Format("2006-01-02")
	levelStr = getLogLevelText(f.level)
	filePath = fmt.Sprintf("%s/%s_%s_%s", f.logPath, dateTime, levelStr, fileName)
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
	filePath, levelStr := f.getFilePath(f.logName)
	// 获取文件句柄并存储在fileMap中
	fileHandle := f.getFileHandle(filePath)
	// 文件写入
	writeFile(fileHandle, 4, levelStr, format, args...)
}
