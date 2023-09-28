package components

import "github.com/rivo/tview"

// Create a component for request options input
func CreateRequestOptionsInput() *tview.TextArea {
	c := tview.NewTextArea()
  c.SetBorder(true)
	c.SetText("", true)

	return c
}
