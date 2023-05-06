package test

import (
	"path/filepath"
	"testing"
	"time"
	"window_handler/common"
	"window_handler/config"
	"window_handler/worker"
)

var batchFileSyncUT = "/batch_file_sync_ut"
var batchFileSyncUTLog = "[local batch sync ut]: "

var partitionSyncUT = "/partition_sync_ut"
var partitionSyncUTLog = "[partition sync ut]: "

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

type periodicCase struct {
	cycle time.Duration
	rate  int
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

var periodicTestCase = []periodicCase{
	{time.Minute, 1},
	{time.Minute, 5},
	{time.Minute, 30},
}

var sourcePath = utRoot + batchFileSyncUT + sourceRoot
var targetPath = utRoot + batchFileSyncUT + targetRoot

func TestLocalBatchSyncCase(t *testing.T) {
	preBatchFileSyncCase()
	defer afterBatchFileSyncCase()
	batchSyncCreateTargetFile(localBatchSyncTestCase, t)
	for _, testCase := range localBatchSyncTestCase {
		startAbsPath, _ := filepath.Abs(sourcePath + testCase.startPath)
		targetAbsPath := targetPath + testCase.startPath
		setBatchConfig(startAbsPath, targetAbsPath)
		startTime := common.GetNowTimeStr()
		worker.LocalBatchSyncSingleTime(false)
		errorInfo := worker.GetBatchSyncError(common.GetCurrentSN(common.TYPE_LOCAL_BATCH))
		if len(errorInfo) == 0 {
			t.Logf(batchFileSyncUTLog+"case {%v} ok!!! start time: %v, over time: %v", testCase.prefixName, startTime, common.GetNowTimeStr())
		} else {
			t.Logf(batchFileSyncUTLog+"case {%v} error!!! start time: %v, over time: %v", testCase.prefixName, startTime, common.GetNowTimeStr())
			for info := range errorInfo {
				t.Errorf(batchFileSyncUTLog+" %v", info)
			}
		}

	}
}

func TestPeriodicLocalBatchSync(t *testing.T) {
	preBatchFileSyncCase()
	defer afterBatchFileSyncCase()
	batchSyncCreateTargetFile(periodicLocalBatchSyncTestCase, t)
	configCache := config.SystemConfigCache.Cache.LocalBatchSync
	for _, testCase := range periodicLocalBatchSyncTestCase {
		startAbs, _ := filepath.Abs(sourcePath + testCase.startPath)
		targetAbs := targetPath + testCase.startPath
		setBatchConfig(startAbs, targetAbs)
		configCache.SyncPolicy.PolicySwitch = true
		configCache.SyncPolicy.PeriodicSync.Enable = true
		for _, periodic := range periodicTestCase {
			notEndFlag := true
			configCache.SyncPolicy.PeriodicSync.Cycle = periodic.cycle
			configCache.SyncPolicy.PeriodicSync.Rate = periodic.rate
			worker.StartPolicySync(time.Minute, &notEndFlag, true, false, false, true)
			errorInfo := worker.GetBatchSyncError(common.GetCurrentSN(common.TYPE_PARTITION))
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
}

func preBatchFileSyncCase() {
	commonPreTest(batchFileSyncUT)
}

func afterBatchFileSyncCase() {
	commonAfterTest(batchFileSyncUT)
}

func batchSyncCreateTargetFile(tc []localBatchSyncCase, t *testing.T) {
	for _, testCase := range tc {
		startTime := time.Now()
		t.Logf("create target file [%v] start, time : %v", testCase.prefixName, common.GetNowTimeStr())
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
		t.Logf("create target file : [%v] over, total size : [%vKB], rate : [%vMB/s], time : %v", testCase.prefixName, totalSize/useTime, totalSize, common.GetNowTimeStr())
		t.Logf("-----------")
	}
}

func setBatchConfig(src, target string) {
	config.SystemConfigCache.Cache.LocalBatchSync.SourcePath = src
	config.SystemConfigCache.Cache.LocalBatchSync.TargetPath = target
}
