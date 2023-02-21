package navigations

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
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
	localBatchStartButton   *widget.Button
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

		localBatchStartButton = getStartLocalBatchButton()

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
			localBatchStartButton,
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
	localBatchStartButton.Disable()
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
		localBatchStartButton.Enable()
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
	return getSyncPolicyBtn(true, win)
}

func getSingleSyncPolicyBtn(win fyne.Window) *widget.Button {
	return getSyncPolicyBtn(false, win)
}

func getSyncPolicyBtn(isBatchSync bool, win fyne.Window) *widget.Button {
	return widget.NewButton("Sync Policy", func() {
		var title string
		var node *worker.FileNode
		var daysCheckCompent [7]*widget.Check
		var rateSelectedValue string
		var disableCacheKey string
		var usePeriodicSyncCheck *widget.Check
		var useTimingSyncCheck *widget.Check
		var policyEnableCheck *widget.Check
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
		disableCacheKey = title
		clearDisableRootCache(disableCacheKey)
		//Periodic sync
		for i := 1; i <= 60; i++ {
			rateList = append(rateList, fmt.Sprintf("%d", i))
		}
		rateSelect := widget.NewSelect(rateList, nil)
		for k, v := range timeCycleMap {
			cycleList = append(cycleList, k)
			if v == configCache.PeriodicSync.Cycle {
				rateSelectedValue = k
			}
		}
		cycleSelect := widget.NewSelect(cycleList, nil)
		rateAndCycleComponent := container.NewHBox(
			rateSelect,
			cycleSelect,
		)

		//Timing sync
		daysContainer := container.NewGridWithColumns(7)
		for index := 0; index < len(daysCheckCompent); index++ {
			daysCheckCompent[index] = widget.NewCheck(dayArrayList[index], nil)
			daysContainer.Add(daysCheckCompent[index])
		}

		usePeriodicSyncCheck = widget.NewCheck("Used periodic sync", func(b bool) {
			if b {
				enableAllChild(disableCacheKey, usePeriodicSyncCheck)
			} else {
				disableAllChild(disableCacheKey, usePeriodicSyncCheck)
			}
		})
		addDisableRoot(disableCacheKey, usePeriodicSyncCheck, rateSelect, cycleSelect)

		useTimingSyncCheck = widget.NewCheck("Used timing sync", func(b bool) {
			if b {
				enableAllChild(disableCacheKey, useTimingSyncCheck)
			} else {
				disableAllChild(disableCacheKey, useTimingSyncCheck)
			}
		})
		addDisableRoot(disableCacheKey, useTimingSyncCheck, daysCheckCompent[0], daysCheckCompent[1], daysCheckCompent[2], daysCheckCompent[3],
			daysCheckCompent[4], daysCheckCompent[5], daysCheckCompent[6])

		policyEnableCheck = widget.NewCheck("Global switch", func(b bool) {
			swapChecked(usePeriodicSyncCheck)
			swapChecked(useTimingSyncCheck)
			if b {
				enableAllChild(disableCacheKey, policyEnableCheck)
			} else {
				disableAllChild(disableCacheKey, policyEnableCheck)
			}
		})
		addDisableRoot(disableCacheKey, policyEnableCheck, usePeriodicSyncCheck, useTimingSyncCheck)
		items := []*widget.FormItem{
			widget.NewFormItem("Select: ", useTimingSyncCheck),
			widget.NewFormItem("Time:  ", daysContainer),
			widget.NewFormItem("Select: ", usePeriodicSyncCheck),
			widget.NewFormItem("Sync cycle: ", rateAndCycleComponent),
			widget.NewFormItem("Select: ", policyEnableCheck),
		}

		dialog.ShowForm(title, "Save & Start", "Cancel", items, func(b bool) {
			if b {
				configCache.PeriodicSync.Rate, _ = strconv.Atoi(rateSelect.Selected)
				configCache.PeriodicSync.Cycle = timeCycleMap[cycleSelect.Selected]
				configCache.PeriodicSync.Enable = usePeriodicSyncCheck.Checked
				configCache.TimingSync.Enable = useTimingSyncCheck.Checked
				configCache.PolicySwitch = policyEnableCheck.Checked
				for index := 0; index < len(daysCheckCompent); index++ {
					configCache.TimingSync.Days[index] = daysCheckCompent[index].Checked
				}
				config.SystemConfigCache.NotifyAll()
				tem := false
				if configCache.PolicySwitch {
					worker.StartPeriodicSync(node,
						configCache.TargetPath,
						time.Duration(configCache.PeriodicSync.Rate)*configCache.PeriodicSync.Cycle,
						&tem,
						isBatchSync,
					)
				}
			}
		}, win)
		//init value
		cycleSelect.SetSelected(rateSelectedValue)
		rateSelect.SetSelected(fmt.Sprintf("%d", configCache.PeriodicSync.Rate))

		usePeriodicSyncCheck.Checked = configCache.PeriodicSync.Enable
		useTimingSyncCheck.Checked = configCache.TimingSync.Enable
		policyEnableCheck.Checked = configCache.PolicySwitch

		if !configCache.PeriodicSync.Enable {
			disableAllChild(disableCacheKey, usePeriodicSyncCheck)
		}
		if !configCache.TimingSync.Enable {
			disableAllChild(disableCacheKey, useTimingSyncCheck)
		}
		if !policyEnableCheck.Checked {
			disableAllChild(disableCacheKey, policyEnableCheck)
		}

		for index := 0; index < len(daysCheckCompent); index++ {
			daysCheckCompent[index].SetChecked(configCache.TimingSync.Days[index])
		}

		batchRefresh(usePeriodicSyncCheck, useTimingSyncCheck, policyEnableCheck, cycleSelect, rateSelect)
	})
}

func watchLocalSyncPolicy() {
	for {
		select {
		case c := <-common.GWChannel:
			if c == common.LOCAL_BATCH_POLICY_RUNNING {
				localBatchPolicySyncBox.Add(localBatchPolicySyncBar)
				batchDisable(localBatchSyncPolicyComponent, localBatchStartButton)
				localBatchSyncComponent.Refresh()
			} else if c == common.LOCAL_BATCH_POLICY_STOP {
				batchEnable(localBatchSyncPolicyComponent, localBatchStartButton)
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
