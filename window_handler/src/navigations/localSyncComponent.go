package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sync"
	"time"
	"window_handler/config"
	"window_handler/worker"
)

var FileSyncComponent fyne.CanvasObject
var ProgressBar *fyne.Container
var progressBox *fyne.Container

var (
	singleLocalSyncComponent fyne.CanvasObject
	slcOnce                  sync.Once
)

var (
	batchLocalSyncComponent fyne.CanvasObject
	blcOnce                 sync.Once
)

func GetSingleLocalSyncComponent(win fyne.Window) fyne.CanvasObject {
	slcOnce.Do(func() {
		sourcePathBind := getBindString(config.SystemConfigCache.Value().LocalSingleSync.SourcePath)
		sourceContainer := loadValue2Label("Source: ", sourcePathBind)

		sourcePathComponent := container.New(layout.NewHBoxLayout(),
			sourceContainer,
			makeOpenFileBtn("Open", win, sourcePathBind, &config.SystemConfigCache.Cache.LocalSingleSync.SourcePath),
		)

		targetPathBind := getBindString(config.SystemConfigCache.Value().LocalSingleSync.TargetPath)
		targetContainer := loadValue2Label("Target: ", targetPathBind)

		targetFileBtn := makeOpenFileBtn("Open", win, targetPathBind, &config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
		targetFolderBtn := makeOpenFolderBtn("Open", win, targetPathBind, &config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
		btnComponent := container.NewMax()
		btnComponent.Add(targetFolderBtn)

		createFileComponent := widget.NewCheck("Create New File", func(b bool) {
			if b {
				btnComponent.Remove(targetFileBtn)
				btnComponent.Add(targetFolderBtn)
			} else {
				btnComponent.Remove(targetFolderBtn)
				btnComponent.Add(targetFileBtn)
			}
		})

		targetPathComponent := container.New(layout.NewHBoxLayout(),
			targetContainer,
			btnComponent,
			createFileComponent,
		)
		singleLocalSyncComponent = container.NewVBox(
			sourcePathComponent,
			targetPathComponent,
			getStartLocalSingleButton(),
		)
	})
	return singleLocalSyncComponent
}

func GetBatchLocalSyncComponent(win fyne.Window) fyne.CanvasObject {
	blcOnce.Do(func() {
		sourcePathBind := getBindString(config.SystemConfigCache.Value().LocalBatchSync.SourcePath)
		sourceContainer := loadValue2Label("Source: ", sourcePathBind)
		sourceComponent := container.New(layout.NewHBoxLayout(), sourceContainer, makeOpenFolderBtn("Open",
			win,
			sourcePathBind,
			&config.SystemConfigCache.Cache.LocalBatchSync.SourcePath))
		targetPathBind := getBindString(config.SystemConfigCache.Value().LocalBatchSync.TargetPath)
		targetContainer := loadValue2Label("Target: ", targetPathBind)
		targetComponent := container.New(layout.NewHBoxLayout(), targetContainer, makeOpenFolderBtn("Open",
			win,
			targetPathBind,
			&config.SystemConfigCache.Cache.LocalBatchSync.TargetPath))

		startButton := getStartLocalBatchButton()

		differecnButton := getDiffAnalysisButton()
		progressBox = container.NewVBox()
		FileSyncComponent = container.NewVBox(
			container.NewVBox(
				sourceComponent,
				targetComponent,
			),
			startButton,
			differecnButton,
			progressBox,
		)
		batchLocalSyncComponent = container.NewBorder(FileSyncComponent, nil, getFileTree(), getFileTree())
	})
	//ProgressBar.Hide()
	return batchLocalSyncComponent
}

func getStartLocalBatchButton() *widget.Button {
	button := widget.NewButton("Start", func() {
		ProgressBar = getTaskProgressBar()
		progressBox.Add(ProgressBar)
		//TODO 优化协程池
		go worker.SyncBatchFileTree(worker.LocalBSFileNode, config.SystemConfigCache.Value().LocalBatchSync.TargetPath)

	})
	return button
}

func getStartLocalSingleButton() *widget.Button {
	button := widget.NewButton("Start", func() {
		go worker.LocalSyncSingleFileGUI()
	})
	return button
}

func getDiffAnalysisButton() *widget.Button {
	button := widget.NewButton("Variance Analysis", func() {
	})
	return button
}

func getTaskProgressBar() *fyne.Container {
	progress := widget.NewProgressBar()
	go func() {
		var progressNum = 0.0
		for progressNum < 1 {
			time.Sleep(time.Millisecond * 500)
			progressNum = worker.GetLocalBatchProgress()
			progress.SetValue(progressNum)
		}
		time.Sleep(time.Second * 10)
		progressBox.Remove(ProgressBar)

	}()
	return container.NewVBox(progress)
}

func getFileTree() fyne.CanvasObject {
	dataM := make(map[string][]string)
	worker.FileNode2TreeMap(&dataM)
	tree := widget.NewTreeWithStrings(dataM)
	size := fyne.Size{
		Height: 600,
	}
	tree.Resize(size)
	return tree
}
