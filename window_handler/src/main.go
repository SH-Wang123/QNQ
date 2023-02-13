package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net"
	"net/url"
	"os"
	"window_handler/Qlog"
	"window_handler/common"
	"window_handler/config"
	"window_handler/navigations"
	"window_handler/network"
	"window_handler/request"
	"window_handler/worker"
)

const preferenceCurrentNavigation = "currentNavigation"

var topWindow fyne.Window

func main() {
	os.Setenv("FYNE_FONT", "msyh.ttc")
	Qlog.MakeLogger()
	common.InitCoroutinesPool()
	worker.LoadWorkerFactory()
	go network.StartQServer()
	go network.StartQClient()
	network.NetChan.StartPump()
	//config.GetTargetSystemInfo()
	worker.InitFileNode(false, false)
	common.GetCoroutinesPool().StartPool()
	startGUI()
	if network.ConnectStauts {
		defer func(ConnectClient net.Conn) {
			err := ConnectClient.Close()
			if err != nil {

			}
		}(network.ConnectClient)
	}
}

func startGUI() {
	a := app.NewWithID("qnq.window_handler")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("QNQ Sync " + config.SystemConfigCache.Value().Version)
	topWindow = w
	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setNavigation := func(t navigations.Navigation) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setNavigation, false))
	} else {
		split := container.NewHSplit(makeNav(setNavigation, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(400, 600))

	w.ShowAndRun()
}

// func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
//
// }
func makeNav(setNavigation func(navigation navigations.Navigation), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return navigations.NavigationIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := navigations.NavigationIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := navigations.Navigations[uid]
			if !ok {
				fyne.LogError("Missing navigations panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedNavigation(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := navigations.Navigations[uid]; ok {
				if unsupportedNavigation(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentNavigation, uid)
				setNavigation(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentNavigation, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)
	lastVersion := request.GetLastVersion()
	if lastVersion == "" || lastVersion == config.SystemConfigCache.Value().Version {
		lastVersion = "Project Link"
	} else {
		lastVersion = "Lastest version is  " + lastVersion
	}
	projectUrl, _ := url.Parse("https://github.com/wangshenghao1/QNQ")
	versionInfo := widget.NewHyperlink(lastVersion, projectUrl)
	flootBox := container.NewVBox(
		versionInfo,
		themes,
	)

	return container.NewBorder(nil, flootBox, nil, nil, tree)
}

func unsupportedNavigation(t navigations.Navigation) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}
