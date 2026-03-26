package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type HomeScreen struct{}

func NewHomeScreen() Screen {
	return HomeScreen{}
}

func (h HomeScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	return h, nil
}

func (h HomeScreen) View(width, height int) string {
	return fmt.Sprintf(
		"Fullscreen TUI\n\nWidth: %d\ntest\ntest\nHeight: %d\n\nPress q to quit",
		width,
		height,
	)
}
