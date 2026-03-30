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

func (h *HomeScreen) Init() tea.Cmd {
	h.activeLoading = true
	return h.fetchProjectsCmd()
}

func (h *HomeScreen) Update(msg tea.Msg, width, height int) (screens.Screen, tea.Cmd) {
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
		if !h.activeLoaded || !h.completedLoaded {
			return h, nil
		}
		s, c, ok := h.handleKeyMsg(msg)

		if ok {
			return s, c
		}
	}

	if !h.activeLoaded || !h.completedLoaded {
		return h, nil
	}

	return h, h.getFocusedTable().Update(msg)
}

func (h *HomeScreen) View(width, height int) string {
	if h.err != nil {
		return components.NewErrorBox(h.err, width, height).View()
	}

	if len(h.projects) == 0 || !h.activeLoaded || !h.completedLoaded {
		return "\n Loading tasks..."
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
