package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
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
	remoteSyncComponent       *fyne.Container
	remoteQNQSelect           *widget.Select
	remoteSyncPolicyComponent *widget.Button
	remoteSyncStartButton     *widget.Button
	remoteSyncProgressBox     *fyne.Container
	remoteSyncCurrentFile     *widget.Label
	remoteSyncTimeRemaining   *widget.Label
	rsOnce                    sync.Once
)

var (
	serverListTable *widget.Table
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
	tabs.CloseIntercept = func(item *container.TabItem) {
		index, _ := strconv.Atoi(item.Text)
		if len(config.SystemConfigCache.Cache.QNQNetCells) < index {
			return
		}
		index--
		ip := config.SystemConfigCache.Cache.QNQNetCells[index].Ip
		go config.SystemConfigCache.DeleteQNQNetCell(index)
		go common.DisconnectTarget(ip)
		tabs.Refresh()
	}
	tabs.SetTabLocation(container.TabLocationTrailing)

	return container.NewAppTabs(
		container.NewTabItem("Target List", container.NewBorder(nil, nil, nil, nil, tabs)),
		container.NewTabItem("Server List", getServerListTable()),
	)
}

func createRemoteQNQItem(index int) *container.TabItem {
	ipComp, ipInput := getLabelInput("IP : ", config.SystemConfigCache.Cache.QNQNetCells[index].Ip)
	netCell := common.GetQNetCell(config.SystemConfigCache.Cache.QNQNetCells[index].Ip)
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
		go common.ConnectTarget(ipInput.Text)
		waitAuthDialog.Show()
	})
	disConnectBtn := widget.NewButton("Disconnect", func() {
		go common.DisconnectTarget(ipInput.Text)
	})
	connectBtn.Importance = widget.HighImportance
	disConnectBtn.Importance = widget.WarningImportance
	comp := container.NewVBox(
		widget.NewLabel(""),
		ipComp,
		targetStatusComp,
		markComp,
		connectBtn,
		disConnectBtn,
		saveBtn,
	)
	if netCell.GetTargetStatus() {
		batchDisable(connectBtn)
		batchEnable(disConnectBtn)
	} else {
		batchEnable(connectBtn)
		batchDisable(disConnectBtn)
	}
	return container.NewTabItem(
		fmt.Sprint(index+1),
		comp,
	)
}

func getRemoteSingleComponent(win fyne.Window) fyne.CanvasObject {
	rsOnce.Do(func() {

		localPathStr := ""
		localPathBind := getBindString(localPathStr)
		localPath := loadValue2Label("Local Path : ", localPathBind)
		localPathComp := container.New(layout.NewHBoxLayout(), localPath, makeOpenFileBtn("Open",
			win,
			localPathBind,
			&localPathStr))
		remotePathInput := widget.NewEntry()
		remotePathComp := container.NewHBox(
			widget.NewLabel("Remote Path : "),
			remotePathInput,
		)
		remoteSyncProgressBox = container.NewVBox()
		remoteSyncPolicyComponent = getSingleSyncPolicyBtn(win, true)
		remoteSyncStartButton = widget.NewButton("Start", func() {
			if remoteQNQSelect.Selected == "" {
				dialog.NewInformation("Start remote sync error", "Please add remote QNQ.", win).Show()
				return
			}
			go worker.RemoteSingleSyncSingleTime(localPathStr, remotePathInput.Text, remoteQNQSelect.Selected)
		})
		remoteSyncCurrentFile = widget.NewLabel(NOT_RUNNING_STR)
		currentFileComp := container.NewHBox(
			widget.NewLabel("Current sync file:"),
			remoteSyncCurrentFile,
		)
		remoteSyncTimeRemaining = widget.NewLabel(NOT_RUNNING_STR)
		rsTimeComp := container.NewHBox(
			widget.NewLabel("Time remaining:"),
			remoteSyncTimeRemaining,
		)
		remoteSyncComponent = container.NewVBox(
			localPathComp,
			remotePathComp,
			currentFileComp,
			rsTimeComp,
			remoteSyncProgressBox,
			remoteSyncPolicyComponent,
			remoteSyncStartButton,
		)
		remotePathComp.Hide()
	})
	targets := common.GetAllQSorT(false)
	targetIps := make([]string, 0)
	for _, v := range targets {
		target := *v
		ip := common.GetIpFromAddr(target.RemoteAddr().String())
		targetIps = append(targetIps, ip)
	}
	remoteQNQSelect = widget.NewSelect(targetIps, func(s string) {

	})
	qnqSelectComp := container.NewHBox(
		widget.NewLabel("Select Remote QNQ : "),
		remoteQNQSelect,
	)
	return container.NewVBox(qnqSelectComp, remoteSyncComponent)
}
func remoteAuthDialogFunc(b bool) {
	if b {
		common.SendSignal2GWChannel(common.SIGNAL_AUTH_PASS)
	} else {
		common.SendSignal2GWChannel(common.SIGNAL_AUTH_NO_PASS)
	}
	remoteAuthShowing = false
}

func getServerListTable() *widget.Table {
	servers := common.GetAllQSorT(true)
	tableData := make([][]string, 0)
	tableData = append(tableData, []string{"IP", "Status"})
	for _, server := range servers {
		s := *server
		data := []string{s.RemoteAddr().String(), "Connected"}
		tableData = append(tableData, data)
	}
	rowNum := len(tableData)
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
			switch id.Col {
			case 0:
				if ids == 0 {
					label.SetText(fmt.Sprint())
				} else {
					label.SetText(fmt.Sprintf("%d", ids))
				}
			case 1:
				label.SetText(tableData[ids][0])
			case 2:
				label.SetText(tableData[ids][1])
			}
		},
	)
	t.SetColumnWidth(0, 40)
	t.SetColumnWidth(1, 160)
	t.SetColumnWidth(2, 80)
	serverListTable = t
	serverListTable.Refresh()
	return serverListTable
}
