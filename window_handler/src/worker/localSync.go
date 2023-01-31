package worker

import (
	"io"
	"log"
	"os"
	"window_handler/common"
	"window_handler/config"
)

func FileNode2TreeMap(data *map[string][]string) {
	datat := *data
	datat[""] = append(datat[""], LocalBSFileNode.AbstractPath)
	getFileTree(LocalBSFileNode, data)
}

func getFileTree(node *FileNode, data *map[string][]string) {
	datat := *data
	if !node.IsDirectory {
		return
	}
	for _, child := range node.ChildrenNodeList {
		var key string
		if node.VarianceType == VARIANCE_ROOT {
			key = node.AbstractPath
		} else {
			key = node.AnchorPointPath
		}
		datat[key] = append(datat[key], child.AnchorPointPath)
		getFileTree(child, data)
	}
}

func MarkFileTree(node *FileNode, rootPath string) {
	currentPath := rootPath + node.AnchorPointPath
	exist, _ := isExist(currentPath)
	if !exist {
		if node.VarianceType == VARIANCE_ROOT {
			node.VarianceType = VARIANCE_ROOT | VARIANCE_ADD
		} else {
			node.VarianceType = VARIANCE_ADD
		}
	} else {
		//TODO check md5,timestamp,lastTime
	}
	if node.IsDirectory {
		for _, child := range node.ChildrenNodeList {
			MarkFileTree(child, rootPath+child.AnchorPointPath)
		}
	}

}

// SyncBatchFileTree Crete folder
func SyncBatchFileTree(node FileNode, startPath string) {
	for _, child := range node.ChildrenNodeList {
		absPath := startPath + fileSeparator + child.AnchorPointPath
		if !child.IsDirectory {
			tf, err := OpenFile(absPath, true)
			if err == nil {
				//common.GetCoroutinesPool().Submit(worker.Execute())
				sf, err := OpenFile(child.AbstractPath, false)
				if err == nil {
					worker := NewLocalSingleWorker(sf, tf)
					common.GetCoroutinesPool().Submit(worker.Execute)
				} else {
					CloseFile(tf)
				}
			}
		} else {
			CreateDir(absPath)
			SyncBatchFileTree(*child, absPath)
		}
	}
}

func LocalSyncSingleFile() bool {
	tempTarget := ""

	sf, err := OpenFile(config.SystemConfigCache.Value().LocalSingleSync.SourcePath, false)
	if err != nil {
		return false
	}

	tf, err := OpenFile(config.SystemConfigCache.Value().LocalSingleSync.TargetPath, true)

	if err != nil {
		if IsOpenDirError(err, config.SystemConfigCache.Value().LocalSingleSync.TargetPath) {
			sfInfo, _ := sf.Stat()
			tempTarget = config.SystemConfigCache.Value().LocalSingleSync.TargetPath + "/" + sfInfo.Name()
			tf, err = OpenFile(tempTarget, true)
			if err != nil {
				return false
			}
		} else {
			return false
		}
	}

	worker := NewLocalSingleWorker(sf, tf)
	common.GetCoroutinesPool().Submit(worker.Execute)

	return true
}

func NewLocalSingleWorker(sourceFile *os.File, targetFile *os.File) *common.QWorker {
	return &common.QWorker{
		Sub:             nil,
		ExecuteFunc:     RemoteSyncSingleFile,
		DeconstructFunc: closeFile,
		PrivateFile:     sourceFile,
		TargetFile:      targetFile,
	}
}

func RemoteSyncSingleFile(msg interface{}, q *common.QWorker) {
	buf := make([]byte, 4096)
	defer CloseFile(q.TargetFile)
	defer CloseFile(q.PrivateFile)
	for {
		n, err := q.PrivateFile.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			break
		}
		_, err = q.TargetFile.Write(buf[:n])
		if err != nil {
			break
		}
	}

	defer singleFileDone()
}

func closeFile(w *common.QWorker) {

}

func singleFileDone() {
	DoneFileNum++
	log.Printf("Done num : %f", DoneFileNum)
}

func GetLocalBatchProgress() float64 {
	return DoneFileNum / TotalFileNum
}
