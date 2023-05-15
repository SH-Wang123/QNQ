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
var menuTree *widget.Tree

const preferenceCurrentNavigation = "currentNavigation"

func init() {
	i18n()
	initRegisterWGFunc()
	go watchWGChannel()
}

func i18n() {
	initTimeCycle()
}

func initGlobalDialog() {
	syncErrorDialog = dialog.NewInformation("Sync task warning!", "Sync task enters repeatedly, please adjust the time interval.", *mainWin)
	waitAuthDialog = dialog.NewInformation("Waiting remote qnq auth", "Please wait.", *mainWin)
	authErrorDialog = dialog.NewInformation("Remote qnq auth error", "The link is blocked or the other party does not agree.", *mainWin)
}

func StartGUI() {
	a := app.NewWithID("qnq.window_handler")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("QNQ Sync " + config.Version)
	topWindow = w
	w.SetMaster()
	content := container.NewMax()
	title := widget.NewLabel("QNQ Sync")
	intro := widget.NewLabel("Welcome to use qnq.")
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
	split := container.NewHSplit(makeNav(setNavigation, true), tutorial)
	split.Offset = 0.2
	w.SetContent(split)
	SetMainWin(&w)
	w.Resize(fyne.NewSize(config.WindowWidth, config.WindowHeight))
	initGlobalDialog()
	w.ShowAndRun()
}

// func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
//
// }
func makeNav(setNavigation func(navigation Navigation), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	menuTree = &widget.Tree{
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
		menuTree.Select(currentPref)
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
		lastVersion = "Latest version is  " + lastVersion
	}
	projectUrl, _ := url.Parse("https://github.com/wangshenghao1/QNQ")
	versionInfo := widget.NewHyperlink(lastVersion, projectUrl)
	flootBox := container.NewVBox(
		versionInfo,
		themes,
	)

	return container.NewBorder(nil, flootBox, nil, nil, menuTree)
}

func unsupportedNavigation(t Navigation) bool {
	return fyne.CurrentDevice().IsBrowser()
}

func SetMainWin(win *fyne.Window) {
	mainWin = win
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

var wgChannelRegisterF = make(map[int]func(), 16)

// registerGWFunc GW管道的响应函数均在此注册
func registerWGFunc(signal int, f func()) {
	wgChannelRegisterF[signal] = f
}

// initRegisterGWFunc 初始化注册GW响应函数
func initRegisterWGFunc() {
	//local batch sync
	registerWGFunc(common.GetRunningSignal(common.TYPE_LOCAL_BATCH), localBatchRunningHandle)
	registerWGFunc(common.GetForceDoneSignal(common.TYPE_LOCAL_BATCH), localBatchDoneHandle)
	//local single sync
	registerWGFunc(common.GetRunningSignal(common.TYPE_LOCAL_SING), localSingleRunningHandle)
	registerWGFunc(common.GetForceDoneSignal(common.TYPE_LOCAL_SING), localSingleDoneHandle)
	//test disk speed
	registerWGFunc(common.GetRunningSignal(common.TYPE_TEST_SPEED), testSpeedRunningHandle)
	registerWGFunc(common.GetForceDoneSignal(common.TYPE_TEST_SPEED), testSpeedDoneHandle)
	//partition sync
	registerWGFunc(common.GetRunningSignal(common.TYPE_PARTITION), partitionRunningHandle)
	registerWGFunc(common.GetForceDoneSignal(common.TYPE_PARTITION), partitionDoneHandle)
	//partition sync
	registerWGFunc(common.GetRunningSignal(common.TYPE_CREATE_TIMEPOINT), createTPRunningHandle)
	registerWGFunc(common.GetForceDoneSignal(common.TYPE_CREATE_TIMEPOINT), createTPDoneHandle)
	//remote qnq auth
	registerWGFunc(common.GetRunningSignal(common.TYPE_REMOTE_QNQ_AUTH), qnqAuthRunningHandle)
}

func watchWGChannel() {
	for {
		select {
		case c := <-common.WGChannel:
			f := wgChannelRegisterF[c]
			if f == nil {
				log.Printf("!!!!!!!!!!!!!!!!!!has a signal doesn't register, num : %v, wg", c)
				continue
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

func qnqAuthRunningHandle() {
	remoteAuthDialog = dialog.NewConfirm("Remote QNQ want to auth", common.CurrentWaitAuthIp, remoteAuthDialogFunc, *mainWin)
	remoteAuthDialog.SetDismissText("Refuse")
	remoteAuthDialog.SetConfirmText("Agree")
	remoteAuthDialog.SetDismissText("Refuse(60s)")
	remoteAuthShowing = true
	go func() {
		i := 0
		for {
			if !remoteAuthShowing {
				return
			}
			i++
			time.Sleep(1 * time.Second)
			remoteAuthDialog.SetDismissText(fmt.Sprintf("Refuse(%vs)", 60-i))
			if i == 60 {
				remoteAuthDialog.Hide()
				common.SendSignal2GWChannel(common.SIGNAL_AUTH_NO_PASS)
				remoteAuthShowing = false
				return
			}
		}
	}()
	remoteAuthDialog.Show()
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
