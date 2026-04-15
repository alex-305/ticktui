package app

import (
	"github.com/alex-305/ticktui/internal/screens"
	tea "github.com/charmbracelet/bubbletea"
)

type TokenExchangedMsg struct {
	err error
}

func (m *Model) handleMessages(msg tea.Msg) (*Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case screens.ChangeScreenMsg:
		m.history = append(m.history, m.current)
		m.current = msg.NewScreen
		return m, m.current.Init(), true
	case screens.ChangeScreenMsgNoHistory:
		m.history = []screens.Screen{}
		m.current = msg.NewScreen
		return m, m.current.Init(), true
	case screens.GoBackScreenMsg:
		if len(m.history) > 0 {
			lastIndex := len(m.history) - 1
			lastPage := m.history[lastIndex]
			m.history = m.history[:lastIndex]

			m.current = lastPage
			return m, m.current.Init(), true
		}
		return m, nil, true
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}
	return m, nil, false
}
