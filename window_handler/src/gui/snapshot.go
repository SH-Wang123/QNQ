package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"
	"sync"
	"window_handler/config"
)

var (
	createTimePoint fyne.CanvasObject
	ctpOnce         sync.Once
)

func getTimePointComponent(win fyne.Window) fyne.CanvasObject {
	component := container.NewAppTabs(
		container.NewTabItem("Create Time Point", getCreateTimePoint(win)),
		container.NewTabItem("Restore", widget.NewLabel("Developing")),
		container.NewTabItem("Time Point Table", getTimePointTable()),
	)
	return component
}

func getCreateTimePoint(win fyne.Window) fyne.CanvasObject {
	ctpOnce.Do(func() {
		sourcePathBind := getBindString(config.SystemConfigCache.Value().LocalBatchSync.SourcePath)
		sourceContainer := loadValue2Label("Source Path: ", sourcePathBind)
		sourceComponent := container.New(layout.NewHBoxLayout(), sourceContainer, makeOpenFolderBtn("Open",
			win,
			sourcePathBind,
			&config.SystemConfigCache.Cache.LocalBatchSync.SourcePath))

		targetPathBind := getBindString(config.SystemConfigCache.Value().LocalBatchSync.TargetPath)
		targetContainer := loadValue2Label("Time Point Path: ", targetPathBind)
		targetComponent := container.New(layout.NewHBoxLayout(), targetContainer, makeOpenFolderBtn("Open",
			win,
			targetPathBind,
			&config.SystemConfigCache.Cache.LocalBatchSync.TargetPath))

	})
	return createTimePoint
}

func getTimePointTable() fyne.CanvasObject {
	log := config.LoadCSV(true)
	rowNum := len(log)

	t := widget.NewTable(
		func() (int, int) {
			return rowNum, 6
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Nothing")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			ids := id.Row
			clos := strings.Split(log[ids], ",")
			switch id.Col {
			case 0:
				if ids == 0 {
					label.SetText("")
				} else {
					label.SetText(fmt.Sprintf("%d", ids))
				}

			case 1:
				label.SetText(clos[0])
			case 2:
				label.SetText(clos[1])
			case 3:
				label.SetText(clos[2])
			case 4:
				label.SetText(clos[3])
			case 5:
				label.SetText(clos[4])
			}
		},
	)
	t.SetColumnWidth(0, 30)
	t.SetColumnWidth(1, 130)
	t.SetColumnWidth(2, 140)
	t.SetColumnWidth(3, 350)
	t.SetColumnWidth(4, 350)
	t.SetColumnWidth(5, 180)
	return t
}
