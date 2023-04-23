package gui

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
	"time"
	"unsafe"
	"window_handler/common"
	"window_handler/config"
	"window_handler/worker"
)

func bindingValue2Label(text string, value interface{}) *fyne.Container {
	var valueLabel *widget.Label
	unsafePtr := unsafe.Pointer(&value)
	switch value.(type) {
	case string:
		data := binding.NewString()
		data.Set(*(*string)(unsafePtr))
		valueLabel = widget.NewLabelWithData(data)
	case int:
		data := binding.NewInt()
		data.Set(*(*int)(unsafePtr))
		valueLabel = widget.NewLabelWithData(binding.IntToString(data))
	case float64:
		data := binding.NewFloat()
		data.Set(*(*float64)(unsafePtr))
		valueLabel = widget.NewLabelWithData(binding.FloatToString(data))
	case bool:
		data := binding.NewBool()
		data.Set(*(*bool)(unsafePtr))
		valueLabel = widget.NewLabelWithData(binding.BoolToString(data))
	}
	return container.New(layout.NewHBoxLayout(),
		widget.NewLabel(text),
		valueLabel,
	)
}

func loadValue2Label(text string, bindValue binding.String) *fyne.Container {
	return container.New(layout.NewHBoxLayout(),
		widget.NewLabel(text),
		widget.NewLabelWithData(bindValue),
	)
}

func getLabelSelect(text string, selectCom *widget.Select) *fyne.Container {
	return container.NewHBox(widget.NewLabel(text), selectCom)
}

func makeOpenFolderBtn(buttonName string, win fyne.Window, bindPath binding.String, configStr *string) *widget.Button {
	return widget.NewButton(buttonName, func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri == nil {
				return
			}
			_, err = uri.List()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			bindPath.Set(uri.Path())
			*configStr = uri.Path()
			config.SystemConfigCache.NotifyAll()
		}, win)
	})
}

func makeBtn(buttonName string, clickFunc func()) *widget.Button {
	return widget.NewButton(buttonName, clickFunc)
}

func makeOpenFileBtn(buttonName string, win fyne.Window, bindPath binding.String, configStr *string) *widget.Button {
	return widget.NewButton(buttonName, func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri == nil {
				log.Println("Cancelled")
				return
			}
			bindPath.Set(uri.URI().Path())
			*configStr = uri.URI().Path()
			config.SystemConfigCache.NotifyAll()
		}, win)
	})
}

func getBindString(value string) binding.String {
	bindD := binding.NewString()
	bindD.Set(value)
	return bindD
}

func batchDisable(ds ...fyne.Disableable) {
	for _, d := range ds {
		d.Disable()
	}
}

func batchEnable(ds ...fyne.Disableable) {
	for _, d := range ds {
		d.Enable()
	}
}

func batchRefresh(bs ...fyne.CanvasObject) {
	for _, v := range bs {
		v.Refresh()
	}
}

func clearDisableRootCache(key string) {
	disableRootCache[key] = make(map[fyne.Disableable]disableRoot)
}

func addDisableRoot(key string, root fyne.Disableable, chidls ...fyne.Disableable) {
	currentRoot := disableRoot{}
	currentRoot.addChild(chidls...)
	disableRootCache[key][root] = currentRoot
}

func disableAllChild(key string, root fyne.Disableable) {
	cache := disableRootCache[key]
	currentRoot := cache[root]
	currentRoot.disableChild()
	for _, v := range currentRoot.child {
		disableAllChild(key, v)
	}
}

func enableAllChild(key string, root fyne.Disableable) {
	cache := disableRootCache[key]
	currentRoot := cache[root]
	currentRoot.enableChild()
}

// swapChecked activation check function
func swapChecked(w *widget.Check) {
	w.SetChecked(!w.Checked)
	w.SetChecked(!w.Checked)
}

func getBatchSyncPolicyBtn(win fyne.Window, isRemote bool) *widget.Button {
	if isRemote {
		return getSyncPolicyBtn(true, false, true, win)
	}
	return getSyncPolicyBtn(true, false, false, win)
}

func getPartitionSyncPolicyBtn(win fyne.Window) *widget.Button {
	return getSyncPolicyBtn(false, false, true, win)
}

func getSingleSyncPolicyBtn(win fyne.Window, isRemote bool) *widget.Button {
	if isRemote {
		return getSyncPolicyBtn(false, true, false, win)
	}
	return getSyncPolicyBtn(false, false, false, win)
}

func getSyncPolicyBtn(isBatchSync bool, isRemoteSync bool, isPartitionSync bool, win fyne.Window) *widget.Button {
	return widget.NewButton("Sync Policy", func() {
		var title string
		var daysCheckCompent [7]*widget.Check
		var rateSelectedValue string
		var disableCacheKey string
		var usePeriodicSyncCheck *widget.Check
		var useTimingSyncCheck *widget.Check
		var policyEnableCheck *widget.Check
		configCache := config.SystemConfigCache.GetSyncPolicy(isBatchSync, isRemoteSync, isPartitionSync)
		rateList := make([]string, 0)
		cycleList := make([]string, 0)
		if isPartitionSync {
			title = "Partition sync policy"
		} else if isBatchSync {
			if isRemoteSync {
				title = "Remote batch sync policy"
			} else {
				title = "Local batch sync policy"
			}
		} else {
			if isRemoteSync {
				title = "Remote single sync policy"
			} else {
				title = "Local single sync policy"
			}
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
		indexArray := make([]string, 0)
		for i := 0; i <= 24; i++ {
			s := fmt.Sprintf("%d", i)
			indexArray = append(indexArray, s)
		}
		hourSelect := widget.NewSelect(indexArray, nil)
		for i := 25; i <= 59; i++ {
			indexArray = append(indexArray, fmt.Sprint(i))
		}
		minSelect := widget.NewSelect(indexArray, nil)
		timeContainer := container.NewHBox(
			widget.NewLabel("Hour:"),
			hourSelect,
			widget.NewLabel("Minute:"),
			minSelect,
		)

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
			daysCheckCompent[4], daysCheckCompent[5], daysCheckCompent[6], minSelect, hourSelect)

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
			widget.NewFormItem("Day:  ", daysContainer),
			widget.NewFormItem("", timeContainer),
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
				minute, _ := strconv.Atoi(minSelect.Selected)
				hour, _ := strconv.Atoi(hourSelect.Selected)
				configCache.TimingSync.Minute = uint8(minute)
				configCache.TimingSync.Hour = uint8(hour)
				for index := 0; index < len(daysCheckCompent); index++ {
					configCache.TimingSync.Days[index] = daysCheckCompent[index].Checked
				}
				config.SystemConfigCache.NotifyAll()
				tem := false
				if configCache.PolicySwitch {
					if configCache.PeriodicSync.Enable {
						worker.StartPolicySync(
							time.Duration(configCache.PeriodicSync.Rate)*configCache.PeriodicSync.Cycle,
							&tem,
							isBatchSync,
							isRemoteSync,
							isPartitionSync,
							true,
						)
					}
					if configCache.TimingSync.Enable {
						nextTime := worker.GetNextSyncTime(
							configCache.TimingSync.Days,
							configCache.TimingSync.Minute,
							configCache.TimingSync.Hour,
						)
						worker.StartPolicySync(
							nextTime,
							&tem,
							isBatchSync,
							isRemoteSync,
							isPartitionSync,
							false,
						)
					}
				}
			}
		}, win)
		//init value
		cycleSelect.SetSelected(rateSelectedValue)
		rateSelect.SetSelected(fmt.Sprintf("%d", configCache.PeriodicSync.Rate))
		minSelect.SetSelected(fmt.Sprintf("%d", configCache.TimingSync.Minute))
		hourSelect.SetSelected(fmt.Sprintf("%d", configCache.TimingSync.Hour))

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

		batchRefresh(usePeriodicSyncCheck, useTimingSyncCheck, policyEnableCheck, cycleSelect, rateSelect, minSelect, hourSelect)
	})
}

func getPartitionSelect(text string) (*fyne.Container, *widget.Select) {
	var partitions []string
	for _, v := range worker.DiskPartitionsCache {
		partitions = append(partitions, v.Name)
	}
	partitionSelect := widget.NewSelect(partitions, nil)
	return getLabelSelect(text, partitionSelect), partitionSelect
}

// startSyncGUI 同步进度条、当前文件、剩余时间显示
func startSyncGUI(progressBox *fyne.Container, currentFileLabel *widget.Label, currentTimeRemaining *widget.Label, businessType int) {
	progress := widget.NewProgressBar()
	progressBox.Add(progress)
	go func() {
		time.Sleep(time.Millisecond * 500)
		var progressNum = 0.0
		currentSN := common.GetCurrentSN(businessType)
		common.GetStartLock(businessType).Wait()

		go func() {
			currentTimeRemaining.SetText("Under calculation")
			for {
				remaining := worker.EstimatedTotalTime(currentSN, 10*time.Second)
				if remaining < 0 {
					return
				} else if remaining == 0 {
					currentTimeRemaining.SetText("Verify MD5, Please wait...")
					return
				}
				if !common.GetRunningFlag(businessType) {
					return
				}
				currentTimeRemaining.SetText(fmt.Sprint(remaining))
			}
		}()

		for progressNum < 1 {
			if currentSN == "" {
				currentSN = common.GetCurrentSN(businessType)
			}
			if !common.GetRunningFlag(businessType) {
				return
			}
			time.Sleep(time.Millisecond * 100)
			currentFileLabel.SetText(common.GetCurrentSyncFile(currentSN))
			progressNum = worker.GetLocalBatchProgress(currentSN)
			currentFileLabel.Refresh()
			progress.SetValue(progressNum)
			err := worker.GetBatchSyncError(currentSN)
			if len(err) != 0 {
				//log.Println()
			}
		}
		progressBox.Refresh()
		currentFileLabel.Refresh()
	}()
}

// overSyncGUI 同步进度条、当前文件、剩余时间隐藏
func overSyncGUI(progressBox *fyne.Container, currentFileLabel *widget.Label, currentTimeRemaining *widget.Label) {
	time.Sleep(time.Second * 1)
	progressBox.RemoveAll()
	currentFileLabel.SetText(NOT_RUNNING_STR)
	currentFileLabel.Refresh()
	currentTimeRemaining.SetText(NOT_RUNNING_STR)
	currentTimeRemaining.Refresh()
}

func getCancelButton(busType int) *widget.Button {
	but := widget.NewButton("Cancel", func() {
		worker.CancelTask(busType)
	})
	but.Disable()
	return but
}
