package authscreen

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/asciiart"
	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
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

	h, c, ok := s.handleMessages(msg)
	if ok {
		return h, c
	}

	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if isKeyMsg {
		h, c, ok := s.handleKeyMsg(keyMsg)
		if ok {
			return h, c
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *AuthScreen) View(width, height int) string {
	if s.err != nil {
		return components.NewErrorBox(s.err, width, height).View()
	}
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	welcomeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	containerStyle := lipgloss.NewStyle().Align(lipgloss.Center)

	content := welcomeStyle.Render("Welcome to") + "\n"
	content += logoStyle.Render(asciiart.Logo) + "\n\n"

	content += fmt.Sprintf(
		"%s\n\nPress [Space] to open TickTick Auth in your browser.\n",
		titleStyle.Render("Authentication Required"),
	)

	if s.submitting {
		content += "\n\nExchanging code for token..."
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, containerStyle.Render(content))
}
