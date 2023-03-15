package worker

const dataMsgPreFix = "00"

const (
	md5CheckError = "File sync error. Md5 is different"
)

const (
	URL_HRED = "http://"
)

const (
	VARIANCE_ROOT = iota
	VARIANCE_ADD
	VARIANCE_EDIT
	VARIANCE_DELETE
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

type FileNode struct {
	IsDirectory      bool
	HasChildren      bool
	AbstractPath     string
	AnchorPointPath  string
	ChildrenNodeList []*FileNode
	HeadFileNode     *FileNode
	VarianceType     int
}

type SyncFileError struct {
	AbsPath string
	Reason  string
}

// LocalBSFileNode /** Local batch source file node
var LocalBSFileNode *FileNode

// LocalBTFileNode /** Local batch target file node
var LocalBTFileNode *FileNode

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

func getFileSeparator() string {
	//if strings.Contains(runtime.GOOS, "window") {
	//	return "\\"
	//} else if strings.Contains(runtime.GOOS, "linux") {
	//	return "/"
	//}
	return "/"
}
