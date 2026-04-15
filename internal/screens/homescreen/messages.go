package homescreen

import (
	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/screens"
	types "github.com/alex-305/ticktui/pkg/tickticktypes"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ActiveTaskListMsg struct {
	tasks []*types.Task
	err   error
}

type CompletedTaskListMsg struct {
	tasks []*types.Task
	err   error
}

type ProjectsLoadedMsg struct {
	projects []*types.Project
	err      error
}

type ActionCompletedMsg struct {
	err error
}

type ShowLoadingMsg struct{}

func (h *HomeScreen) handleMessages(msg tea.Msg, width, height int) (*HomeScreen, tea.Cmd, bool) {

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var spinCmd tea.Cmd
		h.loadingSpinner, spinCmd = h.loadingSpinner.Update(msg)
		return h, spinCmd, true
	case ShowLoadingMsg:
		if h.activeLoading || h.completedLoading {
			h.showLoadingSpinner = true
		}
		return h, nil, true
	case ActiveTaskListMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil, true
		}
		innerWidth := h.tabs.GetWindowWidth(width)
		innerHeight := h.tabs.GetWindowHeight(height)
		h.activeTaskTable = components.NewTaskTable(msg.tasks, innerWidth, (innerHeight/2)-10)
		h.activeLoading = false
		h.activeLoaded = true

		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil, true
	case CompletedTaskListMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil, true
		}
		innerWidth := h.tabs.GetWindowWidth(width)
		innerHeight := h.tabs.GetWindowHeight(height)
		h.completedTaskTable = components.NewTaskTable(msg.tasks, innerWidth, (innerHeight/2)-10)
		h.completedLoading = false
		h.completedLoaded = true

		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil, true

	case ProjectsLoadedMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil, true
		}
		h.projects = msg.projects
		lenProjects := len(h.projects)
		ids := make([]string, lenProjects)

		for i, p := range h.projects {
			ids[i] = p.ID
		}

		var projectNames []string
		for i, p := range h.projects {
			ids[i] = p.ID
			projectNames = append(projectNames, p.Name)
		}
		h.projectIDs = ids

		h.tabs.SetItems(projectNames)
		h.tabs.SetActive(h.activeProject)

		h, c := h.fetchAllTasks()
		return h, c, true

	case screens.GoBackScreenMsg:
		h, c := h.fetchAllTasks()
		return h, c, true

	case ActionCompletedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.activeLoading = false
			return h, nil, true
		}
		h, c := h.fetchAllTasks()
		return h, c, true
	}
	return h, nil, false

}
