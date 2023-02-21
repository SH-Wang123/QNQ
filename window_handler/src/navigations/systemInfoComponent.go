package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"sync"
	"window_handler/config"
	"window_handler/worker"
)

var (
	localSystemInfoComponent fyne.CanvasObject
	lsiOnce                  sync.Once
)

var (
	diskInfoComponent fyne.CanvasObject
	dicOnce           sync.Once
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
	dicOnce.Do(func() {
		loadDiskInfo()
	})
	return diskInfoComponent
}

func loadDiskInfo() {
	diskInfosContainer := container.NewVBox()
	for _, disk := range worker.DiskPartitionsCache {
		name := widget.NewLabel("Disk Name: " + disk.Name)
		totalSize := binding.BindString(&disk.TotalSize)
		totalSizeLab := loadValue2Label("Total Size: ", totalSize)
		freeSize := binding.BindString(&disk.FreeSize)
		freeSizeLab := loadValue2Label("Free Size: ", freeSize)
		fsType := binding.BindString(&disk.FsType)
		fsTypeLabe := loadValue2Label("File System: ", fsType)

		usedPer := binding.BindFloat(&disk.UsedPercent)
		perProgress := widget.NewProgressBarWithData(usedPer)
		diskComp := container.NewVBox(
			name,
			totalSizeLab,
			freeSizeLab,
			fsTypeLabe,
			perProgress,
		)
		diskInfosContainer.Add(diskComp)
	}
	diskInfosContainer.Resize(classicSize)
	diskInfoComponent = container.NewMax(container.NewScroll(diskInfosContainer))

}
