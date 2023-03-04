package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
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
	testSpeedRetLab    *widget.Label
	partitionSelect    *widget.Select
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
		bufferSizeSelect := widget.NewSelect([]string{"512Byte", "1KB", "4KB", "8KB", "1MB", "4MB"}, nil)
		bufferSizeComp := getLabelSelect("Buffer size: ", bufferSizeSelect)
		var partitions []string
		for _, v := range worker.DiskPartitionsCache {
			partitions = append(partitions, v.Name)
		}
		partitionSelect = widget.NewSelect(partitions, nil)
		partitionComp := getLabelSelect("Disk:          ", partitionSelect)
		errorText := widget.NewTextGridFromString("\nPlease select parameters!")
		errorText.SetRowStyle(1, &widget.CustomTextGridStyle{FGColor: &color.NRGBA{R: 255, G: 0, B: 0, A: 255}, BGColor: color.White})
		top := container.NewVBox(
			partitionComp,
			fileSizeComp,
			bufferSizeComp,
			errorText,
		)
		errorText.Hide()
		charts := container.NewMax()
		startBtn := widget.NewButton("Start", func() {
			if partitionSelect.Selected == "" || fileSizeSelect.Selected == "" || bufferSizeSelect.Selected == "" {
				errorText.Show()
				return
			} else {
				errorText.Hide()
			}
			totalPath := partitionSelect.Selected
			fileSize := worker.ConvertCapacity(fileSizeSelect.Selected)
			bufferSize := worker.ConvertCapacity(bufferSizeSelect.Selected)

			go worker.TestDiskSpeed(bufferSize, fileSize, totalPath)
		})
		testSpeedRetLab = widget.NewLabel("Click start button to get result!")
		bottom := container.NewHSplit(charts, container.NewGridWithRows(2, startBtn, testSpeedRetLab))
		testSpeedComponent = container.NewVSplit(top, bottom)
	})
	return testSpeedComponent
}

func loadDiskInfo() {
	diskInfosContainer := container.NewVBox()
	refreshButton := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		worker.GetPartitionsInfo()
	})
	diskInfosContainer.Add(refreshButton)
	for _, disk := range worker.DiskPartitionsCache {
		totalSize := getBindString(disk.TotalSizeStr)
		totalSizeLab := loadValue2Label("Total Size: ", totalSize)
		freeSize := getBindString(disk.FreeSizeStr)
		freeSizeLab := loadValue2Label("Free Size: ", freeSize)
		fsType := getBindString(disk.FsType)
		fsTypeLabe := loadValue2Label("File System: ", fsType)
		refreshButton.Resize(fyne.NewSize(15, 15))
		refreshButton.Refresh()
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
