package DevLogs

import (
	"fmt"
	"runtime"
	"strings"
)

func printConsoleFileAndProc(LevelMsg string, Msg string) {
	pc, _, _, _ := runtime.Caller(4)
	Proc := runtime.FuncForPC(pc).Name()
	dwRet := getColorString(COLOR_FORE_BLUE, "[ ")
	dwRet = dwRet + getColorString(COLOR_FORE_BLACK, rjust("PROC", 6))
	dwRet = dwRet + getColorString(COLOR_FORE_BLUE, " ] ")
	dwRet = dwRet + getColorString(COLOR_FORE_PURPLE_RED, strings.ToUpper(Proc)+"\n")
	dwRet = dwRet + getColorString(COLOR_FORE_BLUE, "[ ")
	dwRet = dwRet + getColorString(COLOR_FORE_YELLOW, rjust(LevelMsg, 6))
	dwRet = dwRet + getColorString(COLOR_FORE_BLUE, " ] ")
	dwRet = dwRet + getColorString(COLOR_FORE_WHITE, Msg+"\n")
	if PRINT_LOG {
		fmt.Print(dwRet)
	}
}

func rjust(str string, need int) string {
	if len(str) >= need {
		return string([]byte(str)[need])
	}
	iCount := need - len(str)
	return strings.Repeat(" ", iCount) + str
}

func ljust(str string, need int) string {
	if len(str) >= need {
		return string([]byte(str)[need])
	}
	iCount := need - len(str)
	return str + strings.Repeat(" ", iCount)
}

// DEBUG len(str)%2 == 1
func center(str string, need int) string {
	if len(str) >= need {
		return string([]byte(str)[:need])
	}
	iCount := need - len(str)
	if iCount%2 == 1 {
		iCount = iCount / 2
		strLeft := iCount - 1
		if strLeft < 0 {
			strLeft = 0
		}
		return strings.Repeat(" ", strLeft) + str + strings.Repeat(" ", iCount)
	}
	iCount = iCount / 2
	return strings.Repeat(" ", iCount) + str + strings.Repeat(" ", iCount)
}
