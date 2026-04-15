package homescreen

import (
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/screens/taskeditscreen"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	SwitchFocus        key.Binding
	SwitchProjectRight key.Binding
	SwitchProjectLeft  key.Binding
	Reload             key.Binding
	NewTask            key.Binding
	DeleteTask         key.Binding
	CompleteTask       key.Binding
}

var DefaultKeyMap = KeyMap{
	SwitchFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	SwitchProjectRight: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "switch project right"),
	),
	SwitchProjectLeft: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "switch project left"),
	),
	Reload: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh data"),
	),
	NewTask: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "create new task"),
	),
	DeleteTask: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete task"),
	),
	CompleteTask: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "complete task"),
	),
}

func (h *HomeScreen) handleKeyMsg(msg tea.KeyMsg) (*HomeScreen, tea.Cmd, bool) {

	switch {
	case key.Matches(msg, DefaultKeyMap.SwitchFocus):
		if h.focus == FocusActive {
			h.focus = FocusCompleted
		} else {
			h.focus = FocusActive
		}
		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil, true

	case key.Matches(msg, DefaultKeyMap.SwitchProjectRight):
		if h.activeProject < len(h.projects)-1 {
			h.activeProject++
			h.tabs.SetActive(h.activeProject)
			h.activeLoading = true
			h.activeLoaded = false
			return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID), true
		}

	case key.Matches(msg, DefaultKeyMap.SwitchProjectLeft):
		if h.activeProject > 0 {
			h.activeProject--
			h.tabs.SetActive(h.activeProject)
			h.activeLoading = true
			h.activeLoaded = false
			return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID), true
		}

	case key.Matches(msg, DefaultKeyMap.Reload):
		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		s, c := h.fetchAllTasks()
		return s, c, true
	case key.Matches(msg, DefaultKeyMap.NewTask):
		return h, func() tea.Msg {
			return screens.ChangeScreenMsg{
				NewScreen: taskeditscreen.NewTaskEditScreen(h.ctx, h.projects[h.activeProject].ID, nil)}
		}, true

	case key.Matches(msg, DefaultKeyMap.DeleteTask):
		selectedTask, ok := h.getFocusedTable().GetSelectedTask()
		if !ok {
			return h, nil, true
		}
		return h, h.deleteTaskCmd(selectedTask), true
	case key.Matches(msg, DefaultKeyMap.CompleteTask):
		selectedTask, ok := h.getFocusedTable().GetSelectedTask()
		if !ok {
			return h, nil, true
		}
		if h.focus == FocusActive {
			return h, h.completeTaskCmd(selectedTask), true
		} else {
			return h, h.decompleteTaskCmd(selectedTask), true
		}

	}

	return h, nil, false
}
