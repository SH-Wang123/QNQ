package common

import (
	"log"
	"os"
)

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

func CreateDir(path string, perm *os.FileMode) {
	exist, err := IsExist(path)
	if err != nil {
		log.Printf("get dir exist error : %v", err)
	}
	if !exist {
		err = os.Mkdir(path, *perm)
		if err != nil {
			log.Printf("create dir error : %v", err)
		}
	} else {
		err := os.Chmod(path, *perm)
		if err != nil {
			log.Printf("Set folder perm err: %v", err)
		}
	}
}

func DeleteFileOrDir(path string) {
	exist, err := IsExist(path)
	f, _ := OpenFile(path, false)
	fChild, _ := f.Readdir(-1)
	CloseFile(f)
	if err != nil {
		log.Printf("get file error : %v", err)
	}
	if exist {
		if len(fChild) == 0 {
			err = os.Remove(path)
		} else {
			err = os.RemoveAll(path)
		}
	}
	if err != nil {
		log.Printf("delte file error : %v", err)
	}
}

func IsOpenDirError(err error, path string) bool {
	return err.Error() == "open "+path+": is a directory"
}

func OpenFile(filePath string, createFile bool) (*os.File, error) {
	var f *os.File
	var err error
	if createFile {
		f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	} else {
		f, err = os.OpenFile(filePath, os.O_RDONLY, 0777)
	}
	if err != nil {
		log.Printf("Open %v err : %v", filePath, err.Error())
		return nil, err
	}
	return f, nil
}

func CloseFile(fs ...*os.File) {
	for _, f := range fs {
		err := f.Close()
		if err != nil {
			log.Printf("close file err : %v", err.Error())
		}
	}
}

func OpenDir(filePath string) (*os.File, error) {
	f, err := os.Open(filePath)
	return f, err
}
