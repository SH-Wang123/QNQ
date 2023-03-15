package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
var localBatchPolicySyncBox *fyne.Container
var localBatchPolicySyncBar *fyne.Container
var localSinglePolicySyncBox *fyne.Container
var localSinglePolicySyncBar *fyne.Container

var (
	localSingleSyncComponent fyne.CanvasObject
	slcOnce                  sync.Once
)

var (
	localBatchSyncComponent fyne.CanvasObject
	localBatchStartButton   *widget.Button
	localBatchProgressBar   *fyne.Container
	localBatchProgressBox   *fyne.Container
	localBatchCurrentFile   *widget.Label
	blcOnce                 sync.Once
)

var (
	localBatchSyncPolicyComponent *widget.Button
	lbspOnce                      sync.Once
)

var (
	localSingleSyncPolicyComponent *widget.Button
	lsspOnce                       sync.Once
)

var (
	partitionSyncComponent fyne.CanvasObject
	partitionProgressBar   *fyne.Container
	partitionProgressBox   *fyne.Container
	partitionStartButton   *widget.Button
	partitionCurrentFile   *widget.Label
	psOnce                 sync.Once
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
		targetContainer := loadValue2Label("Target: ", targetPathBind)

		targetFileBtn := makeOpenFileBtn("Open", win, targetPathBind, &config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
		targetFolderBtn := makeOpenFolderBtn("Open", win, targetPathBind, &config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
		btnComponent := container.NewMax()
		btnComponent.Add(targetFolderBtn)

		createFileComponent := widget.NewCheck("Create New File", func(b bool) {
			if b {
				btnComponent.Remove(targetFileBtn)
				btnComponent.Add(targetFolderBtn)
			} else {
				btnComponent.Remove(targetFolderBtn)
				btnComponent.Add(targetFileBtn)
			}
		})

		targetPathComponent := container.New(layout.NewHBoxLayout(),
			targetContainer,
			btnComponent,
			createFileComponent,
		)
		localSingleSyncPolicyComponent = getSingleSyncPolicyBtn(win, false)

		syncPolicyRunningStatusComp := getPolicyStatusLabel(false, false)

		lspPro := widget.NewProgressBarInfinite()
		localSinglePolicySyncBar = container.NewVBox(lspPro)

		localSinglePolicySyncBox = container.NewVBox()

		localSingleSyncComponent = container.NewVBox(
			sourcePathComponent,
			targetPathComponent,
			syncPolicyRunningStatusComp,
			getStartLocalSingleButton(win),
			localSingleSyncPolicyComponent,
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

		syncPolicyRunningStatusComp := getPolicyStatusLabel(true, false)

		lbspPro := widget.NewProgressBarInfinite()
		localBatchPolicySyncBar = container.NewVBox(lbspPro)

		initStartLocalBatchButton(win)

		localBatchProgressBox = container.NewVBox()
		localBatchPolicySyncBox = container.NewVBox()
		localBatchSyncPolicyComponent = getBatchSyncPolicyBtn(win, false)

		localBatchCurrentFile = widget.NewLabel("Not running")
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			partitionCurrentFile,
		)

		FileSyncComponent = container.NewVBox(
			container.NewVBox(
				sourceComponent,
				targetComponent,
				syncPolicyRunningStatusComp,
				currentFileComp,
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
		partitionCurrentFile = widget.NewLabel("Not running")
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			partitionCurrentFile,
		)
		partitionProgressBox = container.NewVBox()
		partitionSyncComponent = container.NewVBox(
			sPartitionComp,
			tPartitionComp,
			currentFileComp,
			partitionStartButton,
			partitionProgressBox,
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
		common.LocalPartStartLock.Add(1)
		partitionProgressBar = getTaskProgressBar(partitionStartButton, partitionProgressBar, partitionProgressBox, true)
		partitionProgressBox.Add(partitionProgressBar)
		//TODO 优化协程池
		go func() {
			log.Printf("PartitionSync Start Time : %v:%v:%v", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			worker.StartPartitionSync()
			log.Printf("PartitionSync Over Time : %v:%v:%v", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		}()
		common.LocalPartStartLock.Done()
	})
}

func getPolicyStatusLabel(isBatch bool, isRemote bool) *fyne.Container {
	syncPolicy := config.SystemConfigCache.GetSyncPolicy(isBatch, isRemote)
	t := binding.NewBool()
	t.Set(syncPolicy.PolicySwitch)
	return container.New(layout.NewHBoxLayout(),
		widget.NewLabel("Sync Policy Running : "),
		widget.NewLabelWithData(binding.BoolToString(t)),
	)
}

func initStartLocalBatchButton(win fyne.Window) {
	localBatchStartButton = widget.NewButton("Start", func() {
		if config.SystemConfigCache.Cache.LocalBatchSync.TargetPath == config.NOT_SET_STR ||
			config.SystemConfigCache.Cache.LocalBatchSync.SourcePath == config.NOT_SET_STR {
			dialog.ShowInformation("Error", "Please set source and target path!", win)
			return
		}
		common.LocalBatchStartLock.Add(1)
		localBatchProgressBar = getTaskProgressBar(localBatchStartButton, localBatchProgressBar, localBatchProgressBox, false)
		localBatchProgressBox.Add(localBatchProgressBar)
		//TODO 优化协程池
		go worker.StartLocalBatchSync()
		common.LocalBatchStartLock.Done()
	})
}

func getStartLocalSingleButton(win fyne.Window) *widget.Button {
	button := widget.NewButton("Start", func() {
		if config.SystemConfigCache.Cache.LocalSingleSync.TargetPath == config.NOT_SET_STR ||
			config.SystemConfigCache.Cache.LocalSingleSync.SourcePath == config.NOT_SET_STR {
			dialog.ShowInformation("Error", "Please set source and target path!", win)
			return
		}
		go worker.StartLocalSingleSync()
	})
	return button
}

func getTaskProgressBar(startBtn *widget.Button, progressBar *fyne.Container, progressBox *fyne.Container, isPartition bool) *fyne.Container {
	startBtn.Disable()
	progress := widget.NewProgressBar()
	go func() {
		var progressNum = 0.0
		var clock *sync.WaitGroup
		if isPartition {
			clock = common.LocalPartStartLock
		} else {
			clock = common.LocalBatchStartLock
		}
		clock.Wait()
		for progressNum < 1 {
			time.Sleep(time.Millisecond * 100)
			if isPartition {
				partitionCurrentFile.SetText(common.CurrentLocalPartFile)
				progressNum = worker.GetLocalBatchProgress(common.CurrentLocalPartSN)
			} else {
				partitionCurrentFile.SetText(common.CurrentLocalBatchFile)
				progressNum = worker.GetLocalBatchProgress(common.CurrentLocalBatchSN)
			}
			progress.SetValue(progressNum)
			err := worker.GetBatchSyncError()
			if len(err) != 0 {
				log.Println()
			}
			partitionCurrentFile.Refresh()
		}
		partitionCurrentFile.SetText("Not running")
		time.Sleep(time.Second * 3)
		progressBox.Remove(progressBar)
		startBtn.Enable()
	}()
	return container.NewVBox(progress)
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
