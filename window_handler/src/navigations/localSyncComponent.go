package navigations

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"sync"
	"time"
	"window_handler/common"
	"window_handler/config"
	"window_handler/worker"
)

var FileSyncComponent fyne.CanvasObject
var progressBar *fyne.Container
var progressBox *fyne.Container
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

func GetSingleLocalSyncComponent(win fyne.Window) fyne.CanvasObject {
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
		localSingleSyncPolicyComponent = getSingleSyncPolicyBtn(win)

		syncPolicyRunningStatusComp := getPolicyStatusLabel(false)

		lspPro := widget.NewProgressBarInfinite()
		localSinglePolicySyncBar = container.NewVBox(lspPro)

		localSinglePolicySyncBox = container.NewVBox()

		localSingleSyncComponent = container.NewVBox(
			sourcePathComponent,
			targetPathComponent,
			syncPolicyRunningStatusComp,
			getStartLocalSingleButton(),
			localSingleSyncPolicyComponent,
		)
	})

	return localSingleSyncComponent
}

func GetBatchLocalSyncComponent(win fyne.Window) fyne.CanvasObject {
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

		syncPolicyRunningStatusComp := getPolicyStatusLabel(true)

		lbspPro := widget.NewProgressBarInfinite()
		localBatchPolicySyncBar = container.NewVBox(lbspPro)

		startButton := getStartLocalBatchButton()

		differecnButton := getDiffAnalysisButton()
		progressBox = container.NewVBox()
		localBatchPolicySyncBox = container.NewVBox()
		localBatchSyncPolicyComponent = getBatchSyncPolicyBtn(win)

		FileSyncComponent = container.NewVBox(
			container.NewVBox(
				sourceComponent,
				targetComponent,
				syncPolicyRunningStatusComp,
			),
			startButton,
			differecnButton,
			localBatchSyncPolicyComponent,
			progressBox,
			localBatchPolicySyncBox,
		)
		localBatchSyncComponent = container.NewBorder(FileSyncComponent, nil, getFileTree(), getFileTree())
	})
	//ProgressBar.Hide()
	go watchLocalSyncPolicy()
	return localBatchSyncComponent
}

func getPolicyStatusLabel(isBatch bool) *fyne.Container {
	syncPolicyRunningStatus := binding.NewBool()
	if isBatch {
		syncPolicyRunningStatus.Set(worker.BatchSyncPolicyRunFlag)
	} else {
		syncPolicyRunningStatus.Set(worker.SingleSyncPolicyRunFlag)
	}
	return container.New(layout.NewHBoxLayout(),
		widget.NewLabel("Sync Policy Running : "),
		widget.NewLabelWithData(binding.BoolToString(syncPolicyRunningStatus)),
	)
}

func getStartLocalBatchButton() *widget.Button {
	button := widget.NewButton("Start", func() {
		progressBar = getTaskProgressBar()
		progressBox.Add(progressBar)
		//TODO 优化协程池
		go worker.SyncBatchFileTree(worker.LocalBSFileNode, config.SystemConfigCache.Value().LocalBatchSync.TargetPath)

	})
	return button
}

func getStartLocalSingleButton() *widget.Button {
	button := widget.NewButton("Start", func() {
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

func getTaskProgressBar() *fyne.Container {
	progress := widget.NewProgressBar()
	go func() {
		var progressNum = 0.0
		for progressNum < 1 {
			time.Sleep(time.Millisecond * 500)
			progressNum = worker.GetLocalBatchProgress()
			progress.SetValue(progressNum)
		}
		time.Sleep(time.Second * 10)
		progressBox.Remove(progressBar)

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

func getBatchSyncPolicyBtn(win fyne.Window) *widget.Button {
	lbspOnce.Do(func() {
		localBatchSyncPolicyComponent = getSyncPolicyBtn(true, win)
	})
	return localBatchSyncPolicyComponent
}

func getSingleSyncPolicyBtn(win fyne.Window) *widget.Button {
	lsspOnce.Do(func() {
		localSingleSyncPolicyComponent = getSyncPolicyBtn(false, win)
	})
	return localSingleSyncPolicyComponent
}

func getSyncPolicyBtn(isBatchSync bool, win fyne.Window) *widget.Button {
	return widget.NewButton("Sync Policy", func() {
		var title string
		var node *worker.FileNode
		configCache := config.SystemConfigCache.GetLocalPeriodicSyncPolicy(isBatchSync)
		rateList := make([]string, 0)
		cycleList := make([]string, 0)
		if isBatchSync {
			title = "Local batch sync policy"
			node = worker.LocalBSFileNode
		} else {
			title = "Local single sync policy"
			node = worker.GetSingleFileNode(configCache.SourcePath)
		}
		for i := 1; i <= 60; i++ {
			rateList = append(rateList, fmt.Sprintf("%d", i))
		}
		rateSelect := widget.NewSelect(rateList, nil)
		rateSelect.SetSelected(fmt.Sprintf("%d", configCache.PeriodicSync.Rate))
		var cycleStr string
		for k, v := range timeCycleMap {
			cycleList = append(cycleList, k)
			if v == configCache.PeriodicSync.Cycle {
				cycleStr = k
			}
		}
		cycleSelect := widget.NewSelect(cycleList, nil)
		cycleSelect.SetSelected(cycleStr)
		rateAndCycleComponent := container.NewHBox(
			rateSelect,
			cycleSelect,
		)

		enableCheck := widget.NewCheck("Enable", nil)
		enableCheck.Checked = configCache.PeriodicSync.Enable
		items := []*widget.FormItem{
			widget.NewFormItem("Sync cycle: ", rateAndCycleComponent),
			widget.NewFormItem("", enableCheck),
		}

		dialog.ShowForm(title, "Save & Start", "Cancel", items, func(b bool) {
			if b {
				configCache.PeriodicSync.Rate, _ = strconv.Atoi(rateSelect.Selected)
				configCache.PeriodicSync.Cycle = timeCycleMap[cycleSelect.Selected]
				configCache.PeriodicSync.Enable = enableCheck.Checked
				config.SystemConfigCache.NotifyAll()
				tem := false
				if configCache.PeriodicSync.Enable {
					worker.StartPeriodicSync(node,
						configCache.TargetPath,
						time.Duration(configCache.PeriodicSync.Rate)*configCache.PeriodicSync.Cycle,
						&tem,
						isBatchSync,
					)
				}
			}
		}, win)
	})
}

func watchLocalSyncPolicy() {
	for {
		select {
		case c := <-common.GWChannel:
			if c == common.LOCAL_BATCH_POLICY_RUNNING {
				localBatchPolicySyncBox.Add(localBatchPolicySyncBar)
				localBatchSyncComponent.Refresh()
				var progressNum = 0.0
				for progressNum < 1 {
					time.Sleep(time.Millisecond * 500)
					progressNum = worker.GetLocalBatchProgress()
				}
			} else if c == common.LOCAL_BATCH_POLICY_STOP {
				localBatchPolicySyncBox.Remove(localBatchPolicySyncBar)
			} else if c == common.LOCAL_SINGLE_POLICY_RUNNING {
				localSinglePolicySyncBox.Add(localSinglePolicySyncBar)
				localSingleSyncComponent.Refresh()
			} else if c == common.LOCAL_SINGLE_POLICY_STOP {
				localSinglePolicySyncBox.Remove(localSinglePolicySyncBar)
			}

		}
	}
}
