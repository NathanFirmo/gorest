package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Create a component for url input
func CreateUrlInput() *tview.InputField {
	c := tview.NewInputField()
	c.SetBorder(true)
	c.SetFieldTextColor(tcell.Color(tcell.Color.Hex(0)))

	return c
}
