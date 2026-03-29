package ui

import (
"fyne.io/fyne/v2"
"fyne.io/fyne/v2/app"
"fyne.io/fyne/v2/container"
"fyne.io/fyne/v2/layout"
"fyne.io/fyne/v2/widget"
)

// OpenLoginKiosk starts a full-screen login form for kiosk mode.
func OpenLoginKiosk(onSuccess func(user string), onFailure func(err error)) {
myApp := app.New()
myWindow := myApp.NewWindow("Login Kiosk")
myWindow.SetFullScreen(true)

userEntry := widget.NewEntry()
userEntry.SetPlaceHolder("Usuário")
passEntry := widget.NewPasswordEntry()
passEntry.SetPlaceHolder("Senha")

status := widget.NewLabel("")

loginBtn := widget.NewButton("Entrar", func() {
user := userEntry.Text
pass := passEntry.Text
if user == "" || pass == "" {
status.SetText("Preencha usuário e senha")
return
}

// TODO: ajustar a validação real
if user == "admin" && pass == "password" {
status.SetText("Login bem-sucedido")
if onSuccess != nil {
onSuccess(user)
}
myWindow.Close()
return
}

status.SetText("Credenciais inválidas")
if onFailure != nil {
onFailure(nil)
}
})

form := container.NewVBox(
widget.NewLabelWithStyle("Bem-vindo ao modo kiosk", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
userEntry,
passEntry,
loginBtn,
status,
)

myWindow.SetContent(container.NewCenter(container.New(layout.NewGridLayout(1), form)))
myWindow.ShowAndRun()
}
