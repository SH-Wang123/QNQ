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

var (
	testSpeedComponent fyne.CanvasObject
	tsOnce             sync.Once
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

func GetTestDiskSpeedComponent(_ fyne.Window) fyne.CanvasObject {

	tsOnce.Do(func() {
		fileSizeSelect := widget.NewSelect([]string{"128MB", "512MB", "1GB", "4GB"}, nil)
		fileSizeComp := getLabelSelect("Test size:    ", fileSizeSelect)
		bufferSizeSelect := widget.NewSelect([]string{"512B", "1KB", "4KB"}, nil)
		bufferSizeComp := getLabelSelect("Buffer size: ", bufferSizeSelect)

		top := container.NewVBox(
			fileSizeComp,
			bufferSizeComp,
		)
		charts := container.NewMax()
		startBtn := widget.NewButton("Start", func() {

		})
		result := widget.NewLabel("Click start button to get result!")
		bottom := container.NewHSplit(charts, container.NewGridWithRows(2, startBtn, result))
		testSpeedComponent = container.NewVSplit(top, bottom)
	})
	return testSpeedComponent
}

func loadDiskInfo() {
	diskInfosContainer := container.NewVBox()
	for _, disk := range worker.DiskPartitionsCache {
		totalSize := binding.BindString(&disk.TotalSize)
		totalSizeLab := loadValue2Label("Total Size: ", totalSize)
		freeSize := binding.BindString(&disk.FreeSize)
		freeSizeLab := loadValue2Label("Free Size: ", freeSize)
		fsType := binding.BindString(&disk.FsType)
		fsTypeLabe := loadValue2Label("File System: ", fsType)

		usedPer := binding.BindFloat(&disk.UsedPercent)
		perProgress := widget.NewProgressBarWithData(usedPer)
		diskComp := container.NewVBox(
			totalSizeLab,
			freeSizeLab,
			fsTypeLabe,
			perProgress,
		)
		diskInfosContainer.Add(widget.NewCard(disk.Name, "", diskComp))
	}
	diskInfosContainer.Resize(classicSize)
	diskInfoComponent = container.NewMax(container.NewScroll(diskInfosContainer))
}
