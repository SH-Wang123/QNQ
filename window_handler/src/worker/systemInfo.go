package worker

import (
	"github.com/shirou/gopsutil/disk"
	"window_handler/common"
)

var DiskPartitionsCache []Partition
var DiskReadSpeedCache = make(map[string]int)
var DiskWriteSpeedCache = make(map[string]int)

func GetPartitionsInfo() {
	var partitions []Partition
	partitionsInfo, _ := disk.Partitions(true)
	for _, info := range partitionsInfo {
		moreInfo, _ := disk.Usage(info.Device)
		var p = Partition{
			Name:        info.Device,
			FsType:      info.Fstype,
			TotalSize:   uint64ToString(moreInfo.Total/1024/1024/1024) + "GB",
			FreeSize:    uint64ToString(moreInfo.Free/1024/1024/1024) + "GB",
			UsedPercent: moreInfo.UsedPercent / 100,
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
	err := DeleteFile(fileName)
	if err != nil {
		DiskWriteSpeedCache[drive] = -1
		return 0, 0
	}
	DiskWriteSpeedCache[drive] = int(getMb(totalSize) / wirteTime)

	return writeSpeed, 0
}

func getMb(size CapacityUnit) int {
	return int(int64(size) / int64(MB))
}
