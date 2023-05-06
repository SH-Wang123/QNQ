package test

import (
	"path/filepath"
	"testing"
	"window_handler/common"
	"window_handler/worker"
)

var singleFileSyncUTRoot = "/single_file_sync_ut"

var singleFileSyncCase = []struct {
	fileName      string
	fileSize      worker.CapacityUnit
	randomContent bool
	bufferSize    worker.CapacityUnit
}{
	{"1KB", 1 * worker.KB, true, 1 * worker.MB},
	{"4KB", 4 * worker.KB, true, 1 * worker.MB},
	{"512KB", 512 * worker.KB, true, 1 * worker.MB},
	{"1024KB", 1 * worker.MB, true, 1 * worker.MB},
	{"512MB", 512 * worker.MB, true, 1 * worker.MB},
}

func TestSingleSync(t *testing.T) {
	preSingleFileSyncCase()
	defer inhibitLog()()
	var sourcePath = utRoot + singleFileSyncUTRoot + sourceRoot
	var targetPath = utRoot + singleFileSyncUTRoot + targetRoot
	for _, testCase := range singleFileSyncCase {
		sfAbsPath, _ := filepath.Abs(sourcePath + "/" + testCase.fileName)
		worker.CreateFile(testCase.bufferSize, sfAbsPath, testCase.fileSize, testCase.randomContent)
		sf, _ := common.OpenFile(sfAbsPath, false)

		tfAbsPath, _ := filepath.Abs(targetPath + "/" + testCase.fileName)
		tf, _ := common.OpenFile(tfAbsPath, true)
		sn := common.GetTaskCount()
		caseWorker := worker.NewLocalSingleWorker(sf, tf, sn, false)
		caseWorker.Execute()
		if !worker.CompareMd5(sf, tf) {
			t.Errorf("[local single sync ut]: ERROR!!! source file : {%v}, target file : {%v}", sfAbsPath, tfAbsPath)
		} else {
			t.Logf("[local single sync ut]: case {%v} ok!!!", testCase.fileName)
		}
	}
}

func preSingleFileSyncCase() {
	createUtPath(singleFileSyncUTRoot)
}
