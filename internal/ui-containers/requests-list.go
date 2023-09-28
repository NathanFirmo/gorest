package uicontainers

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/rivo/tview"
)

type RequestsListContainer struct {
	Container *tview.Flex
	Component *components.RequestsListComponent
}

// Create a flex container with row alignment to put requests list
func CreateRequestsList() *RequestsListContainer {
	reqList := components.CreateRequestsList()

	container := tview.NewFlex()
	container.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("List")

	container.AddItem(reqList.Component, 0, 1, false)

	return &RequestsListContainer{
		Container: container,
		Component: reqList,
	}
}
