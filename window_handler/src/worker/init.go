package worker

import (
	"runtime"
	"window_handler/common"
)

var CapacityStrMap = make(map[string]CapacityUnit)
var osName string
var linuxOSName = "linux"
var windowsOSName = "windows"
var macOSName = "mac"

func init() {
	CapacityStrMap["Byte"] = Byte
	CapacityStrMap["KB"] = KB
	CapacityStrMap["MB"] = MB
	CapacityStrMap["GB"] = GB
	CapacityStrMap["TB"] = TB
	CapacityStrMap["PB"] = PB
	gcFriend()
	LoadWorkerFactory()
	GetPartitionsInfo()
	osName = runtime.GOOS
}

// gcHelper 定时清理无用数据的引用，GC好帮手
func gcFriend() {
	//TODO 清理totalSizeMap和doneMap（根据是否完成去清理）
}

func LoadWorkerFactory() {
	common.WorkerFactoryMap[common.GetQMQTaskPre(common.TYPE_REMOTE_SINGLE)] = NewRemoteSyncReceiver
}
