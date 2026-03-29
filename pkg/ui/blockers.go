package ui

import (
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func OpenScreenBlocker() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Locker")
	background := canvas.NewRectangle(color.White)
	myWindow.SetContent(container.NewStack(background))
	myWindow.SetFullScreen(true)
	myWindow.ShowAndRun()
}

func OpenLoginScreen(){
	myApp := app.New()
	myWindow := myApp.NewWindow("Login")
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Login")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Digite sua senha:")
	loginBtn := widget.NewButton("Entrar" , func(){})
	loginForm := container.NewVBox(
		widget.NewLabel("Bem Vindo"),
		emailEntry,
		passwordEntry,
		loginBtn,
	)
	content := container.New(layout.NewGridLayout(2), loginForm)
	myWindow.SetContent(content)
	myWindow.SetFullScreen(true)
	myWindow.ShowAndRun()
}
