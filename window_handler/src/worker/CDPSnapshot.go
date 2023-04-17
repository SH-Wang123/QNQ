package worker

import (
	"bufio"
	"window_handler/common"
)

const cdpNameFlag = "/QNQCDPSNAPSHOT0"
const cdpMd5Flag = "/QNQCDPSNAPSHOT1"

var updateFileCache = make([]string, 0)

func CrateCDPSnapshot(sourcePath string, targetPath string, forceInit bool) {
	sn := common.GetSNCount()
	if osName == windowsOSName {
		crateWindowsCDP(sourcePath, targetPath, forceInit, &sn)
	} else if osName == linuxOSName {
		createLinuxCDP(sourcePath, targetPath, forceInit, &sn)
	}
}

func crateWindowsCDP(sourcePath string, targetPath string, forceInit bool, sn *string) {

	syncCDPSnapshot(sourcePath, targetPath, forceInit, sn, false)
}

func createLinuxCDP(sourcePath string, targetPath string, forceInit bool, sn *string) {

	syncCDPSnapshot(sourcePath, targetPath, forceInit, sn, true)
}

func syncCDPSnapshot(sourcePath string, targetPath string, forceInit bool, sn *string, isLinux bool) bool {
	nameFlag := sourcePath + cdpNameFlag
	md5Flag := sourcePath + cdpMd5Flag
	exist0, err0 := common.IsExist(nameFlag)
	exist1, err1 := common.IsExist(md5Flag)
	if !exist0 || err0 != nil || !exist1 || err1 != nil {
		return false
	} else {
		nameFlagFile, err0 := common.OpenFile(nameFlag, false)
		md5FlagFile, err1 := common.OpenFile(md5Flag, false)
		defer common.CloseFile(nameFlagFile, md5FlagFile)
		if err0 != nil || err1 != nil {
			return false
		}
		scannerName := bufio.NewScanner(nameFlagFile)
		scannerMd5 := bufio.NewScanner(md5FlagFile)
		for scannerName.Scan() && scannerMd5.Scan() {
			fileName := scannerName.Text()
			oldMd5 := scannerMd5.Text()
			targetFile, err := common.OpenFile(fileName, false)
			if err != nil {
				continue
			}
			fInfo, _ := targetFile.Stat()
			sourcePathC := sourcePath + fInfo.Name()
			sourceExist, _ := common.IsExist(sourcePathC)
			if !sourceExist {
				common.DeleteFileOrDir(fileName)
			}
			newMd5 := *GetFileMd5(targetFile)
			if newMd5 != oldMd5 {
				sourceFile, _ := common.OpenFile(sourcePathC, false)
				worker := NewLocalSingleWorker(sourceFile, targetFile, *sn, true)
				worker.Execute()
			}
		}
	}
	return true
}

func createTimePoint(sourcePath string, targetPath string, force bool, sn *string) bool {
	md5Path := targetPath + cdpMd5Flag
	namePath := targetPath + cdpNameFlag
	exist0, _ := common.IsExist(md5Path)
	exist1, _ := common.IsExist(namePath)
	if exist0 && exist1 && !force {
		return false
	}
	if exist0 {
		common.DeleteFileOrDir(md5Path)
	}
	if exist1 {
		common.DeleteFileOrDir(namePath)
	}
	batchSyncFile(sourcePath, targetPath, sn, common.TYPE_CDP_SNAPSHOT, true)
	defer clearSfMd5Cache(sn)
	md5Cache := getSfMd5Cache(sn)
	md5File, _ := common.OpenFile(md5Path, true)
	nameFile, _ := common.OpenFile(namePath, true)
	defer common.CloseFile(md5File, nameFile)
	for name, md5 := range md5Cache {
		md5File.Write([]byte(*md5 + "\n"))
		nameFile.Write([]byte(name + "\n"))
	}
	return true
}
