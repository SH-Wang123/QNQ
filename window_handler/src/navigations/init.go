package navigations

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"time"
	"window_handler/common"
	"window_handler/worker"
)

func init() {
	I18n()
	go watchGWChannel()
}

func I18n() {
	initTimeCycle()
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

func watchGWChannel() {
	for {
		select {
		case c := <-common.GWChannel:
			if c == common.LOCAL_BATCH_POLICY_RUNNING {
				if common.LocalBatchPolicyRunningFlag {
					syncErrorDialogOK = true
					if !syncErrorDialogOK {
						syncErrorDialog.Show()
					}
					continue
				}
				common.LocalBatchPolicyRunningFlag = true
				localBatchPolicySyncBox.Add(localBatchPolicySyncBar)
				batchDisable(localBatchSyncPolicyComponent, localBatchStartButton, diffAnalysisButton)
				localBatchSyncComponent.Refresh()
			} else if c == common.LOCAL_BATCH_POLICY_STOP {
				batchEnable(localBatchSyncPolicyComponent, localBatchStartButton, diffAnalysisButton)
				common.LocalBatchPolicyRunningFlag = false
				syncErrorDialogOK = false
				localBatchPolicySyncBox.Remove(localBatchPolicySyncBar)
			} else if c == common.LOCAL_SINGLE_POLICY_RUNNING {
				localSinglePolicySyncBox.Add(localSinglePolicySyncBar)
				localSingleSyncComponent.Refresh()
			} else if c == common.LOCAL_SINGLE_POLICY_STOP {
				localSinglePolicySyncBox.Remove(localSinglePolicySyncBar)
			} else if c == common.TEST_DISK_SPEED_START {
				testSpeedRetLab.SetText("Testing...")
			} else if c == common.TEST_DISK_SPEED_OVER {
				setDiskSpeedRet()
			}

		}
	}
}

func setDiskSpeedRet() {
	partition := partitionSelect.Selected
	rSpeed := fmt.Sprint(worker.DiskReadSpeedCache[partition])
	wSpeed := fmt.Sprint(worker.DiskWriteSpeedCache[partition])
	testSpeedRetLab.SetText("Disk : " + partition + "\n" + "Read speed : " + rSpeed + "MB/S\n" + "Write speed : " + wSpeed + "MB/S\n")
}
