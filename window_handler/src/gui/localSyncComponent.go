package gui

import (
	"fmt"
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
	localSinglePolicySyncBox       *fyne.Container
	localSinglePolicySyncPro       *fyne.Container
	slcOnce                        sync.Once
)

var (
	localBatchSyncComponent       fyne.CanvasObject
	localBatchSyncPolicyComponent *widget.Button
	localBatchStartButton         *widget.Button
	localBatchPolicySyncBox       *fyne.Container
	localBatchPolicySyncPro       *fyne.Container
	localBatchProgressBar         *fyne.Container
	localBatchProgressBox         *fyne.Container
	localBatchCurrentFile         *widget.Label
	localBatchTimeRemaining       *widget.Label
	blcOnce                       sync.Once
)

var (
	partitionSyncComponent       fyne.CanvasObject
	partitionSyncPolicyComponent *widget.Button
	partitionProgressBar         *fyne.Container
	partitionProgressBox         *fyne.Container
	partitionPolicySyncBox       *fyne.Container
	partitionPolicySyncPro       *fyne.Container
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

		lspPro := widget.NewProgressBarInfinite()
		localSinglePolicySyncPro = container.NewVBox(lspPro)

		localSinglePolicySyncBox = container.NewVBox()

		localSingleSyncComponent = container.NewVBox(
			sourcePathComponent,
			targetPathComponent,
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

		lbspPro := widget.NewProgressBarInfinite()
		localBatchPolicySyncPro = container.NewVBox(lbspPro)

		initStartLocalBatchButton(win)

		localBatchProgressBox = container.NewVBox()
		localBatchPolicySyncBox = container.NewVBox()
		localBatchSyncPolicyComponent = getBatchSyncPolicyBtn(win, false)

		localBatchCurrentFile = widget.NewLabel("Not running")
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			localBatchCurrentFile,
		)
		localBatchTimeRemaining = widget.NewLabel("Not running")
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
		partitionCurrentFile = widget.NewLabel("Not running")
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			partitionCurrentFile,
		)
		partitionTimeRemaining = widget.NewLabel("Not running")
		partitionTimeComp := container.NewHBox(
			widget.NewLabel("Time remaining:"),
			partitionTimeRemaining,
		)
		partitionProgressBox = container.NewVBox()
		partitionPolicySyncBox = container.NewVBox()

		partitionSyncPolicyComponent = getPartitionSyncPolicyBtn(win)

		lbspPro := widget.NewProgressBarInfinite()
		partitionPolicySyncPro = container.NewVBox(lbspPro)

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
		common.LocalPartStartLock.Add(1)
		partitionProgressBar = getTaskProgressBar(partitionStartButton, partitionProgressBar, partitionProgressBox, true)
		//TODO 优化协程池
		go func() {
			log.Printf("PartitionSync Start Time : %v:%v:%v", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			worker.PartitionSyncSingleTime()
			log.Printf("PartitionSync Over Time : %v:%v:%v", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		}()
		common.LocalPartStartLock.Done()
	})
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
		//TODO 优化协程池
		go worker.LocalBatchSyncSingleTime()
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
		go worker.LocalSingleSyncSingleTime()
	})
	return button
}

func getTaskProgressBar(startBtn *widget.Button, progressBar *fyne.Container, progressBox *fyne.Container, isPartition bool) *fyne.Container {
	startBtn.Disable()
	progress := widget.NewProgressBar()
	progressBox.Add(progress)
	go func() {
		time.Sleep(time.Millisecond * 500)
		var progressNum = 0.0
		var currentFileLabel *widget.Label
		var currentSN string
		var currentTimeRemaining *widget.Label
		if isPartition {
			common.LocalPartStartLock.Wait()
			currentFileLabel = partitionCurrentFile
			currentTimeRemaining = partitionTimeRemaining
			currentSN = common.CurrentLocalPartSN
		} else {
			common.LocalBatchStartLock.Wait()
			currentFileLabel = localBatchCurrentFile
			currentTimeRemaining = localBatchTimeRemaining
			currentSN = common.CurrentLocalBatchSN
		}
		go func() {
			currentTimeRemaining.SetText("Under calculation")
			for {
				remaining := worker.EstimatedTotalTime(currentSN, 2*time.Second)
				if remaining <= 0 {
					return
				}
				currentTimeRemaining.SetText(fmt.Sprint(remaining))
			}
		}()
		for progressNum < 1 {
			time.Sleep(time.Millisecond * 100)
			currentFileLabel.SetText(common.GetCurrentSyncFile(currentSN))
			progressNum = worker.GetLocalBatchProgress(currentSN)
			currentFileLabel.Refresh()
			progress.SetValue(progressNum)
			err := worker.GetBatchSyncError()
			if len(err) != 0 {
				log.Println()
			}
		}
		currentFileLabel.SetText("Not running")
		currentTimeRemaining.SetText("Not running")
		time.Sleep(time.Second * 3)
		progressBox.Remove(progress)
		progressBox.Refresh()
		startBtn.Enable()
		currentFileLabel.Refresh()
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
