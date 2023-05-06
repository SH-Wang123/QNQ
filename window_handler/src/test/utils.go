package test

import (
	"fmt"
	"io/fs"
	"log"
	_ "log"
	"os"
	"path/filepath"
	"window_handler/common"
	"window_handler/config"
	"window_handler/worker"
)

var tMode fs.FileMode
var tmpConfCache = config.SystemConfigCache

func init() {
	srcPath, _ := filepath.Abs("")
	srcDic, _ := common.OpenDir(srcPath)
	srcInfo, _ := srcDic.Stat()
	tMode = srcInfo.Mode()
}

func createFileTree(bufferSize worker.CapacityUnit, filePrefix *string, startPath string, depth int, layerSize int, fileSize worker.CapacityUnit, randomSize bool, randomContent bool, count *int64) {
	common.CreateDir(startPath, &tMode)
	if depth <= 0 {
		return
	}

	//create file
	for fileIndex := 1; fileIndex <= layerSize; fileIndex++ {
		tempPath1 := startPath + *filePrefix + fmt.Sprintf("%d", *count)
		worker.CreateFile(bufferSize, tempPath1, fileSize, randomContent)
		*count++
	}
	//create folder
	for folderIndex := 1; folderIndex <= layerSize; folderIndex++ {
		tempPath2 := startPath + *filePrefix + fmt.Sprintf("%d", *count)
		common.CreateDir(tempPath2, &tMode)
		createFileTree(bufferSize, filePrefix, tempPath2, depth-1, layerSize, fileSize, randomSize, randomContent, count)
		*count++
	}
}

func commonPreTest(utPath string) {
	createUtPath(utPath)
	config.SystemConfigCache.Cache.SystemSetting.EnableOLog = false
}

func commonAfterTest(utPath string) {
	inhibitLog()()
	delUtPath(utPath)
	reloadConfCache()
}

func createUtPath(utPath string) {
	utAbsPath, _ := filepath.Abs(utRoot)
	common.CreateDir(utAbsPath, &tMode)
	common.DeleteFileOrDir(utAbsPath + utPath)
	common.CreateDir(utAbsPath+utPath, &tMode)
	common.CreateDir(utAbsPath+utPath+sourceRoot, &tMode)
	common.CreateDir(utAbsPath+utPath+targetRoot, &tMode)
}

func delUtPath(utPath string) {
	utAbsPath, _ := filepath.Abs(utRoot)
	common.DeleteFileOrDir(utAbsPath)
}

func reloadConfCache() {
	config.SystemConfigCache = tmpConfCache
}

func inhibitLog() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
