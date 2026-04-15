package authscreen

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/asciiart"
	"github.com/alex-305/ticktui/internal/config"
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/screens/homescreen"
	api "github.com/alex-305/ticktui/pkg/ticktickapi"
	tea "github.com/charmbracelet/bubbletea"
)

func (s *AuthScreen) AuthOnBrowserCmd() tea.Cmd {
	return func() tea.Msg {
		token, err := api.LaunchBrowserAndSaveAuthToken(fmt.Sprintf("%s\n\nSuccessfully authenticated. You can now return to the comfort of your terminal :)", asciiart.Logo))
		config.SaveToken(token)

		return TokenExchangedMsg{err}
	}
}

func (s *AuthScreen) OnSuccessfulAuthCmd() tea.Cmd {

	return func() tea.Msg {
		token, err := config.LoadToken()
		if err != nil {
			s.err = err
		}
		freshClient, err := api.GetClient(token)

		if err != nil {
			s.err = err
		}
		s.ctx.APIClient = freshClient

		return screens.ChangeScreenMsgNoHistory{NewScreen: homescreen.NewHomeScreen(s.ctx)}
	}
}
