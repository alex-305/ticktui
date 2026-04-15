package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Quit   key.Binding
	GoBack key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("[ctrl + c]", "quit app"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (*Model, tea.Cmd, bool) {

	switch {
	case key.Matches(msg, DefaultKeyMap.Quit):
		return m, tea.Quit, true
	case key.Matches(msg, DefaultKeyMap.Quit):
		if len(m.history) > 0 {
			lastIndex := len(m.history) - 1
			lastPage := m.history[lastIndex]
			m.history = m.history[:lastIndex]

			m.current = lastPage

			return m, nil, true
		}
	}

	return m, nil, false
}
