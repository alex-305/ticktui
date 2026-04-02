package homescreen

import (
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/screens/taskeditscreen"
	tea "github.com/charmbracelet/bubbletea"
)

func (h *HomeScreen) handleKeyMsg(msg tea.KeyMsg) (*HomeScreen, tea.Cmd, bool) {

	switch msg.String() {
	case "tab":
		if h.focus == FocusActive {
			h.focus = FocusCompleted
		} else {
			h.focus = FocusActive
		}
		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil, true
	case "l":
		if h.activeProject < len(h.projects)-1 {
			h.activeProject++
			h.tabs.SetActive(h.activeProject)
			h.activeLoading = true
			h.activeLoaded = false
			return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID), true
		}
	case "h":
		if h.activeProject > 0 {
			h.activeProject--
			h.tabs.SetActive(h.activeProject)
			h.activeLoading = true
			h.activeLoaded = false
			return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID), true
		}
	case "r":
		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		s, c := h.fetchAllTasks()
		return s, c, true
	case "n":
		return h, func() tea.Msg {
			return screens.ChangeScreenMsg{
				NewScreen: taskeditscreen.NewTaskEditScreen(h.ctx, h.projects[h.activeProject].ID, nil)}
		}, true
	case "x":
		selectedTask, ok := h.getFocusedTable().GetSelectedTask()
		if !ok {
			return h, nil, true
		}
		return h, h.deleteTaskCmd(selectedTask), true
	case "c":
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
