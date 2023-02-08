//Decoupling gui with worker

package common

const (
	LOCAL_BATCH_POLICY_RUNNING = iota
	LOCAL_BATCH_POLICY_STOP
	LOCAL_SINGLE_POLICY_RUNNING
	LOCAL_SINGLE_POLICY_STOP
)

var GWChannel = make(chan int)

// SetLBSPRunning set local batch sync policy running
func SetLBSPRunning() {
	GWChannel <- LOCAL_BATCH_POLICY_RUNNING
}

// SetLBSPStop set local batch sync policy stop
func SetLBSPStop() {
	GWChannel <- LOCAL_BATCH_POLICY_STOP
}
