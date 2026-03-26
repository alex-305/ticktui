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

	history []screens.Screen
	current screens.Screen
}

func NewModel() Model {
	return Model{
		current: screens.NewHomeScreen(),
		history: []screens.Screen{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screens.ChangeScreenMsg:
		m.history = append(m.history, m.current)
		m.current = msg.NewScreen
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// forward to screen
		return m.updateScreen(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			if len(m.history) > 0 {
				lastIndex := len(m.history) - 1
				lastPage := m.history[lastIndex]
				m.history = m.history[:lastIndex]

				m.current = lastPage

				return m, nil
			}
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
