package test

import (
	"fmt"
	"log"
	"os"
)

// fileSize : KB
func createFile(filePath string, fileSize int, randomContent bool) bool {
	exist, err := isExist(filePath)
	if exist {
		return false
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return false
	}
	defer f.Close()
	content := make([]byte, 1024)
	if randomContent {

	}
	for count := 1; count <= fileSize; count++ {
		f.Write(content)
	}
	return true
}

func createFileTree(startPath string, depth int, layerSize int, fileSize int, randomSize bool, randomContent bool, count *int) {
	createDir(startPath)
	if depth <= 0 {
		return
	}
	//at: auto test
	filePrefix := "/crate_file_tree_at_"

	//create file
	for fileIndex := 1; fileIndex <= layerSize; fileIndex++ {
		tempPath1 := startPath + filePrefix + fmt.Sprintf("%d", *count)
		createFile(tempPath1, fileSize, randomContent)
		*count++
	}
	//create folder
	for folderIndex := 1; folderIndex <= layerSize; folderIndex++ {
		tempPath2 := startPath + filePrefix + fmt.Sprintf("%d", *count)
		createDir(tempPath2)
		createFileTree(tempPath2, depth-1, layerSize, fileSize, randomSize, randomContent, count)
		*count++
	}
}

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createDir(path string) {
	exist, err := isExist(path)
	if err != nil {
		log.Printf("get dir error : %v", err)
	}
	if !exist {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("create dir error : %v", err)
		}
	}
}
