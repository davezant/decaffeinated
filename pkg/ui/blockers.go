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

// ScreenBlockerConfig configures the full-screen blocker.
type ScreenBlockerConfig struct {
	Title       string
	Message     string
	Hint        string
	ButtonText  string
	BgColor     color.NRGBA
	TextColor   color.NRGBA
	OnUnlock    func()
}

// OpenScreenBlocker renders a customizable full-screen blocking overlay.
func OpenScreenBlocker(cfg ScreenBlockerConfig) {
	if cfg.Title == "" {
		cfg.Title = "Screen Blocked"
	}
	if cfg.Message == "" {
		cfg.Message = "Excesso de tempo de uso detectado."
	}
	if cfg.Hint == "" {
		cfg.Hint = "Use o login do Qiosk para desbloquear."
	}
	if cfg.ButtonText == "" {
		cfg.ButtonText = "Abrir Login Kiosk"
	}
	if cfg.BgColor == (color.NRGBA{}) {
		cfg.BgColor = color.NRGBA{R: 0x00, G: 0x24, B: 0x52, A: 0xEE}
	}
	if cfg.TextColor == (color.NRGBA{}) {
		cfg.TextColor = color.White
	}

	myApp := app.New()
	myWindow := myApp.NewWindow(cfg.Title)
	myWindow.SetFullScreen(true)
	myWindow.SetPadded(false)

	background := canvas.NewRectangle(cfg.BgColor)

	label := widget.NewLabelWithStyle(cfg.Message, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	label.TextColor = cfg.TextColor

	hint := widget.NewLabelWithStyle(cfg.Hint, fyne.TextAlignCenter, fyne.TextStyle{})
	hint.TextColor = cfg.TextColor

	unlockBtn := widget.NewButton(cfg.ButtonText, func() {
		if cfg.OnUnlock != nil {
			cfg.OnUnlock()
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

