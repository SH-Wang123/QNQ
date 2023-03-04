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

func GetRemoteSingleComponent(win fyne.Window) fyne.CanvasObject {
	rscOnce.Do(func() {
		bindingIp := binding.NewString()
		bindingIp.Set(config.SystemConfigCache.Value().QnqSTarget.Ip)
		ipAddress := widget.NewEntryWithData(bindingIp)
		ipAdressComp := container.New(
			layout.NewFormLayout(),
			widget.NewLabel("IP:"),
			ipAddress,
		)

		status := binding.NewBool()
		status.Set(network.ConnectStauts)
		connectedStatusComponent := container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Connected Status : "),
			widget.NewLabelWithData(binding.BoolToString(status)),
		)

		errLabel := widget.NewLabel("Connect failed !!!")
		okLabel := widget.NewLabel("Connect OK !!!")
		errLabel.TextStyle = fyne.TextStyle{}

		localPathBind := getBindString(config.SystemConfigCache.Value().QnqSTarget.LocalPath)
		localFilePath := loadValue2Label("Local path:", localPathBind)
		localPathComp := container.NewHBox(
			localFilePath,
			makeOpenFileBtn("Open",
				win,
				localPathBind,
				&config.SystemConfigCache.Cache.QnqSTarget.LocalPath),
		)
		remotePathBind := getBindString(config.SystemConfigCache.Value().QnqSTarget.RemotePath)
		remoteFilePathInput := widget.NewEntryWithData(remotePathBind)
		remotePathComp := container.New(
			layout.NewFormLayout(),
			widget.NewLabel("Remote path:"),
			remoteFilePathInput,
		)

		saveButton := widget.NewButton("Save IP & Remote path", func() {
			ret := checkIpPing(ipAddress.Text, errLabel, okLabel)
			if ret {
				config.SystemConfigCache.Cache.QnqSTarget.Ip = ipAddress.Text
			}
			config.SystemConfigCache.Cache.QnqSTarget.RemotePath = remoteFilePathInput.Text
			config.SystemConfigCache.NotifyAll()
		})
		startButton := widget.NewButton("Start", func() {
			qSender := worker.NewRemoteSyncSender()
			qSender.PrivateVariableMap["local_file_path"] = config.SystemConfigCache.Value().QnqSTarget.LocalPath
			go qSender.ExecuteFunc(qSender)
		})
		connectButton := getConnectQTargetButton()
		remoteSingleSyncPolicyComponent := getSingleSyncPolicyBtn(win, true)
		testButton := &widget.Button{
			Text:       "Test connect",
			Importance: widget.WarningImportance,
			OnTapped: func() {
				ret := checkIpPing(ipAddress.Text, errLabel, okLabel)
				if !ret {
					batchDisable(saveButton, startButton, connectButton, remoteSingleSyncPolicyComponent)
				} else {
					batchEnable(saveButton, startButton, connectButton, remoteSingleSyncPolicyComponent)
				}
			},
		}
		remoteSyncComponent = container.NewVBox(
			localPathComp,
			remotePathComp,
			ipAdressComp,
			connectedStatusComponent,
			errLabel,
			okLabel,
			testButton,
			saveButton,
			connectButton,
			startButton,
			remoteSingleSyncPolicyComponent,
		)

		errLabel.Hide()
		okLabel.Hide()
		batchDisable(saveButton, startButton, connectButton, remoteSingleSyncPolicyComponent)
	})
	return remoteSyncComponent
}

func checkIpPing(ip string, errLabel *widget.Label, okLabel *widget.Label) bool {
	ret := network.TestPing(ip)
	if !ret {
		okLabel.Hide()
		errLabel.Show()
	} else {
		errLabel.Hide()
		okLabel.Show()
	}
	return ret
}

func getConnectQTargetButton() *widget.Button {
	return widget.NewButton("Connect Target", func() {
		network.StartQClient()
	})
}
