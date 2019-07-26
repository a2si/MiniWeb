package DevLogs

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	COLOR_FORE_BLACK      = 30
	COLOR_FORE_RED        = 31
	COLOR_FORE_GREEN      = 32
	COLOR_FORE_YELLOW     = 33
	COLOR_FORE_BLUE       = 34
	COLOR_FORE_PURPLE_RED = 35
	COLOR_FORE_CYAN       = 36
	COLOR_FORE_WHITE      = 37
)
const (
	COLOR_BACK_BLACK      = 40
	COLOR_BACK_RED        = 41
	COLOR_BACK_GREEN      = 42
	COLOR_BACK_YELLOW     = 43
	COLOR_BACK_BLUE       = 44
	COLOR_BACK_PURPLE_RED = 45
	COLOR_BACK_CYAN       = 46
	COLOR_BACK_WHITE      = 47
)

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色

func getColorString(Color int, str string) string {
	if runtime.GOOS == "windows" || PRINT_COLOR == false {
		return str
	} else {
		// b 背景色彩 = 40-47
		// d 显示方式 = 0,1,4
		// f 前景色彩 = 30-37
		return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 40, Color, str, 0x1B)
	}
}

func getNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getFileProc(Level int) (string, string) {
	pc, file, _, ok := runtime.Caller(Level)
	if ok == false {
		return "Err", "Err"
	}
	//pc, file, line, ok := runtime.Caller(Level)
	//fmt.Println(runtime.FuncForPC(pc).Name(), file, line, ok)
	Temp := runtime.FuncForPC(pc).Name()
	Temp = strings.ToUpper(Temp)
	var File string
	var Proc string
	if strings.Contains(Temp, ".") {
		Arr := strings.Split(Temp, ".")
		File = string(Arr[0])
		Proc = string(Arr[1])
	} else {
		Arr := strings.Split(file, string(os.PathSeparator))
		File = Arr[len(Arr)-1]
		if strings.Contains(File, ".") {
			File = strings.Split(File, ".")[0]
		}
		File = strings.ToUpper(File)
		Proc = Temp
	}
	if strings.Contains(File, "/") {
		Arr := strings.Split(File, "/")
		File = Arr[len(Arr)-1]
	}
	return File, Proc
}
