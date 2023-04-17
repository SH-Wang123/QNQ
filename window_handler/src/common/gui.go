//Decoupling gui with worker

package common

import (
	"sync"
)

var CLI_FLAG = false

var (
	currentLocalPartSN string
	localPartStartLock = &sync.WaitGroup{}
)

var (
	currentLocalBatchSN string
	localBatchStartLock = &sync.WaitGroup{}
)

var (
	currentLocalSingleSN string
	localSingleStartLock = &sync.WaitGroup{}
)

var (
	currentRemoteSingleSN string
	remoteSingleStartLock = &sync.WaitGroup{}
)

var (
	currentCDPSnapshotSN string
	cdpSnapshotStartLock = &sync.WaitGroup{}
)

var (
	currentFileMap = make(map[string]string)
	cfLock         = &sync.RWMutex{}
)

var gwLock sync.RWMutex

var GWChannel = make(chan int)

var (
	localBatchRunningFlag   = false
	localSingleRunningFlag  = false
	remoteSingleRunningFlag = false
	partitionRunningFlag    = false
	cdpSnapshotRunningFlag  = false
)

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
	switch businessType {
	case TYPE_LOCAL_BATCH:
		return currentLocalBatchSN
	case TYPE_PARTITION:
		return currentLocalPartSN
	case TYPE_LOCAL_SING:
		return currentLocalSingleSN
	case TYPE_REMOTE_SINGLE:
		return currentRemoteSingleSN
	case TYPE_CDP_SNAPSHOT:
		return currentCDPSnapshotSN
	default:
		return ""
	}
}

func SetCurrentSN(businessType int, SN string) {
	switch businessType {
	case TYPE_LOCAL_BATCH:
		currentLocalBatchSN = SN
	case TYPE_PARTITION:
		currentLocalPartSN = SN
	case TYPE_LOCAL_SING:
		currentLocalSingleSN = SN
	case TYPE_REMOTE_SINGLE:
		currentRemoteSingleSN = SN
	case TYPE_CDP_SNAPSHOT:
		currentCDPSnapshotSN = SN
	}
}

func GetStartLock(businessType int) *sync.WaitGroup {
	switch businessType {
	case TYPE_LOCAL_BATCH:
		return localBatchStartLock
	case TYPE_PARTITION:
		return localPartStartLock
	case TYPE_LOCAL_SING:
		return localSingleStartLock
	case TYPE_REMOTE_SINGLE:
		return remoteSingleStartLock
	case TYPE_CDP_SNAPSHOT:
		return cdpSnapshotStartLock
	default:
		return &sync.WaitGroup{}
	}
}

func SetRunningFlag(businessType int, runningFlag bool) {
	switch businessType {
	case TYPE_LOCAL_BATCH:
		localBatchRunningFlag = runningFlag
	case TYPE_LOCAL_SING:
		localSingleRunningFlag = runningFlag
	case TYPE_PARTITION:
		partitionRunningFlag = runningFlag
	case TYPE_REMOTE_SINGLE:
		remoteSingleRunningFlag = runningFlag
	case TYPE_CDP_SNAPSHOT:
		cdpSnapshotRunningFlag = runningFlag
	}
}

func GetRunningFlag(businessType int) bool {
	switch businessType {
	case TYPE_LOCAL_BATCH:
		return localBatchRunningFlag
	case TYPE_LOCAL_SING:
		return localSingleRunningFlag
	case TYPE_PARTITION:
		return partitionRunningFlag
	case TYPE_REMOTE_SINGLE:
		return remoteSingleRunningFlag
	case TYPE_CDP_SNAPSHOT:
		return cdpSnapshotRunningFlag
	default:
		return true
	}
}
