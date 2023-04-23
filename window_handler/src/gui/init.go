package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"net/url"
	"time"
	"window_handler/common"
	"window_handler/config"
	"window_handler/request"
	"window_handler/worker"
)

var topWindow fyne.Window

const preferenceCurrentNavigation = "currentNavigation"

func init() {
	I18n()
	initRegisterGWFunc()
	go watchGWChannel()
}

func I18n() {
	initTimeCycle()
}

func StartGUI() {
	a := app.NewWithID("qnq.window_handler")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("QNQ Sync " + config.Version)
	topWindow = w
	w.SetMaster()
	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setNavigation := func(t Navigation) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setNavigation, false))
	} else {
		split := container.NewHSplit(makeNav(setNavigation, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	SetMainWin(&w)
	w.Resize(fyne.NewSize(config.WindowWidth, config.WindowHeight))
	w.ShowAndRun()
}

// func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
//
// }
func makeNav(setNavigation func(navigation Navigation), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return NavigationIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := NavigationIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Navigations[uid]
			if !ok {
				fyne.LogError("Missing gui panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedNavigation(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := Navigations[uid]; ok {
				if unsupportedNavigation(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentNavigation, uid)
				setNavigation(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentNavigation, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)
	lastVersion := request.GetLastVersion()
	if lastVersion == "" || lastVersion == config.SystemConfigCache.Value().Version {
		lastVersion = "Project Link"
	} else {
		lastVersion = "Lastest version is  " + lastVersion
	}
	projectUrl, _ := url.Parse("https://github.com/wangshenghao1/QNQ")
	versionInfo := widget.NewHyperlink(lastVersion, projectUrl)
	flootBox := container.NewVBox(
		versionInfo,
		themes,
	)

	return container.NewBorder(nil, flootBox, nil, nil, tree)
}

func unsupportedNavigation(t Navigation) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func SetMainWin(win *fyne.Window) {
	mainWin = win
	syncErrorDialog = dialog.NewInformation("Sync task warning!", "Sync task enters repeatedly, please adjust the time interval.", *mainWin)
}

func initTimeCycle() {
	timeCycleMap["Second"] = time.Second
	timeCycleMap["Minute"] = time.Minute
	timeCycleMap["Hour"] = time.Hour

	dayCycleMap[dayArrayList[0]] = time.Sunday
	dayCycleMap[dayArrayList[1]] = time.Monday
	dayCycleMap[dayArrayList[2]] = time.Tuesday
	dayCycleMap[dayArrayList[3]] = time.Wednesday
	dayCycleMap[dayArrayList[4]] = time.Thursday
	dayCycleMap[dayArrayList[5]] = time.Friday
	dayCycleMap[dayArrayList[6]] = time.Saturday
}

var gwChannelRegisterF = make(map[int]func(), 16)

// registerGWFunc GW管道的响应函数均在此注册
func registerGWFunc(signal int, f func()) {
	gwChannelRegisterF[signal] = f
}

// initRegisterGWFunc 初始化注册GW响应函数
func initRegisterGWFunc() {
	//local batch sync
	registerGWFunc(common.GetRunningSignal(common.TYPE_LOCAL_BATCH), localBatchRunningHandle)
	registerGWFunc(common.GetForceDoneSignal(common.TYPE_LOCAL_BATCH), localBatchDoneHandle)
	//local single sync
	registerGWFunc(common.GetRunningSignal(common.TYPE_LOCAL_SING), localSingleRunningHandle)
	registerGWFunc(common.GetForceDoneSignal(common.TYPE_LOCAL_SING), localSingleDoneHandle)
	//test disk speed
	registerGWFunc(common.GetRunningSignal(common.TYPE_TEST_SPEED), testSpeedRunningHandle)
	registerGWFunc(common.GetForceDoneSignal(common.TYPE_TEST_SPEED), testSpeedDoneHandle)
	//partition sync
	registerGWFunc(common.GetRunningSignal(common.TYPE_PARTITION), partitionRunningHandle)
	registerGWFunc(common.GetForceDoneSignal(common.TYPE_PARTITION), partitionDoneHandle)
	//partition sync
	registerGWFunc(common.GetRunningSignal(common.TYPE_CREATE_TIMEPOINT), createTPRunningHandle)
	registerGWFunc(common.GetForceDoneSignal(common.TYPE_CREATE_TIMEPOINT), createTPDoneHandle)
}

func watchGWChannel() {
	for {
		select {
		case c := <-common.GWChannel:
			f := gwChannelRegisterF[c]
			if f == nil {
				log.Printf("!!!!!!!!!!!!!!!!!!has a signal doesn't register, num : %v", c)
			}
			f()
		}
	}
}

func testSpeedRunningHandle() {
	testSpeedRetLab.SetText("Testing...")
}

func testSpeedDoneHandle() {
	partition := speedPartitionSelect.Selected
	rSpeed := fmt.Sprint(worker.DiskReadSpeedCache[partition])
	wSpeed := fmt.Sprint(worker.DiskWriteSpeedCache[partition])
	testSpeedRetLab.SetText("Disk : " + partition + "\n" + "Read speed : " + rSpeed + "MB/S\n" + "Write speed : " + wSpeed + "MB/S\n")
}

func localBatchRunningHandle() {
	batchDisable(localBatchSyncPolicyComponent, localBatchStartButton)
	batchEnable(localBatchCancelButton)
	showSyncError(common.TYPE_LOCAL_BATCH)
	startSyncGUI(localBatchProgressBox, localBatchCurrentFile, localBatchTimeRemaining, common.TYPE_LOCAL_BATCH)
}

func localBatchDoneHandle() {
	overSyncGUI(localBatchProgressBox, localBatchCurrentFile, localBatchTimeRemaining)
	batchEnable(localBatchSyncPolicyComponent, localBatchStartButton)
	batchDisable(localBatchCancelButton)
}

func localSingleRunningHandle() {
	batchDisable(localSingleSyncPolicyComponent, localSingleStartButton)
	showSyncError(common.TYPE_LOCAL_SING)
	startSyncGUI(localSingleProgressBox, localSingleCurrentFile, localSingleTimeRemaining, common.TYPE_LOCAL_SING)
}

func localSingleDoneHandle() {
	batchEnable(localSingleSyncPolicyComponent, localSingleStartButton)
	overSyncGUI(localSingleProgressBox, localSingleCurrentFile, localSingleTimeRemaining)
}

func partitionRunningHandle() {
	batchDisable(partitionSyncPolicyComponent, partitionStartButton)
	batchEnable(partitionCancelButton)
	startSyncGUI(partitionProgressBox, partitionCurrentFile, partitionTimeRemaining, common.TYPE_PARTITION)
}

func partitionDoneHandle() {
	batchEnable(partitionSyncPolicyComponent, partitionStartButton)
	batchDisable(partitionCancelButton)
	overSyncGUI(partitionProgressBox, partitionCurrentFile, partitionTimeRemaining)
}

func createTPRunningHandle() {

}

func createTPDoneHandle() {

}

func showSyncError(busType int) {
	runningFlag := common.GetRunningFlag(busType)
	if runningFlag {
		syncErrorDialogOK = true
		if !syncErrorDialogOK {
			syncErrorDialog.Show()
		}
	}
}
