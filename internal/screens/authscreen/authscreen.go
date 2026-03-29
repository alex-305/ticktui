package authscreen

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/api"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/screens/homescreen"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthScreen struct {
	ctx        context.AppContext
	textInput  textinput.Model
	err        error
	submitting bool
}

func NewAuthScreen(ctx context.AppContext) screens.Screen {
	ti := textinput.New()
	ti.Placeholder = "Paste your auth code here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return &AuthScreen{
		ctx:       ctx,
		textInput: ti,
	}
}

func (a *AuthScreen) Init() tea.Cmd {
	return nil
}

func (s *AuthScreen) Update(msg tea.Msg, width, height int) (screens.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			return s, func() tea.Msg {
				err := api.LaunchBrowserAndSaveAuthToken()
				return TokenExchangedMsg{err}
			}
		}
	case TokenExchangedMsg:
		if msg.err != nil {
			s.err = msg.err
			s.submitting = false
			return s, nil
		}
		return s, func() tea.Msg {
			return screens.ChangeScreenMsgNoHistory{NewScreen: homescreen.NewHomeScreen(s.ctx)}
		}
	}
	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *AuthScreen) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	content := fmt.Sprintf(
		"%s\n\nPress [Space] to open TickTick Auth in your browser.\n",
		titleStyle.Render("Authentication Required"),
	)

	if s.err != nil {
		content += "\n\n" + errorStyle.Render(fmt.Sprintf("Error: %v", s.err))
	}

	if s.submitting {
		content += "\n\nExchanging code for token..."
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
