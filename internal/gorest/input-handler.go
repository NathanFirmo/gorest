package gorest

import (
	"strings"

	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/NathanFirmo/gorest/internal/db"
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
		case tcell.KeyCtrlSpace:
			split := strings.Split(a.request.UrlComponent.GetText(), " ")
			if len(split) == 1 {
				return event
			}

			res, _ := utils.MakeRequest(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), a.request.HeadersComponent.GetText(), a.request.BodyComponent.GetText())
			a.response.Component.SetText(string(res))
		case tcell.KeyCtrlN:
			length := a.requestsList.Component.GetItemCount()
			req := &components.Request{
				ID:      int64(length),
				Url:     "",
				Name:    "New request",
				Method:  "GET",
				Headers: "",
				Body:    "",
			}
			a.requestsList.AddItem(req, func() {
				a.request.UrlComponent.SetText(req.Url)
				a.request.NameComponent.SetText(req.Name)
				a.tview.SetFocus(a.request.NameComponent)
			})
			a.requestsList.Component.SetCurrentItem(length)

			dbReq := &db.Request{
				ID:      int64(length),
				URL:     req.Url,
				Name:    req.Name,
				Method:  req.Method,
				Headers: req.Headers,
				Body:    req.Body,
			}

			if _, err := db.SaveRequest(dbReq); err != nil {
				panic(err)
			}
		}

		return event
	})

	a.request.NameComponent.SetChangedFunc(func(name string) {
		index := a.requestsList.Component.GetCurrentItem()
		_, url := a.requestsList.Component.GetItemText(index)
		a.requestsList.Component.SetItemText(index, name, url)

		req := &db.Request{
			ID:      int64(index),
			URL:     url,
			Name:    name,
			Method:  "GET", 
			Headers: a.request.HeadersComponent.GetText(),
			Body:    a.request.BodyComponent.GetText(),
		}

		if err := db.UpdateRequest(req); err != nil {
			panic(err)
		}
	})

	a.request.UrlComponent.SetChangedFunc(func(url string) {
		index := a.requestsList.Component.GetCurrentItem()
		name, _ := a.requestsList.Component.GetItemText(index)
		a.requestsList.Component.SetItemText(index, name, url)

		req := &db.Request{
			ID:      int64(index),
			URL:     url,
			Name:    name,
			Method:  "GET", 
			Headers: a.request.HeadersComponent.GetText(),
			Body:    a.request.BodyComponent.GetText(),
		}
		if err := db.UpdateRequest(req); err != nil {
			panic(err)
		}
	})
}
