package test

import (
	"fmt"
	"log"
	_ "log"
	"os"
	"path/filepath"
	"window_handler/worker"
)

func createFileTree(bufferSize worker.CapacityUnit, filePrefix *string, startPath string, depth int, layerSize int, fileSize worker.CapacityUnit, randomSize bool, randomContent bool, count *int64) {
	worker.CreateDir(startPath)
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
		worker.CreateDir(tempPath2)
		createFileTree(bufferSize, filePrefix, tempPath2, depth-1, layerSize, fileSize, randomSize, randomContent, count)
		*count++
	}
}

func createUtPath(utPath string) {
	utAbsPath, _ := filepath.Abs(utRoot)
	worker.CreateDir(utAbsPath)
	worker.DeleteDir(utAbsPath + utPath)
	worker.CreateDir(utAbsPath + utPath)
	worker.CreateDir(utAbsPath + utPath + sourceRoot)
	worker.CreateDir(utAbsPath + utPath + targetRoot)
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
