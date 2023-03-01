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

var gwLock sync.RWMutex

var GWChannel = make(chan int)

var LocalBatchPolicyRunningFlag = false

// SetLBSPRunning set local batch sync policy running
func SetLBSPRunning() {
	gwLock.Lock()
	defer gwLock.Unlock()
	GWChannel <- LOCAL_BATCH_POLICY_RUNNING
}

// SetLBSPStop set local batch sync policy stop
func SetLBSPStop() {
	gwLock.Lock()
	defer gwLock.Unlock()
	GWChannel <- LOCAL_BATCH_POLICY_STOP
}
