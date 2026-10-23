package components

import (
	"github.com/NathanFirmo/gorest/internal/curl"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// URLField is an input field that turns pasted curl commands into requests,
// like Insomnia does.
type URLField struct {
	*tview.InputField
	onCurl func(curl.Request, error)
}

func NewURLField() *URLField {
	u := &URLField{InputField: tview.NewInputField()}
	u.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			u.importCurl(u.GetText())
		}
	})
	return u
}

func (u *URLField) SetCurlHandler(handler func(curl.Request, error)) {
	u.onCurl = handler
}

func (u *URLField) PasteHandler() func(pastedText string, setFocus func(p tview.Primitive)) {
	paste := u.InputField.PasteHandler()
	return func(pastedText string, setFocus func(p tview.Primitive)) {
		if !u.importCurl(pastedText) {
			paste(pastedText, setFocus)
		}
	}
}

func (u *URLField) importCurl(text string) bool {
	if u.onCurl == nil || !curl.IsCurl(text) {
		return false
	}
	u.onCurl(curl.Parse(text))
	return true
}
