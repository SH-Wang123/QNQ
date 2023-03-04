package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"time"
	"window_handler/config"
)

const isDev = true

var mainWin *fyne.Window
var syncErrorDialog dialog.Dialog
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
	SupportWeb   bool
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
			GetBatchLocalSyncComponent,
			true,
		},
		"localBatchSync": {
			"Local Batch Sync",
			"Click start button to begin sync",
			GetBatchLocalSyncComponent,
			true,
		},
		"localSingleSync": {
			"Local Single Sync",
			"Click start button to begin sync.",
			GetSingleLocalSyncComponent,
			true,
		},
		"remoteSync": {
			"Remote Single Sync",
			"Click start button to begin sync.",
			GetRemoteSingleComponent,
			true,
		},
		"remoteSingleSync": {
			"Remote Single Sync",
			"Click start button to begin sync.\nPlease test the connection first!",
			GetRemoteSingleComponent,
			true,
		},
		"systemInfo": {
			"System Information",
			"",
			GetLocalSystemInfoComponent,
			true,
		},
		"diskInfo": {
			"Disk Information",
			"Basic Disk Information",
			GetDiskInfoComponent,
			true,
		},
		"testDiskSpeed": {
			"Test Disk Speed",
			"The recommended buffer size is 4MB.",
			GetTestDiskSpeedComponent,
			true,
		},
	}
	//设置菜单树
	NavigationIndex = map[string][]string{
		"":           {"localSync", "systemInfo", "remoteSync"},
		"localSync":  {"localBatchSync", "localSingleSync"},
		"systemInfo": {"diskInfo", "testDiskSpeed"},
		"remoteSync": {"remoteSingleSync"},
	}
)
