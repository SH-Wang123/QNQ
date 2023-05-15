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

	return container.NewAppTabs(
		container.NewTabItem("Target List", container.NewBorder(nil, nil, nil, nil, tabs)),
		container.NewTabItem("Server List", getServerListTable(win)),
	)
}

func createRemoteQNQItem(index int) *container.TabItem {
	ipComp, ipInput := getLabelInput("IP : ", config.SystemConfigCache.Cache.QNQNetCells[index].Ip)
	netCell := worker.GetQNetCell(config.SystemConfigCache.Cache.QNQNetCells[index].Ip)
	targetStatusComp := container.NewHBox(
		widget.NewLabel("Connect status : "),
		widget.NewLabel(fmt.Sprint(netCell.GetTargetStatus())),
	)
	markComp, markInput := getLabelInput("Mark : ", config.SystemConfigCache.Cache.QNQNetCells[index].Mark)
	saveBtn := widget.NewButton("Save", func() {
		config.SystemConfigCache.Cache.QNQNetCells[index].Ip = ipInput.Text
		config.SystemConfigCache.Cache.QNQNetCells[index].Mark = markInput.Text
		config.SystemConfigCache.NotifyAll()
	})
	connectBtn := widget.NewButton("Connect", func() {
		go worker.ConnectTarget(ipInput.Text)
		waitAuthDialog.Show()
		//go func() {
		//	t := time.NewTicker(60 * time.Second)
		//	select {
		//	case <-t.C:
		//		waitAuthDialog.Hide()
		//		authErrorDialog.Show()
		//	}
		//}()
	})
	connectBtn.Importance = widget.HighImportance
	comp := container.NewVBox(
		widget.NewLabel(""),
		ipComp,
		targetStatusComp,
		markComp,
		connectBtn,
		saveBtn,
	)
	if netCell.GetTargetStatus() {
		batchDisable(connectBtn)
	} else {
		batchEnable(connectBtn)
	}
	return container.NewTabItem(
		fmt.Sprint(index+1),
		comp,
	)
}

func getRemoteSingleComponent(win fyne.Window) fyne.CanvasObject {
	rscOnce.Do(func() {
		remoteSyncComponent = container.NewHBox(widget.NewLabel("wait reconsitution"))
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

func getServerListTable(_ fyne.Window) *widget.Table {
	servers := worker.GetAllQServers()
	rowNum := len(servers)
	t := widget.NewTable(
		func() (int, int) {
			return rowNum, 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Nothing")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			ids := id.Row
			server := servers[ids]
			s := *server
			switch id.Col {
			case 0:
				if ids == 0 {
					label.SetText("")
				} else {
					label.SetText(fmt.Sprintf("%d", ids))
				}
			case 1:
				label.SetText(fmt.Sprint(s.RemoteAddr()))
			case 2:
				label.SetText(fmt.Sprint())
			case 3:
				label.SetText(fmt.Sprint())
			}
		},
	)
	t.SetColumnWidth(0, 130)
	t.SetColumnWidth(1, 60)
	t.SetColumnWidth(2, 60)
	return t
}
