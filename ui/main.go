package clientui

type WindowUI struct {
	Title  string
	Width  int
	Height int
}

func NewWindow() *WindowUI {
	return &WindowUI{
		Title:  "Decaffeinated Client",
		Width:  800,
		Height: 600,
	}
}
