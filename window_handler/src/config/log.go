package config

import (
	"bufio"
	"os"
	"window_handler/common"
)

const LOG_PATH = "OLog.csv"

var logTypeMap = make(map[int]string)

func initOLog() {
	initLogTypeMap()
	isExist, _ := common.IsExist(LOG_PATH)
	if isExist {
		return
	}
	f, _ := common.OpenFile(LOG_PATH, true)
	AddOLog("Name,Start Time,Over Time,Result,Target Path,Source Path")
	defer common.CloseFile(f)
}

func LoadOLog() []string {
	ret := make([]string, 0)
	f, err := os.Open(LOG_PATH)
	if err != nil {
		return ret
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {

		}
	}(f)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	return ret
}

func DeleteOLog(index int) bool {

	return true
}

func AddOLog(log string) bool {
	if log == "" {
		return true
	}
	f, _ := os.OpenFile(LOG_PATH, os.O_RDWR|os.O_APPEND, 0777)
	defer f.Close()
	f.Write([]byte(log + "\n"))
	return true
}

func GetOLogType(busType int) string {
	return logTypeMap[busType]
}

func initLogTypeMap() {
	logTypeMap[common.TYPE_PARTITION] = "Partition Sync"
	logTypeMap[common.TYPE_LOCAL_SING] = "Local Single Sync"
	logTypeMap[common.TYPE_LOCAL_BATCH] = "Local Batch Sync"
	logTypeMap[common.TYPE_CDP_SNAPSHOT] = "Create CDP Snapshot"
	logTypeMap[common.TYPE_REMOTE_BATCH] = "Remote Batch Sync"
	logTypeMap[common.TYPE_REMOTE_SINGLE] = "Remote Single Sync"
}
