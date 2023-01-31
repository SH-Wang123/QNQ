package common

import (
	"fmt"
	"sync"
)

var taskCount = 0
var countLock sync.Mutex

func GetSNCount() string {
	SNCountStr := fmt.Sprintf("%v", GetTaskCount())
	for i := len(SNCountStr); i < 4; i++ {
		SNCountStr = "0" + SNCountStr
	}
	return SNCountStr
}

func GetTaskCount() string {
	countLock.Lock()
	if taskCount == 9999 {
		taskCount = 0
	}
	taskCount++
	countLock.Unlock()
	return fmt.Sprintf("%v", taskCount)
}
