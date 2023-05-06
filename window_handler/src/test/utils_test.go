package test

import (
	"os"
	"testing"
	"window_handler/worker"
)

func TestCreateFile(t *testing.T) {
	//if !worker.CreateFile()reateFile("E:/source/tree_ttt", 1024*512, true) {
	//	t.Error("error")
	//}
}

func TestCreateFileTree(t *testing.T) {
	//count := 0
	//pre := "/8KB_"
	//createFileTree(&pre, "E:/source/tree_ttt", 5, 5, 8, false, true, &count)
}

func TestRandomContent(t *testing.T) {
	//t.Logf("Random result : %s", randomPalindrome(1024))
}

func TestCompareMd5(t *testing.T) {
	sf, _ := os.Open("")
	tf, _ := os.Open("")
	t.Logf("%v", worker.CompareMd5(sf, tf))
	//worker.GetBatchSyncError()
}
