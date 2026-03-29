package screens

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg, width, height int) (Screen, tea.Cmd)
	View(width, height int) string
}
