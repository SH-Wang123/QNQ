package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strings"
	"window_handler/config"
)

func getLogComponent(_ fyne.Window) fyne.CanvasObject {
	log := config.LoadCSV(false)
	rowNum := len(log)

	t := widget.NewTable(
		func() (int, int) {
			return rowNum, 7
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
			case 6:
				label.SetText(clos[5])
			}
		},
	)
	t.SetColumnWidth(0, 30)
	t.SetColumnWidth(1, 130)
	t.SetColumnWidth(2, 140)
	t.SetColumnWidth(3, 140)
	t.SetColumnWidth(4, 80)
	t.SetColumnWidth(5, 350)
	t.SetColumnWidth(6, 350)
	return t
}
