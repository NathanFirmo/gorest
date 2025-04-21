package gorest

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/NathanFirmo/gorest/internal/db"
	"github.com/rivo/tview"
)

// Gorest application
type App struct {
	tview         *tview.Application
	rootContainer *tview.Flex
	requestsList  *components.RequestsListComponent
	request       *components.RequestComponent
	response      *components.ResponseComponent
}

// Create gorest application and initialize components
func CreateApp() *App {
	req := components.CreateRequest()
	res := components.CreateResponse()
	reqList := components.CreateRequestsList()

	a := App{
		tview:         tview.NewApplication(),
		rootContainer: tview.NewFlex(),
		request:       req,
		response:      res,
		requestsList:  reqList,
	}

	a.rootContainer.AddItem(a.requestsList.Container, 0, 1, true).
		AddItem(a.request.Container, 0, 2, false).
		AddItem(a.response.Container, 0, 2, false)
	a.tview.SetRoot(a.rootContainer, true).SetFocus(a.rootContainer).EnableMouse(true)

	a.tview.SetFocus(a.requestsList.Component)
	a.SetInputHandlers()
	a.loadSavedRequests()

	return &a
}

// loadSavedRequests loads all saved requests from the database
func (a *App) loadSavedRequests() {
	requests, err := db.GetAllRequests()
	if err != nil {
		return
	}

	for i, req := range requests {
		a.requestsList.AddItem(&components.Request{
			ID:      int64(i),
			Url:     req.URL,
			Name:    req.Name,
			Method:  req.Method,
			Headers: req.Headers,
			Body:    req.Body,
		}, func() {
			a.request.UrlComponent.SetText(req.URL)
			a.request.NameComponent.SetText(req.Name)
			a.request.HeadersComponent.SetText(req.Headers, true)
			a.request.BodyComponent.SetText(req.Body, true)
			a.tview.SetFocus(a.request.NameComponent)
		})
	}
}

// Start gorest application
func (a App) Start() error {
	return a.tview.Run()
}
