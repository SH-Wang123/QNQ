package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sync"
	"window_handler/config"
	"window_handler/network"
	"window_handler/worker"
)

var (
	remoteSyncComponent fyne.CanvasObject
	rscOnce             sync.Once
)

func GetRemoteSyncComponent(_ fyne.Window) fyne.CanvasObject {
	rscOnce.Do(func() {
		bindingIp := binding.NewString()
		bindingIp.Set(config.SystemConfigCache.Value().QnqTarget.Ip)
		ipAddress := widget.NewEntryWithData(bindingIp)

		status := binding.NewBool()
		status.Set(network.ConnectStauts)
		connectedStatusComponent := container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Connected Status : "),
			widget.NewLabelWithData(binding.BoolToString(status)),
		)

		errLabel := widget.NewLabel("Connect failed !!!")
		okLabel := widget.NewLabel("Connect OK !!!")
		errLabel.TextStyle = fyne.TextStyle{}
		remoteSyncComponent = container.NewVBox(
			connectedStatusComponent,
			ipAddress,
			errLabel,
			okLabel,
			widget.NewButton("Test connect", func() {
				ret := network.TestPing(ipAddress.Text)
				if !ret {
					errLabel.Show()
					okLabel.Hide()
				} else {
					errLabel.Hide()
					okLabel.Show()
				}
			}),
			widget.NewButton("Save", func() {
				ret := network.TestPing(ipAddress.Text)
				if !ret {
					errLabel.Show()
					okLabel.Hide()
				} else {
					errLabel.Hide()
					okLabel.Show()
					config.SystemConfigCache.Cache.QnqTarget.Ip = ipAddress.Text
					config.SystemConfigCache.NotifyAll()
				}
			}),
			getConnectQTargetButton(),
		)
		errLabel.Hide()
		okLabel.Hide()
	})
	return remoteSyncComponent
}

func GetRemoteSingleComponent(win fyne.Window) fyne.CanvasObject {

	localPathBind := getBindString(config.SystemConfigCache.Value().QnqTarget.LocalPath)
	localFilePath := widget.NewLabelWithData(localPathBind)

	startButton := widget.NewButton("Start Sync", func() {
		qSender := worker.NewRemoteSyncSender()
		qSender.PrivateVariableMap["file_path"] = config.SystemConfigCache.Value().QnqTarget.LocalPath
		go qSender.ExecuteFunc(qSender)
	})
	filePathComponent := container.New(layout.NewHBoxLayout(),
		localFilePath,
		makeOpenFolderBtn("Open",
			win,
			localPathBind,
			&config.SystemConfigCache.Cache.QnqTarget.LocalPath),
	)

	return container.NewVBox(
		filePathComponent,
		startButton,
	)
}

func checkIpPing(ip string, errLabel *widget.Label) {
	ret := network.TestPing(ip)
	if !ret {
		errLabel.Show()
	} else {
		errLabel.Hide()
	}
}

func getConnectQTargetButton() *widget.Button {
	return widget.NewButton("Connect Target", func() {
		network.StartQClient()
	})
}
