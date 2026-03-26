package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type HomeScreen struct {
	lastKey string
}

func NewHomeScreen() Screen {
	return HomeScreen{
		lastKey: "None yet",
	}
}

func (h HomeScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		h.lastKey = msg.String()
		switch msg.String() {
		case "n":
			// We caught "n"! Tell the main model to switch screens.
			cmd := func() tea.Msg {
				return ChangeScreenMsg{NewScreen: NewCreateTaskScreen()}
			}
			return h, cmd
		}
	}
	return h, nil
}

func (h HomeScreen) View(width, height int) string {
	return fmt.Sprintf(
		"Fullscreen TUI\n\nWidth: %d\n\nHeight: %d\n\nLast key: %s\nPress q to quit",
		width,
		height,
		h.lastKey,
	)
}
