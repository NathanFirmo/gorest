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
				a.tview.SetFocus(a.request.UrlInputComponent)
			} else if focus == a.request.UrlInputComponent || focus == a.request.OptionsComponent {
				a.tview.SetFocus(a.response.Container)
			}
			return event
		case tcell.KeyCtrlH:
			if focus == a.response.Container {
				a.tview.SetFocus(a.request.UrlInputComponent)
			} else if focus == a.request.UrlInputComponent || focus == a.request.OptionsComponent {
				a.tview.SetFocus(a.requestsList.Component)
			}
			return event
		case tcell.KeyCtrlJ:
			if focus == a.request.UrlInputComponent {
				a.tview.SetFocus(a.request.OptionsComponent)
			}
			return event
		case tcell.KeyCtrlK:
			if focus == a.request.OptionsComponent {
				a.tview.SetFocus(a.request.UrlInputComponent)
			}
			return event
		case tcell.KeyEnter:
			url := a.request.UrlInputComponent.GetText()
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
				_, url := a.requestsList.Component.GetItemText(length)
				a.request.UrlInputComponent.SetText(url)
				a.tview.SetFocus(a.request.UrlInputComponent)
			})
			a.requestsList.Component.SetCurrentItem(length)
		}

		return event
	})

	a.request.UrlInputComponent.SetChangedFunc(func(url string) {
		index := a.requestsList.Component.GetCurrentItem()
		main, _ := a.requestsList.Component.GetItemText(index)
		a.requestsList.Component.SetItemText(index, main, url)
	})
}
