package DevLogs

func Debug(LogString string) {
	logPrintCenter(LEVEL_DEBUG, LogString)
}
func Info(LogString string) {
	logPrintCenter(LEVEL_INFO, LogString)
}
func Warn(LogString string) {
	logPrintCenter(LEVEL_WARN, LogString)
}
func Error(LogString string) {
	logPrintCenter(LEVEL_ERROR, LogString)
}
func Fatal(LogString string) {
	logPrintCenter(LEVEL_FATAL, LogString)
}
func Panic(LogString string) {
	logPrintCenter(LEVEL_PANIC, LogString)
}

func logPrintCenter(Level int, Msg string) {
	var LevelMsg string
	switch Level {
	default:
		LevelMsg = "UNKNOW"
	case LEVEL_DEBUG:
		LevelMsg = "DEBUG"
	case LEVEL_INFO:
		LevelMsg = "INFO"
	case LEVEL_WARN:
		LevelMsg = "WARN"
	case LEVEL_ERROR:
		LevelMsg = "ERROR"
	case LEVEL_FATAL:
		LevelMsg = "FATAL"
	case LEVEL_PANIC:
		LevelMsg = "PANIC"
	}
	printConsoleFileAndProc(LevelMsg, Msg)
}
