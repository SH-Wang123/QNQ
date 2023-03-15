//Decoupling gui with worker

package common

import "sync"

const (
	LOCAL_BATCH_POLICY_RUNNING = iota
	LOCAL_BATCH_POLICY_STOP
	LOCAL_SINGLE_POLICY_RUNNING
	LOCAL_SINGLE_POLICY_STOP
	TEST_DISK_SPEED_START
	TEST_DISK_SPEED_OVER
)

var (
	CurrentLocalPartSN   string
	CurrentLocalPartFile string = "Not running"
	LocalPartStartLock          = &sync.WaitGroup{}
)

var (
	CurrentLocalBatchSN   string
	CurrentLocalBatchFile string = "Not running"
	LocalBatchStartLock          = &sync.WaitGroup{}
)

var gwLock sync.RWMutex

var GWChannel = make(chan int)

var LocalBatchPolicyRunningFlag = false
var LocalSinglePolicyRunningFlag = false

func SendSignal2GWChannel(signal int) {
	gwLock.Lock()
	defer gwLock.Unlock()
	GWChannel <- signal
}
