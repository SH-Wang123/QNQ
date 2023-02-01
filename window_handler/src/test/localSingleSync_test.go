package test

import (
	"path/filepath"
	"testing"
	"window_handler/worker"
)

var singleFileSyncUT = "singleFileSyncUT"

var singleFileSyncCase = []struct {
	fileName       string
	fileSize       int
	sourceFilePath string
	targetFilePath string
}{
	{"1KB", 1, utPath + singleFileSyncUT + "/source/", utPath + singleFileSyncUT + "/target/"},
	{"4KB", 4, utPath + singleFileSyncUT + "/source/", utPath + singleFileSyncUT + "/target/"},
	{"512KB", 512, utPath + singleFileSyncUT + "/source/", utPath + singleFileSyncUT + "/target/"},
	{"1024KB", 1024, utPath + singleFileSyncUT + "/source/", utPath + singleFileSyncUT + "/target/"},
	{"512MB", 1024 * 512, utPath + singleFileSyncUT + "/source/", utPath + singleFileSyncUT + "/target/"},
}

func TestSingleSyncNoCreateFile(t *testing.T) {
	preSingleFileSyncCase()
	for _, testCase := range singleFileSyncCase {
		sfAbsPath, _ := filepath.Abs(testCase.sourceFilePath + testCase.fileName)
		createFile(sfAbsPath, testCase.fileSize, true)
		sf, _ := worker.OpenFile(testCase.sourceFilePath, false)

		tfAbsPath, _ := filepath.Abs(testCase.targetFilePath + testCase.fileName)
		tf, _ := worker.OpenFile(tfAbsPath, true)

		defer worker.CloseFile(sf)
		defer worker.CloseFile(tf)

		caseWorker := worker.NewLocalSingleWorker(sf, tf)
		caseWorker.Execute()
		if !worker.CompareMd5(sf, tf) {
			t.Errorf("[local single sync ut]: error!!! source file : {%v}, target file : {%v}", sfAbsPath, tfAbsPath)
		} else {
			t.Logf("[local single sync ut]: case {%v} ok!!!", testCase.fileName)
		}
	}
}

func preSingleFileSyncCase() {
	utAbsPath, _ := filepath.Abs(utPath)
	worker.CreateDir(utAbsPath)
	worker.DeleteDir(utAbsPath + "/" + singleFileSyncUT)
	worker.CreateDir(utAbsPath + "/" + singleFileSyncUT)
	worker.CreateDir(utAbsPath + "/" + singleFileSyncUT + "/source")
	worker.CreateDir(utAbsPath + "/" + singleFileSyncUT + "/target")
}
