package components

import (
	"github.com/rivo/tview"
)

type ResponseComponent struct {
	Container *tview.Flex
	Component *tview.TextView
}

// Create a flex container with row alignment to put response data
func CreateResponse() *ResponseComponent {
	cp := tview.NewTextView()
  cp.SetText("")

	ct := tview.NewFlex()
	ct.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("Response")

	ct.AddItem(cp, 0, 1, false)

	return &ResponseComponent{
		Container: ct,
		Component: cp,
	}
}
