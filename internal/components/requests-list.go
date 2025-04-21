package components

import (
	"github.com/rivo/tview"
)

type Request struct {
	ID      int64
	URL     string
	Name    string
	Method  string
	Headers string
	Body    string
}

type RequestsListComponent struct {
	Container *tview.Flex
	Component *tview.List
	Items     []Request
}

// AddItem a new request into requests list
func (rl *RequestsListComponent) AddItem(r *Request, cb func()) {
	rl.Component.AddItem(r.Name, r.URL, 0, cb)
	rl.Items = append(rl.Items, *r)
}

// Create a flex container with row alignment to put requests list
func CreateRequestsList() *RequestsListComponent {
	cp := tview.NewList()

	ct := tview.NewFlex()
	ct.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("List")

	ct.AddItem(cp, 0, 1, false)

	return &RequestsListComponent{
		Container: ct,
		Component: cp,
		Items:     make([]Request, 0),
	}
}

func (rl *RequestsListComponent) UpdateItem(r *Request) {
	if r.ID >= 0 && r.ID < int64(len(rl.Items)) {
		rl.Items[r.ID].Name = r.Name
		rl.Items[r.ID].URL = r.URL
		rl.Items[r.ID].Headers = r.Headers
		rl.Items[r.ID].Body = r.Body
	}
}
