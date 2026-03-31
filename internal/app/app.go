package app

import (
	"fmt"
	"os"

	"github.com/alex-305/ticktui/internal/config"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/screens/authscreen"
	"github.com/alex-305/ticktui/internal/screens/homescreen"
	api "github.com/alex-305/ticktui/pkg/ticktickapi"
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

	token, err := config.LoadToken()
	var initialScreen screens.Screen
	client, err2 := api.GetClient(token)

	ctx := context.AppContext{
		APIClient: client,
	}

	if err != nil || err2 != nil {

		initialScreen = authscreen.NewAuthScreen(ctx)
	} else {
		initialScreen = homescreen.NewHomeScreen(ctx)
	}

	return &Model{
		current: initialScreen,
		history: []screens.Screen{},
		ctx:     ctx,
	}
}

func (m *Model) Init() tea.Cmd {
	if m.current != nil {
		return m.current.Init()
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screens.ChangeScreenMsg:
		m.history = append(m.history, m.current)
		m.current = msg.NewScreen
		return m, m.current.Init()
	case screens.ChangeScreenMsgNoHistory:
		m.history = []screens.Screen{}
		m.current = msg.NewScreen
		return m, m.current.Init()
	case screens.GoBackScreenMsg:
		if len(m.history) > 0 {
			lastIndex := len(m.history) - 1
			lastPage := m.history[lastIndex]
			m.history = m.history[:lastIndex]

			m.current = lastPage
			return m, m.current.Init()
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

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
