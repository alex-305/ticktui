package homescreen

import (
	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/types"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskDeletedMsg struct{ err error }

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

func (h *HomeScreen) handleMessages(msg tea.Msg, width, height int) (*HomeScreen, tea.Cmd, bool) {

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var spinCmd tea.Cmd
		h.loadingSpinner, spinCmd = h.loadingSpinner.Update(msg)
		return h, spinCmd, true
	case ActiveTaskListMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil, true
		}
		h.activeTaskTable = components.NewTaskTable(msg.tasks, width)
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
		h.completedTaskTable = components.NewTaskTable(msg.tasks, width)
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
		h.projectIDs = ids
		p := paginator.New()
		p.SetTotalPages(lenProjects)
		p.Type = paginator.Dots
		h.paginator = p

		h, c := h.fullFetch()
		return h, c, true

	case screens.GoBackScreenMsg:
		h, c := h.fullFetch()
		return h, c, true

	case ActionCompletedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.activeLoading = false
			return h, nil, true
		}
		h, c := h.fullFetch()
		return h, c, true
	}
	return h, nil, false

}
