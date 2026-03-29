package homescreen

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
	"github.com/alex-305/ticktui/internal/types"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

type Focus int

const (
	FocusActive Focus = iota
	FocusCompleted
)

type HomeScreen struct {
	ctx           context.AppContext
	projects      []types.Project
	projectIDs    []string
	activeProject int

	activeTaskTable    components.TaskTable
	completedTaskTable components.TaskTable
	paginator          paginator.Model
	focus              Focus

	activeLoaded     bool
	activeLoading    bool
	completedLoaded  bool
	completedLoading bool
	err              error
}

func NewHomeScreen(ctx context.AppContext) screens.Screen {
	return &HomeScreen{
		ctx:              ctx,
		activeLoaded:     false,
		activeLoading:    false,
		completedLoaded:  false,
		completedLoading: false,
	}
}

func (h *HomeScreen) Update(msg tea.Msg, width, height int) (screens.Screen, tea.Cmd) {
	var cmds []tea.Cmd

	if !h.activeLoaded && !h.activeLoading && h.err == nil {
		h.activeLoading = true
		return h, h.fetchProjectsCmd()
	}

	switch msg := msg.(type) {

	case ActiveTaskListMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil
		}
		h.activeTaskTable = components.NewTaskTable(msg.tasks, width)
		h.activeLoading = false
		h.activeLoaded = true

		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil
	case CompletedTaskListMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil
		}
		h.completedTaskTable = components.NewTaskTable(msg.tasks, width)
		h.completedLoading = false
		h.completedLoaded = true

		h.getFocusedTable().ApplyActiveStyle()
		h.getUnfocusedTable().ApplyInactiveStyle()
		return h, nil

	case ProjectsLoadedMsg:
		if msg.err != nil {
			h.err = msg.err
			return h, nil
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

		return h.fullFetch()

	case screens.GoBackScreenMsg:
		return h.fullFetch()
	case TaskDeletedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.activeLoading = false
			return h, nil
		}
		return h.fullFetch()
	case TaskCompletedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.activeLoading = false
			return h, nil
		}
		return h.fullFetch()

	case tea.KeyMsg:
		return h.handleKeyMsg(msg)
	}

	h.getFocusedTable().Update(msg)
	return h, tea.Batch(cmds...)
}

func (h *HomeScreen) View(width, height int) string {
	if h.err != nil {
		return components.NewErrorBox(h.err, width, height).View()
	}

	if len(h.projects) == 0 {
		return "\n  Initializing projects..."
	}

	var paginatorView string
	if len(h.projects) > 0 {
		paginatorView = h.paginator.View()
	} else {
		paginatorView = "\n"
	}

	return fmt.Sprintf(
		"Project: %s\n\n%s\n\n%s\n\nCompleted:\n%s\n\n%s",
		h.projects[h.activeProject].Name,
		h.activeTaskTable.View(),
		paginatorView,
		h.completedTaskTable.View(),
		"Controls: [Tab] Focus • [r] Refresh • [n] New Task • [x] Delete Task • [c] Complete Task",
	)
}
