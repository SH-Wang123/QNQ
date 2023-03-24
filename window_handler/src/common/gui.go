//Decoupling gui with worker

package common

import (
	"sync"
)

var CLI_FALG = false

const (
	LOCAL_BATCH_POLICY_RUNNING = iota
	LOCAL_BATCH_POLICY_STOP
	LOCAL_SINGLE_POLICY_RUNNING
	LOCAL_SINGLE_POLICY_STOP
	TEST_DISK_SPEED_START
	TEST_DISK_SPEED_OVER
)

var (
	CurrentLocalPartSN string
	LocalPartStartLock = &sync.WaitGroup{}
)

var (
	CurrentLocalBatchSN string
	LocalBatchStartLock = &sync.WaitGroup{}
)
var (
	currentFileMap = make(map[string]string)
	cfLock         = &sync.RWMutex{}
)

var gwLock sync.RWMutex

var GWChannel = make(chan int)

var LocalBatchPolicyRunningFlag = false
var LocalSinglePolicyRunningFlag = false

func SendSignal2GWChannel(signal int) {
	gwLock.Lock()
	defer gwLock.Unlock()
	if CLI_FALG {

	} else {
		GWChannel <- signal
	}
}

func SetCurrentSyncFile(sn string, typeStr string, fileName string) {
	cfLock.Lock()
	defer cfLock.Unlock()
	s := []rune(typeStr + fileName)
	if len(s) > 60 {
		currentFileMap[sn] = string(s[0:60]) + " ..."
	} else {
		currentFileMap[sn] = string(s)
	}

}

func GetCurrentSyncFile(sn string) string {
	cfLock.Lock()
	defer cfLock.Unlock()
	return currentFileMap[sn]
}
