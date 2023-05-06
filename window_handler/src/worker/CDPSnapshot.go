package worker

import (
	"log"
	"os/exec"
	"reflect"
	"strings"
	"window_handler/common"
	"window_handler/config"
)

// CreateTimePoint worker对外暴露的创建时间点方法
func CreateTimePoint(name string, sourcePath string, targetPath string, marks string, needLog bool) {
	sn := common.GetTaskCount()
	timePointPath := targetPath + "/" + name + sn
	common.SendSignal2WGChannel(common.GetRunningSignal(common.TYPE_CREATE_TIMEPOINT))
	defer common.SendSignal2WGChannel(common.GetForceDoneSignal(common.TYPE_CREATE_TIMEPOINT))
	startTime := getNowTimeStr()
	createTimePoint(sourcePath, timePointPath, &sn)
	if needLog {
		go config.AddToCsv(
			config.GetCsvStr(name, startTime, sourcePath, timePointPath, marks),
			true,
		)
	}
	go recordOLog(common.TYPE_CREATE_TIMEPOINT, startTime, timePointPath, sourcePath)
}

// createTimePoint 创建简单时间点
func createTimePoint(sourcePath string, timePointPath string, sn *string) {
	sf, err1 := common.OpenDir(sourcePath)
	if err1 != nil {
		log.Printf("create simple time point err: %v", err1)
		return
	}
	sfInfo, _ := sf.Stat()
	sfMode := sfInfo.Mode()
	common.CreateDir(timePointPath, &sfMode)
	tf, err2 := common.OpenDir(timePointPath)
	if err2 != nil {
		log.Printf("create simple time point err: %v", err2)
		return
	}
	defer common.CloseFile(sf, tf)
	children, _ := sf.Readdir(-1)
	for _, child := range children {
		targetAbsPath := timePointPath + fileSeparator + child.Name()
		sourceAbsPath := sourcePath + fileSeparator + child.Name()
		if !child.IsDir() {
			//createWindowsLink(sourceAbsPath, targetAbsPath, sn)
			common.SubmitFunc2Pool(reflect.ValueOf(createWindowsLink), sourceAbsPath, targetAbsPath, sn)
		} else {
			createTimePoint(sourceAbsPath, targetAbsPath, sn)
		}
	}
}

// createWindowsLink SN不为空字符串时，设置当前文件
func createWindowsLink(sourcePath string, targetPath string, sn *string) {
	if *sn != "" {

	}
	source := strings.ReplaceAll(sourcePath, "/", "\\")
	target := strings.ReplaceAll(targetPath, "/", "\\")
	cmd := exec.Command("cmd", "/C", "mklink", "/H", target, source)
	//执行命令
	err := cmd.Run()
	if err != nil {
		//log.Printf("create link err, source : %v, target : %v, err: %v", source, target, out)
	}
}
