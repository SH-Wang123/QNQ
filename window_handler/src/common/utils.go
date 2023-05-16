package common

import (
	"fmt"
	"strings"
	"sync"
	"time"
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

func GetNowTimeStr() string {
	now := time.Now()
	ret := fmt.Sprintf("%v", now.Format("2006/01/02 15:04:05"))
	return ret
}

func GetIpFromAddr(addr string) string {
	s := strings.Split(addr, ":")
	if len(s) == 2 {
		return s[0]
	} else {
		return ""
	}
}
