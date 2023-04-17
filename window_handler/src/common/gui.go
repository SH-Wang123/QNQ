//Decoupling gui with worker

package common

import (
	"sync"
)

var CLI_FLAG = false
var currentLockMap = make(map[int]*sync.WaitGroup)
var currentSNMap = make(map[int]string)
var runningFlagMap = make(map[int]bool)

var (
	localPartStartLock    = &sync.WaitGroup{}
	localBatchStartLock   = &sync.WaitGroup{}
	localSingleStartLock  = &sync.WaitGroup{}
	remoteSingleStartLock = &sync.WaitGroup{}
	cdpSnapshotStartLock  = &sync.WaitGroup{}
)

var (
	currentFileMap = make(map[string]string)
	cfLock         = &sync.RWMutex{}
)

var gwLock sync.RWMutex

var GWChannel = make(chan int)

func SendSignal2GWChannel(signal int) {
	gwLock.Lock()
	defer gwLock.Unlock()
	if CLI_FLAG {

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
	ret := currentFileMap[sn]
	if ret != "" {
		return ret
	}
	return "Starting..."
}

func GetCurrentSN(businessType int) string {
	return currentSNMap[businessType]
}

func SetCurrentSN(businessType int, SN string) {
	currentSNMap[businessType] = SN
}

func GetStartLock(businessType int) *sync.WaitGroup {
	return currentLockMap[businessType]
}

func SetRunningFlag(businessType int, runningFlag bool) {
	runningFlagMap[businessType] = runningFlag
}

func GetRunningFlag(businessType int) bool {
	return runningFlagMap[businessType]
}

func initLockMap() {
	currentLockMap[TYPE_LOCAL_BATCH] = localBatchStartLock
	currentLockMap[TYPE_PARTITION] = localPartStartLock
	currentLockMap[TYPE_LOCAL_SING] = localSingleStartLock
	currentLockMap[TYPE_REMOTE_SINGLE] = remoteSingleStartLock
	currentLockMap[TYPE_CDP_SNAPSHOT] = cdpSnapshotStartLock
}
