package worker

import (
	"log"
	"window_handler/common"
	"window_handler/config"
)

var TotalFileNum = 0.0
var DoneFileNum = 0.0
var CapacityStrMap = make(map[string]CapacityUnit)

func init() {
	CapacityStrMap["Byte"] = Byte
	CapacityStrMap["KB"] = KB
	CapacityStrMap["MB"] = MB
	CapacityStrMap["GB"] = GB
	CapacityStrMap["TB"] = TB
	CapacityStrMap["PB"] = PB
}

func InitFileNode(initAll bool, async bool) {
	LocalBSFileNode = &FileNode{
		IsDirectory:     true,
		HasChildren:     true,
		AbstractPath:    config.SystemConfigCache.Value().LocalBatchSync.SourcePath,
		AnchorPointPath: "",
		HeadFileNode:    nil,
		VarianceType:    VARIANCE_ROOT,
	}
	if async {
		go GetFileTree(LocalBSFileNode)
	} else {
		GetFileTree(LocalBSFileNode)
	}
}

// GetFilePath BFS
func GetFileTree(fNode *FileNode) {
	f, _ := OpenFile(fNode.AbstractPath, false)
	allChild, err := f.Readdir(-1)
	if err != nil {
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
			if child.IsDir() {
				GetFileTree(&childFileNode)
			} else {
				TotalFileNum++
			}
		}
	}
	defer CloseFile(f)
}

func LoadWorkerFactory() {
	common.WorkerFactoryMap[common.RemoteSingleSyncType] = NewRemoteSyncReceiver
}
