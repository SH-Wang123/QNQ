package test

import (
	"path/filepath"
	"testing"
	"window_handler/worker"
)

var singleFileSyncUTRoot = "/single_file_sync_ut"

var singleFileSyncCase = []struct {
	fileName      string
	fileSize      int
	randomContent bool
}{
	{"1KB", 1, true},
	{"4KB", 4, true},
	{"512KB", 512, true},
	{"1024KB", 1024, true},
	{"512MB", 1024 * 512, true},
}

func TestSingleSync(t *testing.T) {
	preSingleFileSyncCase()
	defer inhibitLog()()
	var sourcePath = utRoot + singleFileSyncUTRoot + sourceRoot
	var targetPath = utRoot + singleFileSyncUTRoot + targetRoot
	for _, testCase := range singleFileSyncCase {
		sfAbsPath, _ := filepath.Abs(sourcePath + "/" + testCase.fileName)
		createFile(sfAbsPath, testCase.fileSize, testCase.randomContent)
		sf, _ := worker.OpenFile(sfAbsPath, false)

		tfAbsPath, _ := filepath.Abs(targetPath + "/" + testCase.fileName)
		tf, _ := worker.OpenFile(tfAbsPath, true)
		caseWorker := worker.NewLocalSingleWorker(sf, tf)
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
