package test

import (
	"fmt"
	_ "log"
	"math/rand"
	"os"
	"time"
	"window_handler/worker"
)

// fileSize : KB
func createFile(filePath string, fileSize int, randomContent bool) bool {
	exist, err := worker.IsExist(filePath)
	if exist {
		return false
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return false
	}
	defer f.Close()
	content := make([]byte, 1024)
	for count := 1; count <= fileSize; count++ {
		if randomContent {
			content = randomPalindrome(1024)
		}
		f.Write(content)
	}
	return true
}

func createFileTree(startPath string, depth int, layerSize int, fileSize int, randomSize bool, randomContent bool, count *int) {
	worker.CreateDir(startPath)
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
		worker.CreateDir(tempPath2)
		createFileTree(tempPath2, depth-1, layerSize, fileSize, randomSize, randomContent, count)
		*count++
	}
}

// randomPalindrome size : byte
func randomPalindrome(size int) []byte {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	bytes := make([]byte, size)
	for i := 0; i < (size+1)/2; i++ {
		r := byte(rng.Intn(0x1000)) //random rune up to '\u0999'
		bytes[i] = r
		bytes[size-1-i] = r
	}
	return bytes
}
