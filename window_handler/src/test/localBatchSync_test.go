package test

import (
	"path/filepath"
	"testing"
	"time"
	"window_handler/worker"
)

var batchFileSyncUT = "/batch_file_sync_ut"
var batchFileSyncUTLog = "[local single sync ut]: "

type localBatchSyncCase struct {
	startPath     string
	prefixName    string
	depth         int
	layerSize     int
	fileSize      int
	randomSize    bool
	randomContent bool
	count         int
}

var localBatchSyncTestCase = []localBatchSyncCase{
	{"/1KB", "/1KB_", 5, 5, 1, false, true, 0},
	{"/4KB", "/4KB_", 5, 5, 4, false, true, 0},
	{"/8KB", "/8KB_", 5, 5, 8, false, true, 0},
	{"/512KB", "/512KB_", 5, 5, 512, false, true, 0},
	{"/1024KB", "/1024KB_", 5, 5, 1024, false, true, 0},
}

var periodicLocalBatchSyncTestCase = []localBatchSyncCase{
	{"/1KB", "/1KB_", 2, 2, 1, false, true, 0},
}

var sourcePath = utRoot + batchFileSyncUT + sourceRoot
var targetPath = utRoot + batchFileSyncUT + targetRoot

func TestLocalBatchSyncCase(t *testing.T) {
	preBatchFileSyncCase()
	defer inhibitLog()()
	batchSyncCreateTargetFile(localBatchSyncTestCase, t)
	for _, testCase := range localBatchSyncTestCase {
		startAbsPath, _ := filepath.Abs(sourcePath + testCase.startPath)
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
			t.Logf(batchFileSyncUTLog+"case {%v} ok!!! count : {%v},  time: %v", testCase.prefixName, testCase.count, time.Now())
		} else {
			t.Logf(batchFileSyncUTLog+"case {%v} error!!! count : {%v},  time: %v", testCase.prefixName, testCase.count, time.Now())
			for info := range errorInfo {
				t.Errorf(batchFileSyncUTLog+" %v", info)
			}
		}

	}
}

func TestPeriodicLocalBatchSync(t *testing.T) {
	preBatchFileSyncCase()
	defer inhibitLog()()
	batchSyncCreateTargetFile(periodicLocalBatchSyncTestCase, t)
	for _, testCase := range periodicLocalBatchSyncTestCase {
		startAbsPath, _ := filepath.Abs(sourcePath + testCase.startPath)
		startNode := &worker.FileNode{
			IsDirectory:     true,
			HasChildren:     true,
			AbstractPath:    startAbsPath,
			AnchorPointPath: "",
			HeadFileNode:    nil,
			VarianceType:    worker.VARIANCE_ROOT,
		}
		worker.GetFileTree(startNode)
		notEndFlag := true
		go func() {
			ticker := time.NewTicker(30 * time.Minute)
			<-ticker.C
			notEndFlag = false
		}()
		worker.PeriodicLocalBatchSync(startNode, targetPath+testCase.startPath, 1*time.Minute, &notEndFlag)
		errorInfo := worker.GetBatchSyncError()
		if len(errorInfo) == 0 {
			t.Logf(batchFileSyncUTLog+"periodic policy case {%v} ok!!!  time: %v", testCase.prefixName, time.Now())
		} else {
			t.Logf(batchFileSyncUTLog+"periodic policy  {%v} error!!!  time: %v", testCase.prefixName, time.Now())
			for info := range errorInfo {
				t.Errorf(batchFileSyncUTLog+" %v", info)
			}
		}
	}
}

func preBatchFileSyncCase() {
	createUtPath(batchFileSyncUT)
}
func batchSyncCreateTargetFile(tc []localBatchSyncCase, t *testing.T) {
	for _, testCase := range tc {
		startTime := time.Now()
		t.Logf("create target file [%v] start, time : %v", testCase.prefixName, time.Now())
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
		testCase.count = count
		overTime := time.Now()
		userTime := overTime.Second() - startTime.Second()
		if userTime <= 0 {
			userTime = 1
		}
		totalSize := testCase.count * testCase.fileSize
		t.Logf("create target file : [%v] over, total size : [%vKB], rate : [%vKB/s], time : %v", testCase.prefixName, totalSize/userTime, totalSize, overTime)
		t.Logf("-----------")
	}
}
