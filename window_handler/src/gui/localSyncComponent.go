package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"sync"
	"time"
	"window_handler/common"
	"window_handler/config"
	"window_handler/worker"
)

var FileSyncComponent fyne.CanvasObject

var (
	localSingleSyncComponent       fyne.CanvasObject
	localSingleSyncPolicyComponent *widget.Button
	localSingleStartButton         *widget.Button
	createFileComponent            *widget.Check
	localSingleProgressBox         *fyne.Container
	localSinglePolicySyncBox       *fyne.Container
	localSingleCurrentFile         *widget.Label
	localSingleTimeRemaining       *widget.Label
	slcOnce                        sync.Once
)

var (
	localBatchSyncComponent       fyne.CanvasObject
	localBatchSyncPolicyComponent *widget.Button
	localBatchStartButton         *widget.Button
	localBatchPolicySyncBox       *fyne.Container
	localBatchProgressBox         *fyne.Container
	localBatchCurrentFile         *widget.Label
	localBatchTimeRemaining       *widget.Label
	blcOnce                       sync.Once
)

var (
	partitionSyncComponent       fyne.CanvasObject
	partitionSyncPolicyComponent *widget.Button
	partitionProgressBox         *fyne.Container
	partitionPolicySyncBox       *fyne.Container
	partitionStartButton         *widget.Button
	partitionCurrentFile         *widget.Label
	partitionTimeRemaining       *widget.Label
	psOnce                       sync.Once
)

func getSingleLocalSyncComponent(win fyne.Window) fyne.CanvasObject {
	slcOnce.Do(func() {
		sourcePathBind := getBindString(config.SystemConfigCache.Value().LocalSingleSync.SourcePath)
		sourceContainer := loadValue2Label("Source: ", sourcePathBind)

		sourcePathComponent := container.New(layout.NewHBoxLayout(),
			sourceContainer,
			makeOpenFileBtn("Open", win, sourcePathBind, &config.SystemConfigCache.Cache.LocalSingleSync.SourcePath),
		)

		targetPathBind := getBindString(config.SystemConfigCache.Value().LocalSingleSync.TargetPath)
		localSingleTargetComponent := loadValue2Label("Target: ", targetPathBind)

		targetFileBtn := makeOpenFileBtn("Open", win, targetPathBind, &config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
		targetFolderBtn := makeOpenFolderBtn("Open", win, targetPathBind, &config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
		btnComponent := container.NewMax()
		btnComponent.Add(targetFileBtn)

		createFileComponent = widget.NewCheck("Create New File", func(b bool) {
			btnComponent.RemoveAll()
			if b {
				btnComponent.Add(targetFolderBtn)
			} else {
				btnComponent.Add(targetFileBtn)
			}
			btnComponent.Refresh()
		})

		targetPathComponent := container.New(layout.NewHBoxLayout(),
			localSingleTargetComponent,
			btnComponent,
			createFileComponent,
		)
		localSingleProgressBox = container.NewVBox()
		localSinglePolicySyncBox = container.NewVBox()
		localSingleCurrentFile = widget.NewLabel(NOT_RUNNING_STR)
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			localSingleCurrentFile,
		)
		localSingleTimeRemaining = widget.NewLabel(NOT_RUNNING_STR)
		lbTimeComp := container.NewHBox(
			widget.NewLabel("Time remaining:"),
			localSingleTimeRemaining,
		)
		localSingleSyncPolicyComponent = getSingleSyncPolicyBtn(win, false)

		localSinglePolicySyncBox = container.NewVBox()
		localSingleStartButton = getStartLocalSingleButton(win)
		localSingleSyncComponent = container.NewVBox(
			sourcePathComponent,
			targetPathComponent,
			currentFileComp,
			lbTimeComp,
			localSingleStartButton,
			localSingleSyncPolicyComponent,
			localSingleProgressBox,
		)
	})

	return localSingleSyncComponent
}

func getBatchLocalSyncComponent(win fyne.Window) fyne.CanvasObject {
	blcOnce.Do(func() {
		sourcePathBind := getBindString(config.SystemConfigCache.Value().LocalBatchSync.SourcePath)
		sourceContainer := loadValue2Label("Source: ", sourcePathBind)
		sourceComponent := container.New(layout.NewHBoxLayout(), sourceContainer, makeOpenFolderBtn("Open",
			win,
			sourcePathBind,
			&config.SystemConfigCache.Cache.LocalBatchSync.SourcePath))

		targetPathBind := getBindString(config.SystemConfigCache.Value().LocalBatchSync.TargetPath)
		targetContainer := loadValue2Label("Target: ", targetPathBind)
		targetComponent := container.New(layout.NewHBoxLayout(), targetContainer, makeOpenFolderBtn("Open",
			win,
			targetPathBind,
			&config.SystemConfigCache.Cache.LocalBatchSync.TargetPath))

		initStartLocalBatchButton(win)

		localBatchProgressBox = container.NewVBox()
		localBatchPolicySyncBox = container.NewVBox()
		localBatchSyncPolicyComponent = getBatchSyncPolicyBtn(win, false)

		localBatchCurrentFile = widget.NewLabel(NOT_RUNNING_STR)
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			localBatchCurrentFile,
		)
		localBatchTimeRemaining = widget.NewLabel(NOT_RUNNING_STR)
		lbTimeComp := container.NewHBox(
			widget.NewLabel("Time remaining:"),
			localBatchTimeRemaining,
		)
		FileSyncComponent = container.NewVBox(
			container.NewVBox(
				sourceComponent,
				targetComponent,
				currentFileComp,
				lbTimeComp,
			),
			localBatchStartButton,
			localBatchSyncPolicyComponent,
			localBatchProgressBox,
			localBatchPolicySyncBox,
		)
		localBatchSyncComponent = container.NewBorder(FileSyncComponent, nil, nil, nil)
	})
	return localBatchSyncComponent
}

func getPartitionSyncComponent(win fyne.Window) fyne.CanvasObject {
	psOnce.Do(func() {
		sPartitionComp, sPartitionSelect := getPartitionSelect("Source partition: ")
		tPartitionComp, tPartitionSelect := getPartitionSelect("Target partition: ")
		sPartitionSelect.OnChanged = func(s string) {
			config.SystemConfigCache.Cache.PartitionSync.SourcePath = s
			config.SystemConfigCache.NotifyAll()
		}
		sPartitionSelect.Selected = config.SystemConfigCache.Cache.PartitionSync.SourcePath
		tPartitionSelect.OnChanged = func(s string) {
			config.SystemConfigCache.Cache.PartitionSync.TargetPath = s
			config.SystemConfigCache.NotifyAll()
		}
		tPartitionSelect.Selected = config.SystemConfigCache.Cache.PartitionSync.TargetPath
		initPartitionSyncStartBtn(win)
		partitionCurrentFile = widget.NewLabel(NOT_RUNNING_STR)
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			partitionCurrentFile,
		)
		partitionTimeRemaining = widget.NewLabel(NOT_RUNNING_STR)
		partitionTimeComp := container.NewHBox(
			widget.NewLabel("Time remaining:"),
			partitionTimeRemaining,
		)
		partitionProgressBox = container.NewVBox()
		partitionPolicySyncBox = container.NewVBox()

		partitionSyncPolicyComponent = getPartitionSyncPolicyBtn(win)

		partitionSyncComponent = container.NewVBox(
			sPartitionComp,
			tPartitionComp,
			currentFileComp,
			partitionTimeComp,
			partitionStartButton,
			partitionSyncPolicyComponent,
			partitionProgressBox,
			partitionPolicySyncBox,
		)
	})
	return partitionSyncComponent
}

func initPartitionSyncStartBtn(win fyne.Window) {
	partitionStartButton = widget.NewButton("Start", func() {
		if config.SystemConfigCache.Cache.PartitionSync.TargetPath == config.NOT_SET_STR ||
			config.SystemConfigCache.Cache.PartitionSync.SourcePath == config.NOT_SET_STR {
			dialog.ShowInformation("Error", "Please set source and target path!", win)
			return
		}
		common.GetStartLock(common.TYPE_PARTITION).Add(1)
		//TODO 优化协程池
		go func() {
			log.Printf("PartitionSync Start Time : %v:%v:%v", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			worker.PartitionSyncSingleTime()
			log.Printf("PartitionSync Over Time : %v:%v:%v", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		}()
		common.GetStartLock(common.TYPE_PARTITION).Done()
	})
}

func initStartLocalBatchButton(win fyne.Window) {
	localBatchStartButton = widget.NewButton("Start", func() {
		if config.SystemConfigCache.Cache.LocalBatchSync.TargetPath == config.NOT_SET_STR ||
			config.SystemConfigCache.Cache.LocalBatchSync.SourcePath == config.NOT_SET_STR {
			dialog.ShowInformation("Error", "Please set source and target path!", win)
			return
		}
		common.GetStartLock(common.TYPE_LOCAL_BATCH).Add(1)
		//TODO 优化协程池
		go worker.LocalBatchSyncSingleTime(false)
		common.GetStartLock(common.TYPE_LOCAL_BATCH).Done()
	})
}

func getStartLocalSingleButton(win fyne.Window) *widget.Button {
	button := widget.NewButton("Start", func() {
		targetPath := config.SystemConfigCache.Cache.LocalSingleSync.TargetPath
		sourcePath := config.SystemConfigCache.Cache.LocalSingleSync.SourcePath
		if targetPath == config.NOT_SET_STR ||
			sourcePath == config.NOT_SET_STR {
			dialog.ShowInformation("Error", "Please set source and target path!", win)
			return
		}
		if createFileComponent.Checked {
			fileName := worker.GetFileName(sourcePath)
			config.SystemConfigCache.Cache.LocalSingleSync.TargetPath = targetPath + "/" + fileName
			createFileComponent.Checked = false
		}
		go worker.LocalSingleSyncSingleTime(false)
	})
	return button
}

func getFileTree() fyne.CanvasObject {
	dataM := make(map[string][]string)
	worker.FileNode2TreeMap(&dataM)
	tree := widget.NewTreeWithStrings(dataM)
	size := fyne.Size{
		Height: 600,
	}
	tree.Resize(size)
	return tree
}
