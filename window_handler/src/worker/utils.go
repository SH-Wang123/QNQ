package worker

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
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
	overTime := time.Now()
	usedTime = float64(overTime.Sub(startTime) / time.Second)
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
