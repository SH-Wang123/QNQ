package worker

import (
	"log"
	"runtime"
	"window_handler/common"
	"window_handler/config"
)

var CapacityStrMap = make(map[string]CapacityUnit)
var osName string
var linuxOSName = "linux"
var windowsOSName = "windows"

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

// Deprecated: 没有意义，太浪费内存，无需缓存整棵文件树
func InitFileNode(initAll bool, async bool) {
	LocalBSFileNode = GetNilNode(config.SystemConfigCache.Value().LocalBatchSync.SourcePath)
	if async {
		go GetFileTree(LocalBSFileNode, true)
	} else {
		GetFileTree(LocalBSFileNode, true)
	}
}

// gcHelper 定时清理无用数据的引用，GC好帮手
func gcFriend() {
	//TODO 清理totalSizeMap和doneMap（根据是否完成去清理）
}

// Deprecated: 顽固的面向对象思想，占用内存过多
func GetFileTree(fNode *FileNode, isRecurrence bool) {
	f, _ := common.OpenFile(fNode.AbstractPath, false)
	allChild, err := f.Readdir(-1)
	if err != nil {
		fNode.IsDirectory = false
		log.Printf("open dir error, path : %v, error : %v", fNode.AbstractPath, err)
		return
	}
	if len(allChild) > 0 {
		fNode.HasChildren = true
		for _, child := range allChild {
			childFileNode := FileNode{
				HeadFileNode:    fNode,
				HasChildren:     false,
				AbstractPath:    fNode.AbstractPath + fileSeparator + child.Name(),
				IsDirectory:     child.IsDir(),
				AnchorPointPath: child.Name(),
				VarianceType:    VARIANCE_EDIT,
			}
			fNode.ChildrenNodeList = append(fNode.ChildrenNodeList, &childFileNode)
			if isRecurrence {
				if child.IsDir() {
					GetFileTree(&childFileNode, isRecurrence)
				}
			}
		}
	} else {
		fNode.IsDirectory = false
	}
	defer common.CloseFile(f)
}

func LoadWorkerFactory() {
	common.WorkerFactoryMap[common.RemoteSingleSyncType] = NewRemoteSyncReceiver
}
