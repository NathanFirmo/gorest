package uicontainers

import (
	"github.com/NathanFirmo/gorest/internal/components"
	"github.com/rivo/tview"
)

type RequestContainer struct {
	Container       *tview.Flex
	OptionsComponent  *tview.TextArea
	UrlInputComponent *tview.InputField
}

// Create a flex container with row alignment to put request option components
func CreateRequest() *RequestContainer {
	container := tview.NewFlex()
	container.SetDirection(tview.FlexRow).SetBorder(true).SetTitle("Request")

	reqOpt := components.CreateRequestOptionsInput()
	reqUrlInput := components.CreateUrlInput()

	// Add requests component on its container
	container.
		AddItem(reqUrlInput, 0, 1, false).
		AddItem(reqOpt, 0, 9, false)

	return &RequestContainer{
		Container:       container,
		OptionsComponent:  reqOpt,
		UrlInputComponent: reqUrlInput,
	}
}
