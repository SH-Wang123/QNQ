package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"unsafe"
	"window_handler/config"
)

func bindingValue2Label(text string, value interface{}) *fyne.Container {
	var valueLabel *widget.Label
	unsafePtr := unsafe.Pointer(&value)
	switch value.(type) {
	case string:
		data := binding.NewString()
		data.Set(*(*string)(unsafePtr))
		valueLabel = widget.NewLabelWithData(data)
	case int:
		data := binding.NewInt()
		data.Set(*(*int)(unsafePtr))
		valueLabel = widget.NewLabelWithData(binding.IntToString(data))
	case float64:
		data := binding.NewFloat()
		data.Set(*(*float64)(unsafePtr))
		valueLabel = widget.NewLabelWithData(binding.FloatToString(data))
	case bool:
		data := binding.NewBool()
		data.Set(*(*bool)(unsafePtr))
		valueLabel = widget.NewLabelWithData(binding.BoolToString(data))
	}
	return container.New(layout.NewHBoxLayout(),
		widget.NewLabel(text),
		valueLabel,
	)
}

func loadValue2Label(text string, bindValue binding.String) *fyne.Container {
	return container.New(layout.NewHBoxLayout(),
		widget.NewLabel(text),
		widget.NewLabelWithData(bindValue),
	)
}

func makeOpenFolderBtn(buttonName string, win fyne.Window, bindPath binding.String, configStr *string) *widget.Button {
	return widget.NewButton(buttonName, func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			_, err = uri.List()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			bindPath.Set(uri.Path())
			*configStr = uri.Path()
			config.SystemConfigCache.NotifyAll()
		}, win)
	})
}

func makeBtn(buttonName string, clickFunc func()) *widget.Button {
	return widget.NewButton(buttonName, clickFunc)
}

func makeOpenFileBtn(buttonName string, win fyne.Window, bindPath binding.String, configStr *string) *widget.Button {
	return widget.NewButton(buttonName, func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri == nil {
				log.Println("Cancelled")
				return
			}
			bindPath.Set(uri.URI().Path())
			*configStr = uri.URI().Path()
			config.SystemConfigCache.NotifyAll()
		}, win)
	})
}

func getBindString(value string) binding.String {
	bindD := binding.NewString()
	bindD.Set(value)
	return bindD
}

func batchDisable(ds ...fyne.Disableable) {
	for _, d := range ds {
		d.Disable()
	}
}

func batchEnable(ds ...fyne.Disableable) {
	for _, d := range ds {
		d.Enable()
	}
}

func batchRefresh(bs ...fyne.CanvasObject) {
	for _, v := range bs {
		v.Refresh()
	}
}

func clearDisableRootCache(key string) {
	disableRootCache[key] = make(map[fyne.Disableable]disableRoot)
}

func addDisableRoot(key string, root fyne.Disableable, chidls ...fyne.Disableable) {
	currentRoot := disableRoot{}
	currentRoot.addChild(chidls...)
	disableRootCache[key][root] = currentRoot
}

func disableAllChild(key string, root fyne.Disableable) {
	cache := disableRootCache[key]
	currentRoot := cache[root]
	currentRoot.disableChild()
	for _, v := range currentRoot.child {
		disableAllChild(key, v)
	}
}

func enableAllChild(key string, root fyne.Disableable) {
	cache := disableRootCache[key]
	currentRoot := cache[root]
	currentRoot.enableChild()
}

// swapChecked activation check function
func swapChecked(w *widget.Check) {
	w.SetChecked(!w.Checked)
	w.SetChecked(!w.Checked)
}
