package uicontainers

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/rivo/tview"
)

type ResponseContainer struct {
	Container *tview.Flex
	Component *tview.TextView
}

// Create a flex container with row alignment to put response data
func CreateResponse() *ResponseContainer {
	component := components.CreateResponse()

	container := tview.NewFlex()
	container.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("Response")

	container.AddItem(component, 0, 1, false)

	return &ResponseContainer{
		Container: container,
		Component: component,
	}
}
