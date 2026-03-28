package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

func OpenScreenBlocker() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Locker")
	background := canvas.NewRectangle(color.White)
	myWindow.SetContent(container.NewStack(background))
	myWindow.SetFullScreen(true)
	myWindow.ShowAndRun()
}

