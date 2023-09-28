package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RequestComponent struct {
	Container         *tview.Flex
	OptionsComponent  *tview.TextArea
	UrlInputComponent *tview.InputField
}

// Create a flex container with row alignment to put request option components
func CreateRequest() *RequestComponent {
	ct := tview.NewFlex()
	ct.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("Request")

	reqOpt := tview.NewTextArea()
	reqOpt.SetBorder(true)
	reqOpt.SetText("", true)

	reqUrlInput := tview.NewInputField()
	reqUrlInput.SetBorder(true)
	reqUrlInput.SetFieldTextColor(tcell.Color(tcell.Color.Hex(0)))

	// Add requests component on its container
	ct.
		AddItem(reqUrlInput, 0, 1, false).
		AddItem(reqOpt, 0, 9, false)

	return &RequestComponent{
		Container:         ct,
		OptionsComponent:  reqOpt,
		UrlInputComponent: reqUrlInput,
	}
}
