package app

import (
	"fmt"
	"log"
	"os"

	"github.com/alex-305/ticktui/internal/api"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	history []screens.Screen
	current screens.Screen

	ctx context.AppContext
}

func NewModel() *Model {
	client, err := api.GetClient()

	if err != nil {
		log.Fatal("Client not working")
	}

	ctx := context.AppContext{
		TaskService: services.NewTaskService(client),
		APIClient:   client,
	}
	return &Model{
		current: screens.NewHomeScreen(ctx),
		history: []screens.Screen{},
		ctx:     ctx,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screens.ChangeScreenMsg:
		m.history = append(m.history, m.current)
		m.current = msg.NewScreen
		return m, nil
	case screens.GoBackScreenMsg:
		navigateBackAScreen(m)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m.updateScreen(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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

func (m *Model) updateScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	screen, cmd := m.current.Update(msg, m.width, m.height)
	m.current = screen
	return m, cmd
}

func (m *Model) View() string {
	return m.current.View(m.width, m.height)
}

func Run() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
