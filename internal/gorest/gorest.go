package gorest

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/rivo/tview"
)

// Gorest application
type App struct {
	tview         *tview.Application
	rootContainer *tview.Flex
	requestsList  *components.RequestsListComponent
	request       *components.RequestComponent
	response      *components.ResponseComponent
	mode          int
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

	// Configure root layout
	a.rootContainer.AddItem(a.requestsList.Container, 0, 1, true).
		AddItem(a.request.Container, 0, 2, false).
		AddItem(a.response.Container, 0, 2, false)
	a.tview.SetRoot(a.rootContainer, true).SetFocus(a.rootContainer).EnableMouse(true)

	a.tview.SetFocus(a.requestsList.Component)
	a.SetInputHandlers()

	return &a
}

// Start gorest application
func (a App) Start() error {
	return a.tview.Run()
}
