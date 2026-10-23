package gorest

import (
	"fmt"
	"slices"

	"github.com/NathanFirmo/gorest/internal/curl"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Sets an functions that intercepts keyboard inputs and allows to check pressed key and handle focus
func (a *App) SetInputHandlers() {
	a.tview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if name, _ := a.pages.GetFrontPage(); name != mainPage {
			return event
		}

		focus := a.tview.GetFocus()
		onList := focus == a.requestsList.Component

		switch event.Key() {
		case tcell.KeyCtrlL:
			a.focusRight(focus)
		case tcell.KeyCtrlH:
			a.focusLeft(focus)
		case tcell.KeyCtrlJ:
			a.focusSibling(focus, 1)
		case tcell.KeyCtrlK:
			a.focusSibling(focus, -1)
		case tcell.KeyCtrlSpace, tcell.KeyCtrlR:
			a.sendRequest()
		case tcell.KeyCtrlN:
			if err := a.newRequest(); err != nil {
				a.setMessage(a.theme.Red, err.Error())
			}
			a.tview.SetFocus(a.request.NameComponent)
		case tcell.KeyCtrlT:
			a.cycleTheme()
		case tcell.KeyCtrlO:
			a.showHeaders = !a.showHeaders
			a.renderResponse()
		case tcell.KeyDelete:
			if !onList {
				return event
			}
			a.confirmDelete()
		case tcell.KeyRune:
			if !onList || event.Rune() != 'd' {
				return event
			}
			a.confirmDelete()
		default:
			return event
		}

		return nil
	})

	list := a.requestsList.Component
	list.SetChangedFunc(func(index int, _, _ string, _ rune) {
		a.loadRequest(index)
	})
	list.SetSelectedFunc(func(int, string, string, rune) {
		a.tview.SetFocus(a.request.NameComponent)
	})

	a.request.NameComponent.SetChangedFunc(func(string) { a.saveCurrentRequest() })
	a.request.UrlComponent.SetChangedFunc(func(string) { a.saveCurrentRequest() })
	a.request.MethodComponent.SetSelectedFunc(func(string, int) { a.saveCurrentRequest() })
	a.request.HeadersComponent.SetChangedFunc(a.saveCurrentRequest)
	a.request.BodyComponent.SetChangedFunc(a.saveCurrentRequest)
	a.request.UrlComponent.SetCurlHandler(a.importCurl)

	a.highlightOnFocus(list.Box, a.requestsList.Container.Box, nil)
	a.highlightOnFocus(a.response.Component.Box, a.response.Container.Box, nil)
	for field, box := range map[tview.Primitive]*tview.Box{
		a.request.NameComponent:    a.request.NameComponent.Box,
		a.request.MethodComponent:  a.request.MethodComponent.Box,
		a.request.UrlComponent:     a.request.UrlComponent.Box,
		a.request.HeadersComponent: a.request.HeadersComponent.Box,
		a.request.BodyComponent:    a.request.BodyComponent.Box,
	} {
		a.highlightOnFocus(box, box, func() { a.lastRequestField = field })
	}
}

func (a *App) highlightOnFocus(target, frame *tview.Box, onFocus func()) {
	target.SetFocusFunc(func() {
		frame.SetBorderColor(a.theme.Pink)
		if onFocus != nil {
			onFocus()
		}
	})
	target.SetBlurFunc(func() { frame.SetBorderColor(a.theme.Purple) })
}

func (a *App) focusRight(focus tview.Primitive) {
	if focus == a.requestsList.Component {
		a.tview.SetFocus(a.lastRequestField)
	} else if slices.Contains(a.request.Fields(), focus) {
		a.tview.SetFocus(a.response.Component)
	}
}

func (a *App) focusLeft(focus tview.Primitive) {
	if focus == a.response.Component {
		a.tview.SetFocus(a.lastRequestField)
	} else if slices.Contains(a.request.Fields(), focus) {
		a.tview.SetFocus(a.requestsList.Component)
	}
}

func (a *App) focusSibling(focus tview.Primitive, step int) {
	fields := a.request.Fields()
	index := slices.Index(fields, focus)
	if index < 0 {
		return
	}
	next := index + step
	if next >= 0 && next < len(fields) {
		a.tview.SetFocus(fields[next])
	}
}

func (a *App) cycleTheme() {
	a.ApplyTheme(nextTheme(a.theme))
	if err := a.store.SetSetting(themeKey, a.theme.Name); err != nil {
		a.setMessage(a.theme.Red, "Could not save theme: "+err.Error())
		return
	}
	a.setMessage(a.theme.Cyan, "Theme: "+a.theme.Name)
}

func (a *App) importCurl(req curl.Request, err error) {
	if err != nil {
		a.setMessage(a.theme.Red, "Invalid curl command: "+err.Error())
		return
	}

	a.loading = true
	a.request.SetMethod(req.Method)
	a.request.UrlComponent.SetText(req.URL)
	a.request.HeadersComponent.SetText(req.HeaderText(), false)
	a.request.BodyComponent.SetText(req.Body, false)
	a.loading = false

	a.saveCurrentRequest()
	a.setMessage(a.theme.Green, fmt.Sprintf("Imported curl: %s %s", req.Method, req.URL))
}

func (a *App) confirmDelete() {
	_, current, ok := a.requestsList.Current()
	if !ok {
		return
	}

	t := a.theme
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Delete request %q?", current.Name)).
		AddButtons([]string{"Cancel", "Delete"}).
		SetBackgroundColor(t.Selection).
		SetTextColor(t.Foreground).
		SetButtonStyle(tcell.StyleDefault.Background(t.Background).Foreground(t.Foreground)).
		SetButtonActivatedStyle(tcell.StyleDefault.Background(t.Red).Foreground(t.Background)).
		SetDoneFunc(func(_ int, label string) {
			a.pages.RemovePage(confirmPage)
			if label == "Delete" {
				a.deleteCurrentRequest()
			}
			a.tview.SetFocus(a.requestsList.Component)
		})
	modal.SetBorderColor(t.Red)
	a.pages.AddPage(confirmPage, modal, true, true)
}
