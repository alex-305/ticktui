package app

import (
	"fmt"
	"os"

	"github.com/alex-305/ticktui/internal/screens"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	current screens.Screen
}

func NewModel() Model {
	return Model{
		current: screens.NewHomeScreen(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// forward to screen
		return m.updateScreen(msg)

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m.updateScreen(msg)
}

func (m Model) updateScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	screen, cmd := m.current.Update(msg, m.width, m.height)
	m.current = screen
	return m, cmd
}

func (m Model) View() string {
	return m.current.View(m.width, m.height)
}

func Run() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
