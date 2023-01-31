package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sync"
	"window_handler/config"
)

var (
	localSystemInfoComponent fyne.CanvasObject
	lsiOnce                  sync.Once
)

func GetLocalSystemInfoComponent(_ fyne.Window) fyne.CanvasObject {
	lsiOnce.Do(func() {
		osInfoContainer := loadValue2Label("OS : ", getBindString(config.LocalSystemInfo.OS))
		sysFrameworkInfoContainer := loadValue2Label("System Framework : ", getBindString(config.LocalSystemInfo.SystemFramework))
		machineNameInfoContainer := loadValue2Label("Machine Name : ", getBindString(config.LocalSystemInfo.MachineName))
		localSystemInfoComponent = container.NewVBox(
			container.NewVBox(
				osInfoContainer,
				sysFrameworkInfoContainer,
				machineNameInfoContainer,
			),
		)
	})
	return localSystemInfoComponent
}

func GetDiskInfoComponent(_ fyne.Window) fyne.CanvasObject {
	return container.NewVBox(widget.NewLabel("Coming soon"))
}
