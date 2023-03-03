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
	"time"
	"window_handler/config"
)

func OpenFile(filePath string, createFile bool) (*os.File, error) {
	var f *os.File
	var err error
	if createFile {
		f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	} else {
		f, err = os.Open(filePath)
	}
	if err != nil {
		log.Printf("Open %v err : %v", filePath, err.Error())
		return nil, err
	}
	return f, nil
}

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Printf("close file err : %v", err.Error())
	}
}

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(path string) {
	exist, err := IsExist(path)
	if err != nil {
		log.Printf("get dir error : %v", err)
	}
	if !exist {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("create dir error : %v", err)
		}
	}
}

func DeleteFile(path string) error {
	exist, err := IsExist(path)
	if err != nil {
		log.Printf("get dir error : %v", err)
	}
	if exist {
		err = os.Remove(path)
		return err
	}
	return nil
}

func DeleteDir(path string) error {
	exist, err := IsExist(path)
	if err != nil {
		log.Printf("get dir error : %v", err)
	}
	if exist {
		err = os.RemoveAll(path)
		return err
	}
	return nil
}

func IsOpenDirError(err error, path string) bool {
	return err.Error() == "open "+path+": is a directory"
}

func GetFileMd5(f *os.File) *string {
	md5h := md5.New()
	io.Copy(md5h, f)
	md5Str := hex.EncodeToString(md5h.Sum(nil))
	return &md5Str
}

func CompareMd5(sf *os.File, tf *os.File) bool {
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

func GetSingleFileNode(path string) *FileNode {
	return &FileNode{
		IsDirectory:     false,
		HasChildren:     false,
		AbstractPath:    path,
		AnchorPointPath: "",
		HeadFileNode:    nil,
		VarianceType:    VARIANCE_ROOT,
	}
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

// fileSize : KB
func CreateFile(bufferSize CapacityUnit, filePath string, fileSize CapacityUnit, randomContent bool) (success bool, usedTime float64) {
	exist, err := IsExist(filePath)
	if exist {
		return false, 1
	}

	f, err := OpenFile(filePath, true)
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
	usedTime = float64(time.Now().Sub(startTime)) / float64(time.Second)
	if usedTime == 0 {
		usedTime++
	}
	return true, usedTime
}

func ReadFile(filePath string, bufferSize CapacityUnit) (success bool, usedTime float64) {
	exist, _ := IsExist(filePath)
	if !exist {
		return false, 1
	}
	startTime := time.Now()
	f, err := OpenFile(filePath, true)
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
	regFindNum, _ := regexp.Compile(`\d+`)
	numStr := regFindNum.FindAllString(str, -1)[0]
	regFindUnit, _ := regexp.Compile(`\D+`)
	unit := regFindUnit.FindAllString(str, -1)[0]
	var totalCap CapacityUnit
	for k, v := range CapacityStrMap {
		if k == unit {
			totalCap = v
		}
	}
	num, _ := strconv.Atoi(numStr)
	return CapacityUnit(int64(num)) * totalCap
}

func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func GetSuitableCapacityStr(c uint64) string {
	var ret string
	if c/uint64(GB) == 0 {
		if c/uint64(MB) == 0 {
			ret = fmt.Sprintf("%vKB", c/uint64(KB))
		} else {
			ret = fmt.Sprintf("%vMB", c/uint64(MB))
		}
	} else {
		ret = fmt.Sprintf("%vGB", c/uint64(GB))
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
	retTime = retTime + time.Duration(hourSub*int(time.Hour))
	retTime = retTime + time.Duration(minSub*int(time.Minute))
	return retTime
}

func GetNextTimeFromConfig(isBatchSync bool, isRemoteSync bool) time.Duration {
	configCache := config.SystemConfigCache.GetSyncPolicy(isBatchSync, isRemoteSync)
	return GetNextSyncTime(
		configCache.TimingSync.Days,
		configCache.TimingSync.Minute,
		configCache.TimingSync.Hour,
	)
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
