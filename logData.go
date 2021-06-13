package logger

type LogData struct {
	LogFilePath string
	LogTime     string
	Level       int
	LevelStr    string
	Message     string
	FileName    string
	LineNo      int
	FuncName    string
}
