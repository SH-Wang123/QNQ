package test

import (
	"path/filepath"
	"testing"
	"window_handler/worker"
)

var batchFileSyncUT = "/batch_file_sync_ut"
var batchFileSyncUTLog = "[local single sync ut]: "

var localBatchSyncCase = []struct {
	startPath     string
	prefixName    string
	depth         int
	layerSize     int
	fileSize      int
	randomSize    bool
	randomContent bool
}{
	{"/1KB", "/1KB_", 5, 5, 1, false, true},
	{"/4KB", "/4KB_", 5, 5, 4, false, true},
	{"/8KB", "/8KB_", 5, 5, 8, false, true},
	{"/512KB", "/512KB_", 5, 5, 512, false, true},
	{"/1024KB", "/1024KB_", 5, 5, 1024, false, true},
}

func TestLocalBatchSyncCase(t *testing.T) {
	preBatchFileSyncCase()
	defer inhibitLog()()
	var sourcePath = utRoot + batchFileSyncUT + sourceRoot
	var targetPath = utRoot + batchFileSyncUT + targetRoot

	for _, testCase := range localBatchSyncCase {
		startAbsPath, _ := filepath.Abs(sourcePath + testCase.startPath)
		count := 0
		createFileTree(&(testCase.prefixName),
			startAbsPath,
			testCase.depth,
			testCase.layerSize,
			testCase.fileSize,
			testCase.randomSize,
			testCase.randomContent,
			&count)

		startNode := &worker.FileNode{
			IsDirectory:     true,
			HasChildren:     true,
			AbstractPath:    startAbsPath,
			AnchorPointPath: "",
			HeadFileNode:    nil,
			VarianceType:    worker.VARIANCE_ROOT,
		}
		worker.GetFileTree(startNode)
		worker.SyncBatchFileTree(startNode, targetPath+testCase.startPath)
		errorInfo := worker.GetBatchSyncError()
		if len(errorInfo) == 0 {
			t.Logf(batchFileSyncUTLog+"case {%v} ok!!!", testCase.prefixName)
		} else {
			t.Logf(batchFileSyncUTLog+"case {%v} error!!!", testCase.prefixName)
			for info := range errorInfo {
				t.Errorf(batchFileSyncUTLog+" %v", info)
			}
		}

	}
}

func preBatchFileSyncCase() {
	createUtPath(batchFileSyncUT)
}
