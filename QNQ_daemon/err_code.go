package QNQ

const (
	OK_CODE  = 1
	ERR_CODE = 1<<32 - 1
)

//0x00000000

// local sync failed
const (
	SourceFileOpenFailed = iota + 1<<10
	TargetFileOpenFailed
	DirSync2FileError
	SyncFailed
	TargetNotConnected
)

// remote sync failed
const ()

var msgMap = make(map[int]string)

func initErrCode() {
	msgMap[SourceFileOpenFailed] = "Source File Open Failed"
	msgMap[TargetFileOpenFailed] = "Target File Open Failed"
	msgMap[DirSync2FileError] = "Source is dir but target is file"
	msgMap[SyncFailed] = "Sync Failed"
}

func GetMsg(key int) string {
	return msgMap[key]
}
