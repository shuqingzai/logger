package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// getLineInfo
// 获取函数调用栈信息
// @return fileName
// @return funcName
// @return lineNo
// @Author ShuQingZai<overbeck.jack@qq.com>
func getLineInfo(skip int) (fileName string, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}

// writeFile
// 写入文件
// @param fileHandle
// @param levelStr
// @param format
// @param args
// @Author ShuQingZai<overbeck.jack@qq.com>
func writeFile(fileHandle *os.File, skip int, levelStr string, format string, args ...interface{}) {
	// 文件写入
	nowTime := time.Now().Format("2006-01-02 15:04:05.999")
	fileName, funcName, lineNo := getLineInfo(skip)
	msg := fmt.Sprintf(format, args...)
	_, err := fmt.Fprintf(fileHandle, "%s %s [%s::%d %s()]\n%s\n", nowTime, levelStr, fileName, lineNo, funcName, msg)
	if err != nil {
		panic(fmt.Sprintf("Write Log error: %s", err))
	}
}

func checkDir(path string) error {
	dirPath, err := os.Stat(path)
	if err == nil {
		if dirPath.IsDir() {
			return nil
		}
		return fmt.Errorf("%s is Exist, but it is not a dir", path)
	}
	if os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return fmt.Errorf("checkDir error: %s", err)
}

func checkFileName(fileName string, fileSuffix string) string {
	if strings.HasSuffix(fileName, fileSuffix) {
		return fileName
	}

	return fmt.Sprintf("%s%s", fileName, fileSuffix)
}
