package config

import (
	"bufio"
	"os"
	"window_handler/common"
)

const oLogPath = "OLog.csv"
const timePointPath = "QTP.csv"

var logTypeMap = make(map[int]string)

func initOLog() {
	initLogTypeMap()
	isExist, _ := common.IsExist(oLogPath)
	if isExist {
		return
	}
	f, _ := common.OpenFile(oLogPath, true)
	AddToCsv("Name,Start Time,Over Time,Result,Target Path,Source Path", false)
	defer common.CloseFile(f)
}

func initQTP() {
	isExist, _ := common.IsExist(timePointPath)
	if isExist {
		return
	}
	f, _ := common.OpenFile(timePointPath, true)
	AddToCsv("Name,Time,Source Path,TimePoint Path,Marks", true)
	defer common.CloseFile(f)
}

func LoadCSV(isTimePoint bool) []string {
	path := ""
	if isTimePoint {
		path = timePointPath
	} else {
		path = oLogPath
	}
	ret := make([]string, 0)
	f, err := os.Open(path)
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

func AddToCsv(str string, isTimePoint bool) bool {
	var path string
	if isTimePoint {
		path = timePointPath
	} else {
		path = oLogPath
	}
	if str == "" {
		return true
	}
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0777)
	defer f.Close()
	f.Write([]byte(str + "\n"))
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
