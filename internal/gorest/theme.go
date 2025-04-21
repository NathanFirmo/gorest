package gorest

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Colors based on Dracula theme
var (
	colorBackground = tcell.GetColor("#282a36") // Dark background
	colorForeground = tcell.GetColor("#f8f8f2") // Light foreground
	colorSelection  = tcell.GetColor("#44475a") // Selection background
	colorComment    = tcell.GetColor("#6272a4") // Comment/secondary text
	colorPurple     = tcell.GetColor("#bd93f9") // Purple
	colorPink       = tcell.GetColor("#ff79c6") // Pink
	colorGreen      = tcell.GetColor("#50fa7b") // Green
	colorOrange     = tcell.GetColor("#ffb86c") // Orange
	colorRed        = tcell.GetColor("#ff5555") // Red
	colorYellow     = tcell.GetColor("#f1fa8c") // Yellow
	colorCyan       = tcell.GetColor("#8be9fd") // Cyan
)

// ApplyDraculaTheme applies a Dracula-inspired theme to the tview application
func (a *App) ApplyDraculaTheme() {
	// Create a new theme
	theme := tview.Theme{
		PrimitiveBackgroundColor:    colorBackground,
		ContrastBackgroundColor:     colorSelection,
		MoreContrastBackgroundColor: colorBackground,
		BorderColor:                 colorPurple,
		TitleColor:                  colorCyan,
		GraphicsColor:               colorPurple,
		PrimaryTextColor:            colorForeground,
		SecondaryTextColor:          colorGreen,
		TertiaryTextColor:           colorYellow,
		InverseTextColor:            colorBackground,
		ContrastSecondaryTextColor:  colorPink,
	}

	// Apply the theme to the application
	tview.Styles = theme

	// Apply consistent style to each flex container
	a.rootContainer.SetBackgroundColor(colorBackground)
	a.requestsList.Container.SetBackgroundColor(colorBackground)
	a.requestsList.Container.SetBorderColor(colorPurple)
	a.requestsList.Container.SetTitleColor(colorCyan)

	a.request.Container.SetBackgroundColor(colorBackground)
	a.request.Container.SetBorderColor(colorPurple)
	a.request.Container.SetTitleColor(colorCyan)

	a.response.Container.SetBackgroundColor(colorBackground)
	a.response.Container.SetBorderColor(colorPurple)
	a.response.Container.SetTitleColor(colorCyan)

	// Style the requestsList component
	a.requestsList.Component.SetBackgroundColor(colorBackground)
	a.requestsList.Component.SetSelectedBackgroundColor(colorSelection)
	a.requestsList.Component.SetSelectedTextColor(colorGreen)
	a.requestsList.Component.SetMainTextColor(colorForeground)
	a.requestsList.Component.SetSecondaryTextColor(colorComment)

	// Style input fields
	a.request.NameComponent.SetBackgroundColor(colorBackground)
	a.request.NameComponent.SetFieldBackgroundColor(colorBackground)
	a.request.NameComponent.SetFieldTextColor(colorYellow)
	a.request.NameComponent.SetBorderColor(colorPurple)
	a.request.NameComponent.SetTitleColor(colorCyan)
	a.request.NameComponent.SetLabelColor(colorForeground)

	a.request.UrlComponent.SetBackgroundColor(colorBackground)
	a.request.UrlComponent.SetFieldBackgroundColor(colorBackground)
	a.request.UrlComponent.SetFieldTextColor(colorGreen)
	a.request.UrlComponent.SetBorderColor(colorPurple)
	a.request.UrlComponent.SetTitleColor(colorCyan)
	a.request.UrlComponent.SetLabelColor(colorForeground)

	// Style text areas with the exact same background color
	headerStyle := tcell.StyleDefault.
		Background(colorBackground).
		Foreground(colorCyan)
	a.request.HeadersComponent.SetTextStyle(headerStyle)
	a.request.HeadersComponent.SetBorderColor(colorPurple)
	a.request.HeadersComponent.SetTitleColor(colorCyan)

	bodyStyle := tcell.StyleDefault.
		Background(colorBackground).
		Foreground(colorOrange)
	a.request.BodyComponent.SetTextStyle(bodyStyle)
	a.request.BodyComponent.SetBorderColor(colorPurple)
	a.request.BodyComponent.SetTitleColor(colorCyan)

	// Style response component
	a.response.Component.SetBackgroundColor(colorBackground)
	a.response.Component.SetTextColor(colorForeground)
	a.response.Component.SetBorderColor(colorPurple)
	a.response.Component.SetTitleColor(colorCyan)
}
