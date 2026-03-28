package components

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

type ErrorBox struct {
	Err    error
	Width  int
	Height int
}

func NewErrorBox(err error, width, height int) ErrorBox {
	return ErrorBox{
		Err:    err,
		Width:  width,
		Height: height,
	}
}

func (e ErrorBox) View() string {
	if e.Err == nil {
		return ""
	}

	errorStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")). // Red
		Padding(1, 2).
		Width(e.Width - 4).
		Align(lipgloss.Center)

	content := fmt.Sprintf("ERROR OCCURRED\n\n%v\n\n[ctrl + c] Quit", e.Err)

	return lipgloss.Place(e.Width, e.Height, lipgloss.Center, lipgloss.Center, errorStyle.Render(content))

}
