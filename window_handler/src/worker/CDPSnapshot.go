package worker

import (
	"log"
	"os/exec"
	"strings"
	"window_handler/common"
	"window_handler/config"
)

// CreateTimePoint worker对外暴露的创建时间点方法
func CreateTimePoint(name string, sourcePath string, targetPath string, marks string, needLog bool, sn *string) {
	timePointPath := targetPath + "/" + *sn
	if needLog {
		config.AddToCsv(
			config.GetCsvStr(name, getNowTimeStr(), sourcePath, timePointPath, marks),
			true,
		)
	}
	createTimePoint(sourcePath, timePointPath, sn)
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
			createWindowsLink(sourceAbsPath, targetAbsPath)
		} else {
			createTimePoint(sourceAbsPath, targetAbsPath, sn)
		}
	}
}

func createWindowsLink(sourcePath string, targetPath string) {
	source := strings.ReplaceAll(sourcePath, "/", "\\")
	target := strings.ReplaceAll(targetPath, "/", "\\")
	cmd := exec.Command("cmd", "/C", "mklink", "/H", target, source)
	//执行命令
	err := cmd.Run()
	if err != nil {
		//log.Printf("create link err, source : %v, target : %v, err: %v", source, target, out)
	}
}
