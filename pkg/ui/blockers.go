package ui

import (
"image/color"

"fyne.io/fyne/v2"
"fyne.io/fyne/v2/app"
"fyne.io/fyne/v2/canvas"
"fyne.io/fyne/v2/container"
"fyne.io/fyne/v2/layout"
"fyne.io/fyne/v2/widget"
)

// OpenScreenBlocker displays a full-screen blocking overlay with an unlock action.
func OpenScreenBlocker(onUnlock func()) {
myApp := app.New()
myWindow := myApp.NewWindow("Screen Blocker")
myWindow.SetFullScreen(true)
myWindow.SetPadded(false)

background := canvas.NewRectangle(color.NRGBA{R: 0x00, G: 0x24, B: 0x52, A: 0xEE})
label := widget.NewLabelWithStyle("Sistema Bloqueado", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
label.TextColor = color.White

hint := widget.NewLabelWithStyle("Para desbloquear use o modo kiosk", fyne.TextAlignCenter, fyne.TextStyle{})
hint.TextColor = color.White

unlockBtn := widget.NewButton("Iniciar Login (Kiosk)", func() {
if onUnlock != nil {
onUnlock()
}
myWindow.Close()
})

content := container.NewVBox(
layout.NewSpacer(),
label,
hint,
unlockBtn,
layout.NewSpacer(),
)

stack := container.NewMax(background, container.NewCenter(content))
myWindow.SetContent(stack)
myWindow.ShowAndRun()
}
