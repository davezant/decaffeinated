package policies

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func CreateBlockWall() {
	var app = app.New()
	var w = app.NewWindow("Access Blocked")
	w.SetFullScreen(true)
	w.SetContent(canvas.NewText("Blocked access, ask your manager.", nil))
	w.ShowAndRun()
}
