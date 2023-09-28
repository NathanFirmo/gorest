package gorest

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/NathanFirmo/gorest/internal/utils"
	"github.com/gdamore/tcell/v2"
)

const (
	KeyI = 256
)

func (app *App) SetInputHandlers() {
	// Sets an functions that intercepts keyboard inputs and allows to check pressed key and handle focus
	app.rootContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focus := app.tview.GetFocus()

		switch event.Key() {
		case tcell.KeyCtrlL:
			if focus == app.requestsList.Component.Component || focus == app.requestsList.Container {
				app.tview.SetFocus(app.request.UrlInputComponent)
			} else if focus == app.request.UrlInputComponent || focus == app.request.OptionsComponent {
				app.tview.SetFocus(app.response.Container)
			}
			return event
		case tcell.KeyCtrlH:
			if focus == app.response.Container {
				app.tview.SetFocus(app.request.UrlInputComponent)
			} else if focus == app.request.UrlInputComponent || focus == app.request.OptionsComponent {
				app.tview.SetFocus(app.requestsList.Component.Component)
			}
			return event
		case tcell.KeyCtrlJ:
			if focus == app.request.UrlInputComponent {
				app.tview.SetFocus(app.request.OptionsComponent)
			}
			return event
		case tcell.KeyCtrlK:
			if focus == app.request.OptionsComponent {
				app.tview.SetFocus(app.request.UrlInputComponent)
			}
			return event
		case tcell.KeyEnter:
			url := app.request.UrlInputComponent.GetText()
			res, _ := utils.MakeRequest(url)
			app.response.Component.SetText(string(res))
		case tcell.KeyCtrlN:
			length := app.requestsList.Component.Component.GetItemCount()
			app.requestsList.Component.AddItem(&components.Request{
				Url:     "",
				Name:    "New request",
				Method:  "GET",
				Headers: "",
				Body:    "",
			}, func() {
				_, url := app.requestsList.Component.Component.GetItemText(length)
				app.request.UrlInputComponent.SetText(url)
				app.tview.SetFocus(app.request.UrlInputComponent)
			})
			app.requestsList.Component.Component.SetCurrentItem(length)
		}

		return event
	})

	app.request.UrlInputComponent.SetChangedFunc(func(url string) {
		index := app.requestsList.Component.Component.GetCurrentItem()
		main, _ := app.requestsList.Component.Component.GetItemText(index)
		app.requestsList.Component.Component.SetItemText(index, main, url)
	})
}
