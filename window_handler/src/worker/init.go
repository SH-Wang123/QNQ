package worker

import (
	"log"
	"runtime"
	"window_handler/common"
	"window_handler/network"
)

var CapacityStrMap = make(map[string]CapacityUnit)
var osName string
var linuxOSName = "linux"
var windowsOSName = "windows"
var macOSName = "mac"

var gwChannelRegisterF = make(map[int]func(), 16)

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
	initRegisterGWFunc()
	go watchGWChannel()
}

// gcFriend 定时清理无用数据的引用，GC好帮手
func gcFriend() {
	//TODO 清理totalSizeMap和doneMap（根据是否完成去清理）
}

func LoadWorkerFactory() {
	common.WorkerFactoryMap[common.GetQMQTaskPre(common.TYPE_REMOTE_SINGLE)] = NewRemoteSyncReceiver
}

func watchGWChannel() {
	for {
		select {
		case signal := <-common.GWChannel:
			f := gwChannelRegisterF[signal]
			if f == nil {
				log.Printf("!!!!!!!!!!!!!!!!!!has a signal doesn't register, num : %v, gw", signal)
				continue
			}
			f()
		}
	}
}

// initRegisterGWFunc 注册GW通道响应函数
func initRegisterGWFunc() {
	registerGWFunc(common.SIGNAL_AUTH_PASS, authPassHandler)
	registerGWFunc(common.SIGNAL_AUTH_NO_PASS, authNoPassHandler)
}

func registerGWFunc(signal int, f func()) {
	gwChannelRegisterF[signal] = f
}

func authPassHandler() {
	network.AuthLock.Unlock()
	network.AuthFlag = true
}

func authNoPassHandler() {
	network.AuthLock.Unlock()
	network.AuthFlag = false
}
