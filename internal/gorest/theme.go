package gorest

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	Name       string
	Background tcell.Color
	Foreground tcell.Color
	Selection  tcell.Color
	Muted      tcell.Color
	Purple     tcell.Color
	Pink       tcell.Color
	Green      tcell.Color
	Orange     tcell.Color
	Red        tcell.Color
	Yellow     tcell.Color
	Cyan       tcell.Color
}

func newTheme(name string, colors ...string) Theme {
	c := make([]tcell.Color, len(colors))
	for i, hex := range colors {
		c[i] = tcell.GetColor(hex)
	}
	return Theme{
		Name:       name,
		Background: c[0],
		Foreground: c[1],
		Selection:  c[2],
		Muted:      c[3],
		Purple:     c[4],
		Pink:       c[5],
		Green:      c[6],
		Orange:     c[7],
		Red:        c[8],
		Yellow:     c[9],
		Cyan:       c[10],
	}
}

// Palette order: background, foreground, selection, muted, purple, pink, green, orange, red, yellow, cyan.
var Themes = []Theme{
	newTheme("dracula", "#282a36", "#f8f8f2", "#44475a", "#6272a4", "#bd93f9", "#ff79c6", "#50fa7b", "#ffb86c", "#ff5555", "#f1fa8c", "#8be9fd"),
	newTheme("catppuccin", "#1e1e2e", "#cdd6f4", "#45475a", "#7f849c", "#cba6f7", "#f5c2e7", "#a6e3a1", "#fab387", "#f38ba8", "#f9e2af", "#89dceb"),
	newTheme("tokyo-night", "#1a1b26", "#c0caf5", "#33467c", "#565f89", "#bb9af7", "#ff007c", "#9ece6a", "#ff9e64", "#f7768e", "#e0af68", "#7dcfff"),
	newTheme("gruvbox", "#282828", "#ebdbb2", "#504945", "#928374", "#d3869b", "#d3869b", "#b8bb26", "#fe8019", "#fb4934", "#fabd2f", "#8ec07c"),
	newTheme("nord", "#2e3440", "#eceff4", "#434c5e", "#7b88a1", "#b48ead", "#88c0d0", "#a3be8c", "#d08770", "#bf616a", "#ebcb8b", "#88c0d0"),
	newTheme("solarized-light", "#fdf6e3", "#586e75", "#eee8d5", "#93a1a1", "#6c71c4", "#d33682", "#859900", "#cb4b16", "#dc322f", "#b58900", "#2aa198"),
}

func themeByName(name string) Theme {
	for _, t := range Themes {
		if t.Name == name {
			return t
		}
	}
	return Themes[0]
}

func nextTheme(current Theme) Theme {
	for i, t := range Themes {
		if t.Name == current.Name {
			return Themes[(i+1)%len(Themes)]
		}
	}
	return Themes[0]
}

func (t Theme) MethodColor(method string) tcell.Color {
	switch method {
	case "GET":
		return t.Purple
	case "POST":
		return t.Green
	case "PUT":
		return t.Orange
	case "PATCH":
		return t.Yellow
	case "DELETE":
		return t.Red
	default:
		return t.Cyan
	}
}

func (t Theme) StatusColor(code int) tcell.Color {
	switch {
	case code >= 500:
		return t.Red
	case code >= 400:
		return t.Orange
	case code >= 300:
		return t.Cyan
	default:
		return t.Green
	}
}

// ApplyTheme applies the theme to the tview application
func (a *App) ApplyTheme(t Theme) {
	a.theme = t
	base := tcell.StyleDefault.Background(t.Background).Foreground(t.Foreground)

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    t.Background,
		ContrastBackgroundColor:     t.Selection,
		MoreContrastBackgroundColor: t.Background,
		BorderColor:                 t.Purple,
		TitleColor:                  t.Cyan,
		GraphicsColor:               t.Purple,
		PrimaryTextColor:            t.Foreground,
		SecondaryTextColor:          t.Green,
		TertiaryTextColor:           t.Yellow,
		InverseTextColor:            t.Background,
		ContrastSecondaryTextColor:  t.Pink,
	}

	for _, box := range a.boxes() {
		box.SetBackgroundColor(t.Background)
		box.SetBorderColor(t.Purple)
		box.SetTitleColor(t.Cyan)
	}

	list := a.requestsList.Component
	list.SetMainTextStyle(base)
	list.SetSecondaryTextStyle(base.Foreground(t.Muted))
	list.SetSelectedStyle(base.Background(t.Selection).Foreground(t.Foreground))

	placeholder := base.Foreground(t.Muted)
	for _, input := range []*tview.InputField{a.request.NameComponent, a.request.UrlComponent.InputField} {
		input.SetFieldStyle(base.Foreground(t.Green))
		input.SetPlaceholderStyle(placeholder)
	}

	method := a.request.MethodComponent
	method.SetFieldStyle(base.Foreground(t.Pink))
	method.SetFocusedStyle(base.Foreground(t.Pink).Bold(true))
	method.SetPrefixStyle(base.Foreground(t.Muted))
	method.SetListStyles(base.Background(t.Selection), base.Background(t.Purple).Foreground(t.Background))

	for textArea, color := range map[*tview.TextArea]tcell.Color{
		a.request.HeadersComponent: t.Cyan,
		a.request.BodyComponent:    t.Orange,
	} {
		textArea.SetTextStyle(base.Foreground(color))
		textArea.SetSelectedStyle(base.Background(t.Selection).Foreground(color))
		textArea.SetPlaceholderStyle(placeholder)
	}

	a.response.Component.SetTextColor(t.Foreground)
	a.statusBar.SetTextColor(t.Muted)

	a.requestsList.Refresh()
	a.renderResponse()
	a.renderStatus()

	// Re-focusing refreshes the highlighted border with the new colors.
	if focus := a.tview.GetFocus(); focus != nil {
		a.tview.SetFocus(focus)
	}
}

func (a *App) boxes() []*tview.Box {
	return []*tview.Box{
		a.rootContainer.Box,
		a.panels.Box,
		a.requestsList.Container.Box,
		a.requestsList.Component.Box,
		a.request.Container.Box,
		a.request.NameComponent.Box,
		a.request.MethodComponent.Box,
		a.request.UrlComponent.Box,
		a.request.HeadersComponent.Box,
		a.request.BodyComponent.Box,
		a.response.Container.Box,
		a.response.Component.Box,
		a.statusBar.Box,
	}
}
