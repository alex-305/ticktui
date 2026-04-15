package authscreen

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	OpenBrowser key.Binding
}

var DefaultKeyMap = KeyMap{
	OpenBrowser: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("[space]", "open browser"),
	),
}

func (s *AuthScreen) handleKeyMsg(msg tea.KeyMsg) (*AuthScreen, tea.Cmd, bool) {

	switch {
	case key.Matches(msg, DefaultKeyMap.OpenBrowser):
		if s.submitting {
			return s, nil, true
		}
		s.submitting = true
		return s, s.AuthOnBrowserCmd(), true
	}

	return s, nil, false
}
