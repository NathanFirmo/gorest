package gorest

import (
	"context"
	"fmt"
	"time"

	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/NathanFirmo/gorest/internal/db"
	"github.com/NathanFirmo/gorest/internal/utils"
	"github.com/rivo/tview"
)

const (
	mainPage    = "main"
	confirmPage = "confirm"
	themeKey    = "theme"
	keyHints    = "^Space/^R send  ^N new  d delete  ^O headers/body  ^T theme  ^H/^J/^K/^L move  ^C quit"
)

type result struct {
	response *utils.Response
	err      error
	pending  bool
}

// Gorest application
type App struct {
	tview         *tview.Application
	pages         *tview.Pages
	rootContainer *tview.Flex
	panels        *tview.Flex
	requestsList  *components.RequestsListComponent
	request       *components.RequestComponent
	response      *components.ResponseComponent
	statusBar     *tview.TextView

	store            *db.Store
	theme            Theme
	loading          bool
	lastRequestField tview.Primitive
	results          map[int64]result
	showHeaders      bool
	cancelRequest    context.CancelFunc
	message          string
	messageTimer     *time.Timer
}

// Create gorest application and initialize components
func CreateApp(store *db.Store) (*App, error) {
	a := App{
		tview:         tview.NewApplication(),
		pages:         tview.NewPages(),
		rootContainer: tview.NewFlex(),
		panels:        tview.NewFlex(),
		request:       components.CreateRequest(),
		response:      components.CreateResponse(),
		requestsList:  components.CreateRequestsList(),
		statusBar:     tview.NewTextView(),
		store:         store,
		results:       make(map[int64]result),
	}
	a.lastRequestField = a.request.NameComponent
	a.requestsList.Format = a.formatListItem
	a.statusBar.SetDynamicColors(true)

	a.panels.AddItem(a.requestsList.Container, 0, 1, true).
		AddItem(a.request.Container, 0, 2, false).
		AddItem(a.response.Container, 0, 2, false)
	a.rootContainer.SetDirection(tview.FlexRow).
		AddItem(a.panels, 0, 1, true).
		AddItem(a.statusBar, 1, 0, false)
	a.pages.AddPage(mainPage, a.rootContainer, true, true)
	a.tview.SetRoot(a.pages, true).EnablePaste(true).EnableMouse(true)

	themeName, err := store.Setting(themeKey)
	if err != nil {
		return nil, err
	}
	a.ApplyTheme(themeByName(themeName))

	a.SetInputHandlers()
	if err := a.loadSavedRequests(); err != nil {
		return nil, err
	}
	a.tview.SetFocus(a.requestsList.Component)

	return &a, nil
}

// loads all saved requests from the database
func (a *App) loadSavedRequests() error {
	requests, err := a.store.ListRequests()
	if err != nil {
		return err
	}

	for _, req := range requests {
		a.requestsList.AddItem(components.Request(req))
	}

	if len(requests) == 0 {
		return a.newRequest()
	}
	a.loadRequest(0)
	return nil
}

func (a *App) formatListItem(r components.Request) (string, string) {
	method := r.Method
	switch method {
	case "DELETE":
		method = "DEL"
	case "OPTIONS":
		method = "OPT"
	}
	main := fmt.Sprintf("[%s::b]%-5s[-::-] %s", a.theme.MethodColor(r.Method).CSS(), method, tview.Escape(r.Name))
	return main, tview.Escape(r.URL)
}

func (a *App) loadRequest(index int) {
	if index < 0 || index >= len(a.requestsList.Items) {
		return
	}
	a.loading = true
	a.request.SetValue(a.requestsList.Items[index])
	a.loading = false
	a.renderResponse()
}

func (a *App) saveCurrentRequest() {
	if a.loading {
		return
	}
	index, current, ok := a.requestsList.Current()
	if !ok {
		return
	}

	req := a.request.Value(current.ID)
	a.requestsList.UpdateItem(index, req)
	if err := a.store.UpdateRequest(db.Request(req)); err != nil {
		a.setMessage(a.theme.Red, "Could not save request: "+err.Error())
	}
}

func (a *App) newRequest() error {
	req := db.Request{Name: "New request", Method: "GET"}
	if err := a.store.CreateRequest(&req); err != nil {
		return err
	}
	a.requestsList.AddItem(components.Request(req))
	a.requestsList.Component.SetCurrentItem(len(a.requestsList.Items) - 1)
	a.loadRequest(len(a.requestsList.Items) - 1)
	return nil
}

func (a *App) deleteCurrentRequest() {
	index, current, ok := a.requestsList.Current()
	if !ok {
		return
	}
	if err := a.store.DeleteRequest(current.ID); err != nil {
		a.setMessage(a.theme.Red, "Could not delete request: "+err.Error())
		return
	}
	delete(a.results, current.ID)
	a.requestsList.RemoveItem(index)

	if len(a.requestsList.Items) == 0 {
		if err := a.newRequest(); err != nil {
			a.setMessage(a.theme.Red, err.Error())
		}
		return
	}
	a.loadRequest(a.requestsList.Component.GetCurrentItem())
	a.setMessage(a.theme.Green, fmt.Sprintf("Deleted %q", current.Name))
}

func (a *App) sendRequest() {
	_, current, ok := a.requestsList.Current()
	if !ok {
		return
	}
	req := a.request.Value(current.ID)
	if req.URL == "" {
		a.setMessage(a.theme.Orange, "Missing URL")
		return
	}

	if a.cancelRequest != nil {
		a.cancelRequest()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelRequest = cancel
	a.results[req.ID] = result{pending: true}
	a.renderResponse()

	go func() {
		res, err := utils.MakeRequest(ctx, req.Method, req.URL, req.Headers, req.Body)
		a.tview.QueueUpdateDraw(func() {
			if ctx.Err() != nil && err != nil {
				return
			}
			a.results[req.ID] = result{response: res, err: err}
			a.renderResponse()
		})
	}()
}

func (a *App) renderResponse() {
	_, current, ok := a.requestsList.Current()
	if !ok {
		a.response.Show("Response", "")
		return
	}

	r, ok := a.results[current.ID]
	switch {
	case !ok:
		a.response.Show("Response", "")
	case r.pending:
		a.response.Show("Response", "Sending request...")
	case r.err != nil:
		a.response.Show(fmt.Sprintf("Response · [%s]Error[-]", a.theme.Red.CSS()), r.err.Error())
	default:
		res := r.response
		view, content := "body", res.FormattedBody()
		if a.showHeaders {
			view, content = "headers", res.FormattedHeaders()
		}
		title := fmt.Sprintf("Response · [%s::b]%s[-::-] · %s · %s · %s",
			a.theme.StatusColor(res.StatusCode).CSS(), tview.Escape(res.Status),
			utils.FormatDuration(res.Duration), utils.FormatSize(len(res.Body)), view)
		a.response.Show(title, content)
	}
}

func (a *App) setMessage(color interface{ CSS() string }, message string) {
	a.message = fmt.Sprintf("[%s]%s[-]", color.CSS(), tview.Escape(message))
	a.renderStatus()

	if a.messageTimer != nil {
		a.messageTimer.Stop()
	}
	a.messageTimer = time.AfterFunc(4*time.Second, func() {
		a.tview.QueueUpdateDraw(func() {
			a.message = ""
			a.renderStatus()
		})
	})
}

func (a *App) renderStatus() {
	text := a.message
	if text == "" {
		text = keyHints
	}
	a.statusBar.SetText(" " + text)
}

// Start gorest application
func (a *App) Start() error {
	return a.tview.Run()
}
