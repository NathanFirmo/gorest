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
	Format    func(Request) (main, secondary string)
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
		Format: func(r Request) (string, string) {
			return tview.Escape(r.Name), tview.Escape(r.URL)
		},
	}
}

func (rl *RequestsListComponent) AddItem(r Request) {
	rl.Items = append(rl.Items, r)
	main, secondary := rl.Format(r)
	rl.Component.AddItem(main, secondary, 0, nil)
}

func (rl *RequestsListComponent) UpdateItem(index int, r Request) {
	if index < 0 || index >= len(rl.Items) {
		return
	}
	rl.Items[index] = r
	main, secondary := rl.Format(r)
	rl.Component.SetItemText(index, main, secondary)
}

func (rl *RequestsListComponent) RemoveItem(index int) {
	if index < 0 || index >= len(rl.Items) {
		return
	}
	rl.Items = append(rl.Items[:index], rl.Items[index+1:]...)
	rl.Component.RemoveItem(index)
}

func (rl *RequestsListComponent) Current() (int, Request, bool) {
	index := rl.Component.GetCurrentItem()
	if index < 0 || index >= len(rl.Items) {
		return index, Request{}, false
	}
	return index, rl.Items[index], true
}

// Refresh re-renders every item, e.g. after the theme changes.
func (rl *RequestsListComponent) Refresh() {
	for i, r := range rl.Items {
		main, secondary := rl.Format(r)
		rl.Component.SetItemText(i, main, secondary)
	}
}
