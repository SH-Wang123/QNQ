package worker

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
)

var DiskPartitionsCache []Partition

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

func TestDiskSpeed(bufferSize CapacityUnit, totalSize CapacityUnit, drive string) (writeSpeed string, readSpeed string) {
	fileName := drive + "/test_speed"
	_, wirteTime := CreateFile(bufferSize, fileName, totalSize, true)
	err := DeleteFile(fileName)
	if err != nil {
		return "N/A MB/s", "N/A MB/s"
	}
	writeSpeed = fmt.Sprint(int64(totalSize)/int64(MB)/int64(wirteTime)) + " MB/s"
	return writeSpeed, ""
}
