package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"time"
	"window_handler/config"
)

const isDev = true

const NOT_RUNNING_STR = "Not running"

var mainWin *fyne.Window
var (
	syncErrorDialog dialog.Dialog
	waitAuthDialog  dialog.Dialog
	authErrorDialog dialog.Dialog
)
var syncErrorDialogOK = false

var timeCycleMap = make(map[string]time.Duration)
var dayCycleMap = make(map[string]time.Weekday)
var dayArrayList = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

var classicSize = fyne.Size{
	Height: config.WindowHeight * 0.8,
	Width:  config.WindowWidth,
}
var disableRootCache = make(map[string]map[fyne.Disableable]disableRoot)

/*
*
Selectivity between some components presents strong correlation
*/
type disableRoot struct {
	child []fyne.Disableable
}

func (p *disableRoot) addChild(childs ...fyne.Disableable) {
	for _, v := range childs {
		p.child = append(p.child, v)
	}
}

func (p *disableRoot) disableChild() {
	for _, v := range p.child {
		v.Disable()
	}
}

func (p *disableRoot) enableChild() {
	for _, v := range p.child {
		v.Enable()
	}
}

type Navigation struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

var (
	Navigations = map[string]Navigation{
		"localSync": {
			"Local Sync",
			"",
			getBatchLocalSyncComponent,
		},
		"localBatchSync": {
			"Local Batch Sync",
			"Click start button to begin sync",
			getBatchLocalSyncComponent,
		},
		"localSingleSync": {
			"Local Single Sync",
			"Click start button to begin sync.",
			getSingleLocalSyncComponent,
		},
		"partitionSync": {
			"Partition Sync",
			"Click start button to begin sync",
			getPartitionSyncComponent,
		},
		"remoteSync": {
			"Remote Single Sync",
			"Click start button to begin sync.",
			getManageRemoteQNQComponent,
		},
		"manageRemoteQNQ": {
			"Mange Remote QNQ",
			"Manage remote qnq.\nPlease auth the connection first!",
			getManageRemoteQNQComponent,
		},
		"remoteSingleSync": {
			"Remote Single Sync",
			"This is a demo!",
			getRemoteSingleComponent,
		},
		"systemInfo": {
			"System Information",
			"",
			getLocalSystemInfoComponent,
		},
		"diskInfo": {
			"Disk Information",
			"Basic Disk Information",
			getDiskInfoComponent,
		},
		"testDiskSpeed": {
			"Test Disk Speed",
			"The recommended buffer size is 4MB.",
			getTestDiskSpeedComponent,
		},
		"systemSetting": {
			"QNQ System Setting",
			"",
			getTestDiskSpeedComponent,
		},
		"setting": {
			"QNQ System Setting",
			"",
			getTestDiskSpeedComponent,
		},
		"olog": {
			"QNQ Log",
			"",
			getLogComponent,
		},
		"snapshot": {
			"Snapshot",
			"",
			getTimePointComponent,
		},
		"timePoint": {
			"Time Point",
			"Currently, only the same disk creation time point is supported.",
			getTimePointComponent,
		},
	}
	//设置菜单树
	NavigationIndex = map[string][]string{
		"":           {"localSync", "remoteSync", "systemInfo", "systemSetting"},
		"localSync":  {"localBatchSync", "localSingleSync", "partitionSync"},
		"remoteSync": {"manageRemoteQNQ", "remoteSingleSync"},
		//"snapshot":      {"timePoint"},
		"systemInfo":    {"diskInfo", "testDiskSpeed"},
		"systemSetting": {"olog"},
	}
)
