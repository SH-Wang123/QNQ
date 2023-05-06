package worker

const dataMsgPreFix = "00"

const (
	md5CheckError = "File sync error. Md5 is different"
)

const (
	URL_HRED = "http://"
)

type CapacityUnit uint64

const (
	Byte CapacityUnit = 1
	KB                = 1024 * Byte
	MB                = 1024 * KB
	GB                = 1024 * MB
	TB                = 1024 * GB
	PB                = 1024 * TB
)

const (
	GET_FILE_ROOT_URI = "/fileRootMap"
	GET_DISK_INFO_URI = "/disk/info"
	TEST_CONNECT      = "/debug/connect"
)

var fileSeparator = getFileSeparator()

type SyncFileError struct {
	AbsPath string
	Reason  string
}

type Disk struct {
	Name       string
	Partitions []Partition
	MediaType  int
	Speed      uint64
}

type Partition struct {
	Name         string  `json:"name"`
	FsType       string  `json:"fs_type"`
	TotalSizeStr string  `json:"total_size_str"`
	TotalSize    uint64  `json:"total_size"`
	FreeSizeStr  string  `json:"free_size_str"`
	FreeSize     uint64  `json:"free_size"`
	UsedPercent  float64 `json:"used_percent"`
}

type cdpInfo struct {
	Name string `json:"name"`
	Md5  string `json:"md5"`
}

func getFileSeparator() string {
	//if strings.Contains(runtime.GOOS, "window") {
	//	return "\\"
	//} else if strings.Contains(runtime.GOOS, "linux") {
	//	return "/"
	//}
	return "/"
}
