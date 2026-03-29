package screens

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

func (h *HomeScreen) fullFetch() (*HomeScreen, tea.Cmd) {
	h.activeLoading = true
	h.completedLoading = true
	h.activeLoaded = false
	h.completedLoaded = false

	return h, h.fetchAllData(h.projects[h.activeProject].ID, h.projectIDs)
}

func (h *HomeScreen) fetchAllData(projectID string, allProjectIDs []string) tea.Cmd {
	return tea.Batch(
		h.fetchCompletedTasksCmd(allProjectIDs),
		h.fetchActiveTasksCmd(projectID),
	)
}
