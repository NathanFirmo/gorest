package gorest

import (
	uicontainers "github.com/NathanFirmo/gorest/internal/ui-containers"
	"github.com/rivo/tview"
)

// Gorest application
type App struct {
	tview         *tview.Application
	rootContainer *tview.Flex
	requestsList  *uicontainers.RequestsListContainer
	request       *uicontainers.RequestContainer
	response      *uicontainers.ResponseContainer
	mode          int
}

// Create gorest application and initialize components
func CreateApp() *App {
	requestContainer := uicontainers.CreateRequest()
	responseContainer := uicontainers.CreateResponse()
	requestsListContainer := uicontainers.CreateRequestsList()

	a := App{
		tview:         tview.NewApplication(),
		rootContainer: tview.NewFlex(),
		request:       requestContainer,
		response:      responseContainer,
		requestsList:  requestsListContainer,
	}

	// Configure root layout
	a.rootContainer.AddItem(a.requestsList.Container, 0, 1, true).
		AddItem(a.request.Container, 0, 2, false).
		AddItem(a.response.Container, 0, 2, false)
	a.tview.SetRoot(a.rootContainer, true).SetFocus(a.rootContainer).EnableMouse(true)

	a.tview.SetFocus(a.requestsList.Component.Component)
	a.SetInputHandlers()

	return &a
}

// Start gorest application
func (a App) Start() error {
	return a.tview.Run()
}
