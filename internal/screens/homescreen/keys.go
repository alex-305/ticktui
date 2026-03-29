package homescreen

import (
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/screens/createtaskscreen"
	tea "github.com/charmbracelet/bubbletea"
)

func (h *HomeScreen) handleKeyMsg(msg tea.KeyMsg) (screens.Screen, tea.Cmd) {

	switch msg.String() {
	case "tab":
		if h.focus == FocusActive {
			h.focus = FocusCompleted
		} else {
			h.focus = FocusActive
		}
		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil
	case "l":
		if h.activeProject < len(h.projects)-1 {
			h.activeProject++
			h.paginator.Page++
			h.activeLoading = true
			h.activeLoaded = false
			return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID)
		}
	case "h":
		if h.activeProject > 0 {
			h.activeProject--
			h.paginator.Page--
			h.activeLoading = true
			h.activeLoaded = false
			return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID)
		}
	case "r":
		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h.fullFetch()
	case "n":
		return h, func() tea.Msg {
			return screens.ChangeScreenMsg{NewScreen: createtaskscreen.NewCreateTaskScreen(h.ctx, h.projects[h.activeProject].ID)}
		}
	case "x":
		selectedTask, ok := h.getFocusedTable().GetSelectedTask()
		if !ok {
			return h, nil
		}
		return h, h.deleteTaskCmd(selectedTask)
	case "c":
		selectedTask, ok := h.activeTaskTable.GetSelectedTask()
		if !ok {
			return h, nil
		}

		return h, h.completeTaskCmd(selectedTask)
	}

	return h, nil
}
