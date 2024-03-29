package worker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
	"window_handler/common"
	"window_handler/config"
)

func GetFileMd5(f *os.File) *string {
	md5h := md5.New()
	_, err := io.Copy(md5h, f)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	md5Str := hex.EncodeToString(md5h.Sum(nil))
	return &md5Str
}

func CompareMd5(sf *os.File, tf *os.File) bool {
	sfMd5Ptr := GetFileMd5(sf)
	tfMd5Ptr := GetFileMd5(tf)
	return sfMd5Ptr == tfMd5Ptr
}

func CompareAndCacheMd5(sf *os.File, tf *os.File) bool {
	sfMd5Ptr := GetFileMd5(sf)
	tfMd5Ptr := GetFileMd5(tf)
	return *sfMd5Ptr == *tfMd5Ptr
}

func CompareModifyTime(sf *os.File, tf *os.File) bool {
	sfInfo, err := sf.Stat()
	if err != nil {
		log.Printf("get file stat error : %v", err)
		return false
	}
	tfInfo, err := tf.Stat()
	if err != nil {
		log.Printf("get file stat error : %v", err)
		return false
	}
	return sfInfo.ModTime() == tfInfo.ModTime()
}

func uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func stringToUint64(s string) (uint64, error) {
	intNum, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint64(intNum), nil
}

// CreateFile fileSize : KB
func CreateFile(bufferSize CapacityUnit, filePath string, fileSize CapacityUnit, randomContent bool) (success bool, usedTime float64) {
	exist, err := common.IsExist(filePath)
	if exist {
		return false, 1
	}

	f, err := common.OpenFile(filePath, true)
	if err != nil {
		return false, 1
	}
	defer f.Close()
	content := make([]byte, bufferSize)
	if randomContent {
		content = randomPalindrome(bufferSize)
	}
	startTime := time.Now()
	countSum := int(fileSize / bufferSize)
	for count := 1; count <= countSum; count++ {
		f.Write(content)
	}
	f.Sync()
	usedTime = float64(time.Now().Sub(startTime)) / float64(time.Second)
	if usedTime == 0 {
		usedTime++
	}
	return true, usedTime
}

func readFile(filePath string, bufferSize CapacityUnit) (success bool, usedTime float64) {
	exist, _ := common.IsExist(filePath)
	if !exist {
		return false, 1
	}
	startTime := time.Now()
	f, err := common.OpenFile(filePath, true)
	if err != nil {
		return false, 1
	}
	buffer := make([]byte, bufferSize)
	for {
		_, err = f.Read(buffer)
		if err == io.EOF {
			break
		}
	}
	f.Close()
	usedTime = float64(time.Now().Sub(startTime)) / float64(time.Second)
	if usedTime == 0 {
		usedTime++
	}
	return true, usedTime
}

// randomPalindrome size : byte
func randomPalindrome(size CapacityUnit) []byte {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	bytes := make([]byte, size)
	for i := 0; i < (int(size)+1)/2; i++ {
		r := byte(rng.Intn(0x1000)) //random rune up to '\u0999'
		bytes[i] = r
		bytes[int(size)-1-i] = r
	}
	return bytes
}

func ConvertCapacity(str string) CapacityUnit {
	regFindNum, _ := regexp.Compile(`\d+\.?\d*`)
	numStr := regFindNum.FindAllString(str, -1)[0]
	regFindUnit, _ := regexp.Compile(`[A-Z]+`)
	unit := regFindUnit.FindAllString(str, -1)[0]
	var totalCap CapacityUnit
	for k, v := range CapacityStrMap {
		if k == unit {
			totalCap = v
		}
	}
	num, _ := strconv.ParseFloat(numStr, 64)
	return CapacityUnit(int64(num)) * totalCap
}

func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func GetSuitableCapacityStr(c uint64) string {
	var ret string
	floatC := float64(c)
	gbNum := floatC / float64(GB)
	if gbNum < 1 {
		if floatC/float64(MB) < 1 {
			ret = fmt.Sprintf("%vKB", FloatRound(floatC/float64(KB), 2))
		} else {
			ret = fmt.Sprintf("%vMB", FloatRound(floatC/float64(MB), 2))
		}
	} else if gbNum > 1024 {
		ret = fmt.Sprintf("%vTB", FloatRound(floatC/float64(TB), 2))
	} else {
		ret = fmt.Sprintf("%vGB", FloatRound(floatC/float64(GB), 2))
	}
	return ret
}

func GetNextSyncTime(dayArray [7]bool, min uint8, hour uint8) time.Duration {
	var everyDayFlag = true
	var subs [7]int
	for i := 0; i < len(dayArray); i++ {
		value := dayArray[i]
		everyDayFlag = everyDayFlag && value
		if value {
			subs[i] = i - int(time.Now().Weekday())
		} else {
			subs[i] = -10
		}
	}
	hourSub := int(hour) - time.Now().Hour()
	minSub := int(min) - time.Now().Minute()
	if everyDayFlag {
		return GetTimeSum(0, hourSub, minSub)
	} else {
		nextDay := getClosetDaySub(subs, minSub, hourSub)
		return GetTimeSum(nextDay, hourSub, minSub)
	}
}

func GetTimeSum(daySub int, hourSub int, minSub int) time.Duration {
	var retTime time.Duration
	var dayFlag = false
	if minSub < 0 {
		daySub = 7
		dayFlag = true
	}

	if daySub < 0 {
		daySub = daySub + 7
	}

	if hourSub < 0 {
		hourSub = hourSub + 24
	}

	if minSub < 0 {
		minSub = minSub + 60
	}

	retTime = retTime + time.Duration(daySub*int(time.Hour*24))
	if dayFlag {
		retTime = retTime - time.Hour
	}
	retTime = retTime + time.Duration(hourSub*int(time.Hour))
	retTime = retTime + time.Duration(minSub*int(time.Minute))
	return retTime
}

func getNextTimeFromConfig(isBatchSync bool, isRemoteSync bool, isPartition bool) time.Duration {
	configCache := config.SystemConfigCache.GetLocalSyncPolicy(isBatchSync, isPartition)
	nextTime := GetNextSyncTime(
		configCache.TimingSync.Days,
		configCache.TimingSync.Minute,
		configCache.TimingSync.Hour,
	)
	if nextTime == 0 {
		time.Sleep(61 * time.Second)
	}
	nextTime = GetNextSyncTime(
		configCache.TimingSync.Days,
		configCache.TimingSync.Minute,
		configCache.TimingSync.Hour,
	)
	return nextTime
}

// getClosetDaySub 比较日期差，获取最近的那个日期
func getClosetDaySub(subs [7]int, minSub int, hourSub int) int {
	shift := false
	if (hourSub == 0 && minSub < 0) || hourSub < 0 {
		shift = true
	}
	var minNum = getMinTimeSubNum(&subs)
	var secondNum = getMinTimeSubNum(&subs)

	if shift {
		return secondNum
	}
	return minNum
}

// getMinNum 获取最小的时间差数字（正数：返回最小值，负数：返回最大值，不比较-10这个特殊数字）
func getMinTimeSubNum(subs *[7]int) int {
	var minNum = subs[0]
	var minIndex = -1
	for i := 0; i < len(subs); i++ {
		if subs[i] == -10 {
			continue
		} else if minNum == -10 {
			minNum = subs[i]
		}
		if subs[i] == 0 {
			minNum = subs[i]
			minIndex = i
			subs[i] = -10
			break
		}
		if subs[i] > 0 {
			if minNum < 0 {
				minNum = subs[i]
				minIndex = i
			} else if minNum > subs[i] {
				minNum = subs[i]
				minIndex = i

			}
		} else if subs[i] < minNum {
			if minNum > 0 {
				continue
			}
			minNum = subs[i]
			minIndex = i
		}
	}
	if minIndex != -1 {
		subs[minIndex] = -10
	}
	return minNum
}

func ReverseCompareAndDelete(sourcePath string, targetPath string) {
	exist0, err0 := common.IsExist(sourcePath)
	exist1, err1 := common.IsExist(targetPath)
	if err0 != nil || err1 != nil || !exist0 || !exist1 {
		return
	}
	sf, _ := common.OpenFile(sourcePath, false)
	tf, _ := common.OpenFile(targetPath, false)
	defer common.CloseFile(tf, sf)
	sfChildMap := make(map[string]int)
	tfChild, _ := tf.Readdir(-1)
	sfChild, _ := sf.Readdir(-1)
	for _, child := range sfChild {
		sfChildMap[child.Name()] = 1
	}
	for _, child := range tfChild {
		if sfChildMap[child.Name()] == 0 {
			common.DeleteFileOrDir(targetPath + fileSeparator + child.Name())
		}
	}
}

func GetTotalSize(sn *string, startPath string, isRoot bool, lock *sync.WaitGroup) {
	lock.Add(1)
	defer lock.Done()
	for _, v := range DiskPartitionsCache {
		if v.Name == startPath {
			GetPartitionsInfo()
			addSizeToTotalMap(*sn, v.TotalSize-v.FreeSize)
			return
		}
	}
	f, err := common.OpenDir(startPath)
	defer common.CloseFile(f)
	if err != nil {
		log.Printf("GetTotalSize err: %v", err)
		return
	}
	children, _ := f.Readdir(-1)
	if len(children) == 0 {
		fs, _ := f.Stat()
		addSizeToTotalMap(*sn, uint64(fs.Size()))
		return
	}

	for _, child := range children {
		if child.IsDir() {
			if isRoot {
				go GetTotalSize(sn, startPath+fileSeparator+child.Name(), false, lock)
			} else {
				GetTotalSize(sn, startPath+fileSeparator+child.Name(), false, lock)
			}
		} else {
			addSizeToTotalMap(*sn, uint64(child.Size()))
		}
	}
}

// EstimatedTotalTime 估计完成所需的时间
func EstimatedTotalTime(sn string, timeCell time.Duration) time.Duration {
	if getTotalSize(sn) == getDoneSize(sn) {
		return 0
	}
	doneSizeStart := getDoneSize(sn)
	time.Sleep(timeCell)
	doneSizeOver := getDoneSize(sn) - doneSizeStart
	speed := doneSizeOver / uint64(timeCell/time.Second)
	if speed == 0 {
		return EstimatedTotalTime(sn, timeCell)
	}
	tTime := (getTotalSize(sn) - getDoneSize(sn)) / speed
	return time.Duration(tTime * uint64(time.Second))
}

func GetFileName(path string) string {
	f, err := common.OpenFile(path, false)
	fInfo, _ := f.Stat()
	defer common.CloseFile(f)
	if err == nil {
		return fInfo.Name()
	}
	return ""
}

// recordOLog 记录操作日志
func recordOLog(busType int, startTime string, target string, source string) {
	if !config.SystemConfigCache.Value().SystemSetting.EnableOLog {
		return
	}
	overTime := getNowTimeStr()
	logStr := config.GetOLogType(busType) + "," + startTime + "," + overTime + ",success," + target + "," + source
	config.AddToCsv(logStr, false)
}

// recordTPLog 记录时间点日志
func recordTPLog(busType int, startTime string, target string, source string) {
	if !config.SystemConfigCache.Value().SystemSetting.EnableOLog {
		return
	}
	overTime := getNowTimeStr()
	logStr := config.GetOLogType(busType) + "," + startTime + "," + overTime + ",success," + target + "," + source
	config.AddToCsv(logStr, false)
}

func getNowTimeStr() string {
	now := time.Now()
	ret := fmt.Sprintf("%v", now.Format("2006/01/02 15:04:05"))
	return ret
}

// GetFileTree 获取文件树
func GetFileTree(startPath string, layer int) {

}
