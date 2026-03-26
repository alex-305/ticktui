package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateTaskScreen struct {
	titleInput textinput.Model
	descInput  textarea.Model
	focusIndex int
}

func NewCreateTaskScreen() Screen {
	ti := textinput.New()
	ti.Placeholder = "Enter task title..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	ta := textarea.New()
	ta.Placeholder = "Enter task description..."
	ta.SetHeight(5)
	ta.SetWidth(40)

	return CreateTaskScreen{
		titleInput: ti,
		descInput:  ta,
		focusIndex: 0,
	}
}

func (h CreateTaskScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab:
			if msg.Type == tea.KeyShiftTab {
				h.focusIndex--
			} else {
				h.focusIndex++
			}

			if h.focusIndex > 1 {
				h.focusIndex = 0
			} else if h.focusIndex < 0 {
				h.focusIndex = 1
			}

			if h.focusIndex == 0 {
				cmd = h.titleInput.Focus()
				h.descInput.Blur()
			} else {
				h.titleInput.Blur()
				cmd = h.descInput.Focus()
			}
			cmds = append(cmds, cmd)

			return h, tea.Batch(cmds...)
		}
	}

	if h.focusIndex == 0 {
		h.titleInput, cmd = h.titleInput.Update(msg)
	} else {
		h.descInput, cmd = h.descInput.Update(msg)
	}
	cmds = append(cmds, cmd)

	return h, tea.Batch(cmds...)
}

func (h CreateTaskScreen) View(width, height int) string {
	return fmt.Sprintf(
		"New Task\n\nTitle:\n%s\n\nDescription:\n%s\n\n[Tab] Next field • [Shift+Tab] Prev field",
		h.titleInput.View(),
		h.descInput.View(),
	)
}
