package test

import (
	"path/filepath"
	"window_handler/worker"
)

func preCase() {
	utAbsPath, _ := filepath.Abs(utPath)
	worker.CreateDir(utAbsPath)
	worker.CreateDir(utAbsPath + "/" + singleFileSyncUT)
	worker.CreateDir(utAbsPath + "/" + singleFileSyncUT + "/source")
	worker.CreateDir(utAbsPath + "/" + singleFileSyncUT + "/target")
}
