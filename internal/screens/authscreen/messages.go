package authscreen

import tea "github.com/charmbracelet/bubbletea"

type TokenExchangedMsg struct {
	err error
}

func (s *AuthScreen) handleMessages(msg tea.Msg) (*AuthScreen, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case TokenExchangedMsg:
		if msg.err != nil {
			s.err = msg.err
			s.submitting = false
			return s, nil, true
		}
		return s, s.OnSuccessfulAuthCmd(), true
	}
	return s, nil, false

}
