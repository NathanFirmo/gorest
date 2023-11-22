package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RequestComponent struct {
	Container        *tview.Flex
	UrlComponent     *tview.InputField
	NameComponent    *tview.InputField
	HeadersComponent *tview.TextArea
	BodyComponent    *tview.TextArea
}

// Create a flex container with row alignment to put request option components
func CreateRequest() *RequestComponent {
	ct := tview.NewFlex()
	ct.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("Request")

	name := tview.NewInputField()
	name.SetBorder(true)
	name.SetTitle("Name")
	name.SetFieldTextColor(tcell.Color(tcell.Color.Hex(0)))

	url := tview.NewInputField()
	url.SetBorder(true)
	url.SetTitle("URL")
	url.SetFieldTextColor(tcell.Color(tcell.Color.Hex(0)))

	headers := tview.NewTextArea()
	headers.SetBorder(true)
	headers.SetTitle("Headers")
	headers.SetText("", true)

	body := tview.NewTextArea()
	body.SetBorder(true)
	body.SetTitle("Body")
	body.SetText("", true)

	// Add requests component on its container
	ct.
		AddItem(name, 0, 1, false).
		AddItem(url, 0, 1, false).
		AddItem(headers, 0, 4, false).
		AddItem(body, 0, 4, false)

	return &RequestComponent{
		Container:        ct,
		HeadersComponent: headers,
		BodyComponent:    body,
		UrlComponent:     url,
		NameComponent:    name,
	}
}
