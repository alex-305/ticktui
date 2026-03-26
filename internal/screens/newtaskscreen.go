package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type CreateTaskScreen struct{}

func NewCreateTaskScreen() Screen {
	return CreateTaskScreen{}
}

func (h CreateTaskScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	return h, nil
}

func (h CreateTaskScreen) View(width, height int) string {
	return "New Task Screen"
}
