package navigations

import (
	"fyne.io/fyne/v2"
	"time"
)

const isDev = true

var timeCycleMap = make(map[string]time.Duration)

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
			"Remote",
			"QNQ Target Info",
			GetRemoteSyncComponent,
			true,
		},
		"remoteSingleSync": {
			"Remote sync",
			"Remote",
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
	}
	//设置菜单树
	NavigationIndex = map[string][]string{
		"":           {"localSync", "systemInfo", "remoteSync"},
		"localSync":  {"localBatchSync", "localSingleSync"},
		"systemInfo": {"diskInfo"},
		"remoteSync": {"remoteSingleSync"},
	}
)

func initTimeCycle() {
	timeCycleMap["Second"] = time.Second
	timeCycleMap["Minute"] = time.Minute
	timeCycleMap["Hour"] = time.Hour
}
