package policies

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/davezant/decafein/src/server/processes"
)

func UpdateRunningBinaryByGui(a *processes.Activity) {
	gui := app.New()
	win := gui.NewWindow("Update Binary")

	// pegamos snapshot atual do bucket
	processList := processes.GlobalSnapshot.Processes

	listWidget := widget.NewList(
		func() int {
			return len(processList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(processList[i])
		},
	)

	selected := ""

	listWidget.OnSelected = func(id widget.ListItemID) {
		selected = processList[id]
	}

	saveBtn := widget.NewButton("Update Binary", func() {
		if selected != "" {
			a.ExecutionBinary = selected
			win.Close()
		}
	})
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		listWidget.Refresh()
	})

	bottomContainer := container.NewHBox(
		saveBtn,
		refreshBtn,
	)

	win.SetContent(
		container.NewBorder(nil, bottomContainer, nil, nil, listWidget),
	)

	win.Resize(fyne.NewSize(300, 400))
	win.ShowAndRun()
}
