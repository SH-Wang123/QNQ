package worker

import (
	"github.com/shirou/gopsutil/disk"
	"window_handler/common"
)

var DiskPartitionsCache []Partition
var DiskReadSpeedCache = make(map[string]float64)
var DiskWriteSpeedCache = make(map[string]float64)

func GetPartitionsInfo() {
	var partitions []Partition
	partitionsInfo, _ := disk.Partitions(true)

	for _, info := range partitionsInfo {
		moreInfo, _ := disk.Usage(info.Device)
		totalStr := GetSuitableCapacityStr(moreInfo.Total)
		freeStr := GetSuitableCapacityStr(moreInfo.Free)
		var p = Partition{
			Name:         info.Device,
			FsType:       info.Fstype,
			TotalSize:    moreInfo.Total,
			TotalSizeStr: totalStr,
			FreeSizeStr:  freeStr,
			FreeSize:     moreInfo.Free,
			UsedPercent:  moreInfo.UsedPercent / 100,
		}
		partitions = append(partitions, p)
	}
	DiskPartitionsCache = partitions
}

func TestDiskSpeed(bufferSize CapacityUnit, totalSize CapacityUnit, drive string) (writeSpeed int, readSpeed int) {
	common.GWChannel <- common.TEST_DISK_SPEED_START
	defer func() {
		common.GWChannel <- common.TEST_DISK_SPEED_OVER
	}()
	fileName := drive + "/test_speed"
	DeleteFile(fileName)

	_, wirteTime := CreateFile(bufferSize, fileName, totalSize, true)
	defer DeleteFile(fileName)
	DiskWriteSpeedCache[drive] = FloatRound(float64(getMb(totalSize))/wirteTime, 2)

	_, readTime := ReadFile(fileName, bufferSize)
	DiskReadSpeedCache[drive] = FloatRound(float64(getMb(totalSize))/readTime, 2)
	return writeSpeed, 0
}

func getMb(size CapacityUnit) int {
	return int(int64(size) / int64(MB))
}
