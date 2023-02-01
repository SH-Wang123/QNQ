package worker

import (
	"log"
	"os"
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
	err := os.Remove(path)
	return err
}

func DeleteDIr(path string) error {
	err := os.RemoveAll(path)
	return err
}

func IsOpenDirError(err error, path string) bool {
	return err.Error() == "open "+path+": is a directory"
}
