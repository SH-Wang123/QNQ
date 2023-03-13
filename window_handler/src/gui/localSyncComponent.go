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
	diffAnalysisButton      *widget.Button
	localBatchProgressBar   *fyne.Container
	localBatchProgressBox   *fyne.Container
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

		diffAnalysisButton = getDiffAnalysisButton()
		localBatchProgressBox = container.NewVBox()
		localBatchPolicySyncBox = container.NewVBox()
		localBatchSyncPolicyComponent = getBatchSyncPolicyBtn(win, false)

		FileSyncComponent = container.NewVBox(
			container.NewVBox(
				sourceComponent,
				targetComponent,
				syncPolicyRunningStatusComp,
			),
			localBatchStartButton,
			diffAnalysisButton,
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
		tPartitionComp, tPartitionSelect := getPartitionSelect("Source partition: ")
		sPartitionSelect.OnChanged = func(s string) {
			config.SystemConfigCache.Cache.PartitionSync.SourcePath = s
			config.SystemConfigCache.NotifyAll()
		}
		tPartitionSelect.OnChanged = func(s string) {
			config.SystemConfigCache.Cache.PartitionSync.TargetPath = s
			config.SystemConfigCache.NotifyAll()
		}
		initPartitionSyncStartBtn(win)
		partitionProgressBox = container.NewVBox()
		partitionSyncComponent = container.NewVBox(
			sPartitionComp,
			tPartitionComp,
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
		partitionProgressBar = getTaskProgressBar(partitionStartButton, partitionProgressBar, partitionProgressBox)
		partitionProgressBox.Add(partitionProgressBar)
		//TODO 优化协程池
		go func() {
			sourceNode := worker.GetNilNode(config.SystemConfigCache.Cache.PartitionSync.SourcePath)
			worker.GetFileTree(sourceNode, true)
			worker.SyncBatchFileTree(worker.LocalBSFileNode, config.SystemConfigCache.Value().PartitionSync.TargetPath)
		}()

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
		localBatchProgressBar = getTaskProgressBar(localBatchStartButton, localBatchProgressBar, localBatchProgressBox)
		localBatchProgressBox.Add(localBatchProgressBar)
		//TODO 优化协程池
		go worker.SyncBatchFileTree(worker.LocalBSFileNode, config.SystemConfigCache.Value().LocalBatchSync.TargetPath)

	})
}

func getStartLocalSingleButton(win fyne.Window) *widget.Button {
	button := widget.NewButton("Start", func() {
		if config.SystemConfigCache.Cache.LocalSingleSync.TargetPath == config.NOT_SET_STR ||
			config.SystemConfigCache.Cache.LocalSingleSync.SourcePath == config.NOT_SET_STR {
			dialog.ShowInformation("Error", "Please set source and target path!", win)
			return
		}
		go worker.LocalSyncSingleFileGUI()
	})
	return button
}

func getDiffAnalysisButton() *widget.Button {
	button := widget.NewButton("Variance Analysis", func() {
		worker.MarkFileTree(worker.LocalBSFileNode, config.SystemConfigCache.Value().LocalBatchSync.TargetPath)
	})
	return button
}

func getTaskProgressBar(startBtn *widget.Button, progressBar *fyne.Container, progressBox *fyne.Container) *fyne.Container {
	startBtn.Disable()
	progress := widget.NewProgressBar()
	go func() {
		var progressNum = 0.0
		for progressNum < 1 {
			time.Sleep(time.Millisecond * 500)
			progressNum = worker.GetLocalBatchProgress()
			progress.SetValue(progressNum)
			err := worker.GetBatchSyncError()
			if len(err) != 0 {
				log.Println()
			}
		}
		time.Sleep(time.Second * 1)
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
