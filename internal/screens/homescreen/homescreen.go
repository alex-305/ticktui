package homescreen

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
	types "github.com/alex-305/ticktui/pkg/tickticktypes"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Focus int

const (
	FocusActive Focus = iota
	FocusCompleted
)

type HomeScreen struct {
	ctx           context.AppContext
	projects      []*types.Project
	projectIDs    []string
	activeProject int

	activeTaskTable    components.TaskTable
	completedTaskTable components.TaskTable
	paginator          paginator.Model
	loadingSpinner     spinner.Model
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
		loadingSpinner:   spinner.New(),
	}
}

func (h *HomeScreen) Init() tea.Cmd {
	h.activeLoading = true
	h.loadingSpinner.Spinner = spinner.Dot

	return tea.Batch(
		h.fetchProjectsCmd(),
		h.loadingSpinner.Tick)
}

func (h *HomeScreen) Update(msg tea.Msg, width, height int) (screens.Screen, tea.Cmd) {
	h, c, ok := h.handleMessages(msg, width, height)
	if ok {
		return h, c
	}

	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if isKeyMsg {
		if !h.activeLoaded || !h.completedLoaded {
			return h, nil
		}
		h, c, ok := h.handleKeyMsg(keyMsg)
		if ok {
			return h, c
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
		return fmt.Sprintf("\n %s Loading tasks...", h.loadingSpinner.View())
	}

	var paginatorView string
	if len(h.projects) > 0 {
		paginatorView = h.paginator.View()
	} else {
		paginatorView = "\n"
	}

	var cKeyString string
	if h.focus == FocusActive {
		cKeyString = "Complete Task"
	} else {
		cKeyString = "Undo Completion"
	}

	return fmt.Sprintf(
		"Project: %s\n\n%s\n\n%s\n\nCompleted:\n%s\n\n%s",
		h.projects[h.activeProject].Name,
		h.activeTaskTable.View(),
		paginatorView,
		h.completedTaskTable.View(),
		fmt.Sprintf("Controls: [Tab] Focus • [r] Refresh • [n] New Task • [x] Delete Task • [c] %s", cKeyString),
	)
}
