package test

import (
	"path/filepath"
	"testing"
	"time"
	"window_handler/config"
	"window_handler/worker"
)

var batchFileSyncUT = "/batch_file_sync_ut"
var batchFileSyncUTLog = "[local single sync ut]: "

type localBatchSyncCase struct {
	startPath     string
	prefixName    string
	depth         int
	layerSize     int
	fileSize      worker.CapacityUnit
	randomSize    bool
	randomContent bool
	count         int64
	bufferSize    worker.CapacityUnit
}

var localBatchSyncTestCase = []localBatchSyncCase{
	{"/1KB", "/1KB_", 5, 5, 1 * worker.KB, false, true, 0, 1 * worker.KB},
	{"/4KB", "/4KB_", 5, 5, 4 * worker.KB, false, true, 0, 1 * worker.KB},
	{"/8KB", "/8KB_", 5, 5, 8 * worker.KB, false, true, 0, 1 * worker.KB},
	{"/512KB", "/512KB_", 5, 5, 512 * worker.KB, false, true, 0, 1 * worker.KB},
	{"/1024KB", "/1024KB_", 5, 5, 1 * worker.MB, false, true, 0, 1 * worker.KB},
}

var periodicLocalBatchSyncTestCase = []localBatchSyncCase{
	{"/1KB", "/1KB_", 2, 2, 1 * worker.KB, false, true, 0, 1 * worker.KB},
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
	configCache := config.SystemConfigCache.Cache.LocalBatchSync
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
		configCache.TargetPath = targetPath + testCase.startPath
		configCache.SyncPolicy.PolicySwitch = true
		configCache.SyncPolicy.PeriodicSync.Cycle = time.Minute
		worker.StartPolicySync(time.Minute, &notEndFlag, true, false, true)
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
		count := int64(0)
		createFileTree(testCase.bufferSize,
			&(testCase.prefixName),
			startAbsPath,
			testCase.depth,
			testCase.layerSize,
			testCase.fileSize,
			testCase.randomSize,
			testCase.randomContent,
			&count)
		testCase.count = count
		overTime := time.Now()
		useTime := int64(overTime.Sub(startTime) / time.Second)
		if useTime <= 0 {
			useTime = 1
		}
		totalSize := (testCase.count * int64(testCase.fileSize)) / int64(worker.MB)
		t.Logf("create target file : [%v] over, total size : [%vKB], rate : [%vKB/s], time : %v", testCase.prefixName, totalSize/useTime, totalSize, overTime)
		t.Logf("-----------")
	}
}
