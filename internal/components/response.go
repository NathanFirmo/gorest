package components

import "github.com/rivo/tview"

// Create a component to show response
func CreateResponse() *tview.TextView {
  c := tview.NewTextView()
  c.SetText("")

  return c
}
