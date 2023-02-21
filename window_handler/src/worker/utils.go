package worker

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strconv"
)

func OpenFile(filePath string, createFile bool) (*os.File, error) {
	var f *os.File
	var err error
	if createFile {
		f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	} else {
		f, err = os.Open(filePath)
	}
	if err != nil {
		log.Printf("Open %v err : %v", filePath, err.Error())
		return nil, err
	}
	return f, nil
}

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Printf("close file err : %v", err.Error())
	}
}

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(path string) {
	exist, err := IsExist(path)
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

func DeleteFile(path string) error {
	exist, err := IsExist(path)
	if err != nil {
		log.Printf("get dir error : %v", err)
	}
	if exist {
		err = os.Remove(path)
		return err
	}
	return nil
}

func DeleteDir(path string) error {
	exist, err := IsExist(path)
	if err != nil {
		log.Printf("get dir error : %v", err)
	}
	if exist {
		err = os.RemoveAll(path)
		return err
	}
	return nil
}

func IsOpenDirError(err error, path string) bool {
	return err.Error() == "open "+path+": is a directory"
}

func GetFileMd5(f *os.File) *string {
	md5h := md5.New()
	io.Copy(md5h, f)
	md5Str := hex.EncodeToString(md5h.Sum(nil))
	return &md5Str
}

func CompareMd5(sf *os.File, tf *os.File) bool {
	sfMd5Ptr := GetFileMd5(sf)
	tfMd5Ptr := GetFileMd5(tf)
	return *sfMd5Ptr == *tfMd5Ptr
}

func CompareModifyTime(sf *os.File, tf *os.File) bool {
	sfInfo, err := sf.Stat()
	if err != nil {
		log.Printf("get file stat error : %v", err)
		return false
	}
	tfInfo, err := tf.Stat()
	if err != nil {
		log.Printf("get file stat error : %v", err)
		return false
	}
	return sfInfo.ModTime() == tfInfo.ModTime()
}

func GetSingleFileNode(path string) *FileNode {
	return &FileNode{
		IsDirectory:     false,
		HasChildren:     false,
		AbstractPath:    path,
		AnchorPointPath: "",
		HeadFileNode:    nil,
		VarianceType:    VARIANCE_ROOT,
	}
}

func uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func stringToUint64(s string) (uint64, error) {
	intNum, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint64(intNum), nil
}
