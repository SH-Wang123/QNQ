package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sync"
	"window_handler/common"
	"window_handler/config"
	"window_handler/worker"
)

var (
	remoteAuthDialog  *dialog.ConfirmDialog
	remoteAuthShowing = false
)

var (
	remoteSyncComponent fyne.CanvasObject
	rscOnce             sync.Once
)

func getManageRemoteQNQComponent(win fyne.Window) fyne.CanvasObject {
	items := make([]*container.TabItem, 0)
	itemsLen := len(config.SystemConfigCache.Cache.QNQNetCells)
	for index := 0; index < itemsLen; index++ {
		item := createRemoteQNQItem(index)
		items = append(items, item)
	}
	tabs := container.NewDocTabs(items...)
	tabs.OnClosed = func(item *container.TabItem) {
		i := 1
		fmt.Sprint(i)
	}
	tabs.CreateTab = func() *container.TabItem {
		config.SystemConfigCache.Cache.AddNilNetCell()
		itemsLen = len(config.SystemConfigCache.Cache.QNQNetCells) - 1
		return createRemoteQNQItem(itemsLen)
	}

	tabs.SetTabLocation(container.TabLocationLeading)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func createRemoteQNQItem(index int) *container.TabItem {
	ipComp, ipInput := getLabelInput("IP : ", config.SystemConfigCache.Cache.QNQNetCells[index].Ip)
	serverStatusComp := container.NewHBox(
		widget.NewLabel("Server status : "),
		widget.NewLabel(fmt.Sprint(config.SystemConfigCache.Cache.QNQNetCells[index].GetServerStatus())),
	)
	targetStatusComp := container.NewHBox(
		widget.NewLabel("Target status : "),
		widget.NewLabel(fmt.Sprint(config.SystemConfigCache.Cache.QNQNetCells[index].GetTargetStatus())),
	)
	markComp, markInput := getLabelInput("Mark : ", config.SystemConfigCache.Cache.QNQNetCells[index].Mark)
	saveBtn := widget.NewButton("Save", func() {
		config.SystemConfigCache.Cache.QNQNetCells[index].Ip = ipInput.Text
		config.SystemConfigCache.Cache.QNQNetCells[index].Mark = markInput.Text
		config.SystemConfigCache.NotifyAll()
	})
	authBtn := widget.NewButton("Auth", func() {
		go worker.ConnectTarget(ipInput.Text)
	})
	comp := container.NewVBox(
		ipComp,
		serverStatusComp,
		targetStatusComp,
		markComp,
		authBtn,
		saveBtn,
	)
	return container.NewTabItem(
		fmt.Sprint(index+1),
		comp,
	)
}

func getRemoteSingleComponent(win fyne.Window) fyne.CanvasObject {
	rscOnce.Do(func() {

	})
	return remoteSyncComponent
}
func remoteAuthDialogFunc(b bool) {
	if b {
		common.SendSignal2GWChannel(common.SIGNAL_AUTH_PASS)
	} else {
		common.SendSignal2GWChannel(common.SIGNAL_AUTH_NO_PASS)
	}
	remoteAuthShowing = false
}
