package gorest

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/NathanFirmo/gorest/internal/utils"
	"github.com/gdamore/tcell/v2"
)

// Sets an functions that intercepts keyboard inputs and allows to check pressed key and handle focus
func (a *App) SetInputHandlers() {
	a.rootContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focus := a.tview.GetFocus()

		switch event.Key() {
		case tcell.KeyCtrlL:
			if focus == a.requestsList.Component || focus == a.requestsList.Container {
				a.tview.SetFocus(a.request.NameComponent)
			} else {
				a.tview.SetFocus(a.response.Container)
			}
			return event
		case tcell.KeyCtrlH:
			if focus == a.response.Container {
				a.tview.SetFocus(a.request.NameComponent)
			} else {
				a.tview.SetFocus(a.requestsList.Component)
			}
			return event
		case tcell.KeyCtrlJ:
			if focus == a.request.NameComponent {
				a.tview.SetFocus(a.request.UrlComponent)
			}
			if focus == a.request.UrlComponent {
				a.tview.SetFocus(a.request.HeadersComponent)
			}
			if focus == a.request.HeadersComponent {
				a.tview.SetFocus(a.request.BodyComponent)
			}
			return event
		case tcell.KeyCtrlK:
			if focus == a.request.BodyComponent {
				a.tview.SetFocus(a.request.HeadersComponent)
			}
			if focus == a.request.HeadersComponent {
				a.tview.SetFocus(a.request.UrlComponent)
			}
			if focus == a.request.UrlComponent {
				a.tview.SetFocus(a.request.NameComponent)
			}
			return event
		case tcell.KeyEnter:
			url := a.request.UrlComponent.GetText()
			res, _ := utils.MakeRequest(url)
			a.response.Component.SetText(string(res))
		case tcell.KeyCtrlN:
			length := a.requestsList.Component.GetItemCount()
			a.requestsList.AddItem(&components.Request{
				Url:     "",
				Name:    "New request",
				Method:  "GET",
				Headers: "",
				Body:    "",
			}, func() {
				name, url := a.requestsList.Component.GetItemText(length)
				a.request.UrlComponent.SetText(url)
				a.request.NameComponent.SetText(name)
				a.tview.SetFocus(a.request.NameComponent)
			})
			a.requestsList.Component.SetCurrentItem(length)
		}

		return event
	})

	a.request.NameComponent.SetChangedFunc(func(name string) {
		index := a.requestsList.Component.GetCurrentItem()
		_, url := a.requestsList.Component.GetItemText(index)
		a.requestsList.Component.SetItemText(index, name, url)
	})

	a.request.UrlComponent.SetChangedFunc(func(url string) {
		index := a.requestsList.Component.GetCurrentItem()
		main, _ := a.requestsList.Component.GetItemText(index)
		a.requestsList.Component.SetItemText(index, main, url)
	})
}
