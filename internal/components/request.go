package components

import (
	"slices"

	"github.com/rivo/tview"
)

var Methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

type RequestComponent struct {
	Container        *tview.Flex
	NameComponent    *tview.InputField
	MethodComponent  *tview.DropDown
	UrlComponent     *URLField
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

	method := tview.NewDropDown()
	method.SetBorder(true)
	method.SetTitle("Method")
	method.SetOptions(Methods, nil)
	method.SetCurrentOption(0)

	url := NewURLField()
	url.SetBorder(true)
	url.SetTitle("URL")
	url.SetPlaceholder("https://... or paste a curl command")

	headers := tview.NewTextArea()
	headers.SetBorder(true)
	headers.SetTitle("Headers")
	headers.SetPlaceholder("Content-Type: application/json")

	body := tview.NewTextArea()
	body.SetBorder(true)
	body.SetTitle("Body")

	target := tview.NewFlex().
		AddItem(method, 11, 0, false).
		AddItem(url, 0, 1, false)

	ct.
		AddItem(name, 3, 0, false).
		AddItem(target, 3, 0, false).
		AddItem(headers, 0, 1, false).
		AddItem(body, 0, 1, false)

	return &RequestComponent{
		Container:        ct,
		NameComponent:    name,
		MethodComponent:  method,
		UrlComponent:     url,
		HeadersComponent: headers,
		BodyComponent:    body,
	}
}

// Fields returns the request inputs in navigation order.
func (rc *RequestComponent) Fields() []tview.Primitive {
	return []tview.Primitive{rc.NameComponent, rc.MethodComponent, rc.UrlComponent, rc.HeadersComponent, rc.BodyComponent}
}

func (rc *RequestComponent) Method() string {
	_, method := rc.MethodComponent.GetCurrentOption()
	return method
}

func (rc *RequestComponent) SetMethod(method string) {
	rc.MethodComponent.SetCurrentOption(max(slices.Index(Methods, method), 0))
}

func (rc *RequestComponent) Value(id int64) Request {
	return Request{
		ID:      id,
		Name:    rc.NameComponent.GetText(),
		Method:  rc.Method(),
		URL:     rc.UrlComponent.GetText(),
		Headers: rc.HeadersComponent.GetText(),
		Body:    rc.BodyComponent.GetText(),
	}
}

func (rc *RequestComponent) SetValue(r Request) {
	rc.NameComponent.SetText(r.Name)
	rc.SetMethod(r.Method)
	rc.UrlComponent.SetText(r.URL)
	rc.HeadersComponent.SetText(r.Headers, false)
	rc.BodyComponent.SetText(r.Body, false)
}
