package DevLogs

import (
	"sync"
)

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_PANIC
)

var (
	PRINT_COLOR bool = false
	PRINT_LOG   bool = true
)

var lock *sync.RWMutex // 读写锁

func init() {
	PRINT_COLOR = true
	lock = new(sync.RWMutex)
}

func LogsEnable(Enable bool) {
	PRINT_LOG = Enable
}
