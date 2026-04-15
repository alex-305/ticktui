package homescreen

import (
	"github.com/alex-305/ticktui/internal/components"
	tea "github.com/charmbracelet/bubbletea"
)

func (h *HomeScreen) getFocusedTable() *components.TaskTable {
	if h.focus == FocusActive {
		return &h.activeTaskTable
	}
	return &h.completedTaskTable
}

func (h *HomeScreen) getUnfocusedTable() *components.TaskTable {
	if h.focus != FocusActive {
		return &h.activeTaskTable
	}
	return &h.completedTaskTable
}

func (h *HomeScreen) fetchAllTasks() (*HomeScreen, tea.Cmd) {
	if len(h.projects) == 0 {
		return h, nil
	}

	h.activeLoading = true
	h.completedLoading = true

	return h, tea.Batch(
		h.fetchProjectsAndTasks(h.projects[h.activeProject].ID, h.projectIDs),
		h.showLoadingCmd(),
	)
}

func (h *HomeScreen) fetchProjectsAndTasks(projectID string, allProjectIDs []string) tea.Cmd {
	return tea.Batch(
		h.fetchCompletedTasksCmd(allProjectIDs),
		h.fetchActiveTasksCmd(projectID),
	)
}
