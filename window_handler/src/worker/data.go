package worker

const dataMsgPreFix = "00"

const (
	VARIANCE_ROOT = iota
	VARIANCE_ADD
	VARIANCE_EDIT
	VARIANCE_DELETE
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

// LocalBSFileNode /** Local batch source file node
var LocalBSFileNode *FileNode

// LocalBTFileNode /** Local batch target file node
var LocalBTFileNode *FileNode

func getFileSeparator() string {
	//if strings.Contains(runtime.GOOS, "window") {
	//	return "\\"
	//} else if strings.Contains(runtime.GOOS, "linux") {
	//	return "/"
	//}
	return "/"
}
