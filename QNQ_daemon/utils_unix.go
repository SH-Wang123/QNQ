//go:build unix || linux

package QNQ

import (
	"QNQ/common"
	"QNQ/model"
	"log/slog"
	"os"
	"syscall"
)

var sparator = "/"

func syncFile(source string, target string, ctx context.Context) *model.TaskResult {
	res := &model.TaskResult{
		Code: common.OK_CODE,
	}
	src, err := os.Open(source)
	if err != nil {
		slog.Error("source file open err : " + err.Error())
		return res
	}
	tar, err := os.Create(target)
	if err != nil {
		slog.Error("source file open err : " + err.Error())
		return res
	}

	defer func(src *os.File, tar *os.File) {
		err := tar.Sync()
		if err != nil {

		}
		err = src.Close()
		if err != nil {

		}
		err = tar.Close()
		if err != nil {

		}
	}(src, tar)

	srcInfo, _ := src.Stat()
	if srcInfo.IsDir() {
		children, _ := src.Readdir(-1)
		for _, _ = range children {

		}
	} else {
		srcFd := src.Fd()
		tarFD := tar.Fd()
		n, _ := syscall.Sendfile(int(tarFD), int(srcFd), nil, int(srcInfo.Size()))
		if int64(n) != srcInfo.Size() {
			return res
		}
	}

	return res
}

func GetAllDiskInfo() {

}
