package components

import (
	"github.com/rivo/tview"
)

type Request struct {
	Url     string
	Name    string
	Method  string
	Headers string
	Body    string
}

type RequestsListComponent struct {
	Items     []Request
	Component *tview.List
}

// AddItem a new request into requests list
func (rl *RequestsListComponent) AddItem(r *Request, cb func()) {
  rl.Component.AddItem(r.Name, r.Url, 0, cb)

	rl.Items = append(rl.Items, *r)
}

// Create a component to show requests list
func CreateRequestsList() *RequestsListComponent {
	c := tview.NewList()

	return &RequestsListComponent{
		Component: c,
		Items:     make([]Request, 0),
	}
}
