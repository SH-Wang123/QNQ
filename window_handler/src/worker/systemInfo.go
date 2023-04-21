package worker

import (
	"github.com/shirou/gopsutil/disk"
	"net/http"
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
		if moreInfo == nil {
			continue
		}
		totalStr := GetSuitableCapacityStr(moreInfo.Total)
		freeStr := GetSuitableCapacityStr(moreInfo.Free)
		for i, _ := range DiskPartitionsCache {
			if DiskPartitionsCache[i].Name == info.Device {
				DiskPartitionsCache[i].FreeSizeStr = freeStr
				DiskPartitionsCache[i].FreeSize = moreInfo.Free
				DiskPartitionsCache[i].UsedPercent = moreInfo.UsedPercent / 100
				break
			}
		}

		if DiskPartitionsCache == nil {
			var p = Partition{
				Name:         info.Device,
				FsType:       info.Fstype,
				TotalSizeStr: totalStr,
				TotalSize:    moreInfo.Total,
				FreeSize:     moreInfo.Free,
				FreeSizeStr:  freeStr,
				UsedPercent:  FloatRound(moreInfo.UsedPercent/100, 4),
			}
			partitions = append(partitions, p)
		}
	}
	if DiskPartitionsCache == nil {
		DiskPartitionsCache = partitions
	}
}

func TestDiskSpeed(bufferSize CapacityUnit, totalSize CapacityUnit, drive string) (writeSpeed int, readSpeed int) {
	common.GWChannel <- common.TEST_DISK_SPEED_START
	defer func() {
		common.GWChannel <- common.TEST_DISK_SPEED_OVER
	}()
	fileName := drive + "/test_speed"
	common.DeleteFileOrDir(fileName)

	_, wirteTime := CreateFile(bufferSize, fileName, totalSize, true)
	defer common.DeleteFileOrDir(fileName)
	DiskWriteSpeedCache[drive] = FloatRound(float64(getMb(totalSize))/wirteTime, 2)

	_, readTime := readFile(fileName, bufferSize)
	DiskReadSpeedCache[drive] = FloatRound(float64(getMb(totalSize))/readTime, 2)
	return writeSpeed, 0
}

func getMb(size CapacityUnit) int {
	return int(int64(size) / int64(MB))
}

func GetRemoteDiskInfo(ip string) *[]Partition {
	r, err := http.Get(URL_HRED + ip + common.QNQ_TARGET_REST_PORT + GET_DISK_INFO_URI)
	if err != nil {
		return nil
	}
	var disks []Partition
	getObjFromResponse(r, &disks)
	return &disks
}
