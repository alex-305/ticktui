package homescreen

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
	types "github.com/alex-305/ticktui/pkg/tickticktypes"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	tabs               components.Tabs

	loadingSpinner spinner.Model

	focus            Focus
	activeLoaded     bool
	activeLoading    bool
	completedLoaded  bool
	completedLoading bool
	err              error
}

func NewHomeScreen(ctx context.AppContext) screens.Screen {
	return &HomeScreen{
		ctx:              ctx,
		tabs:             components.NewTabs("h", "l"),
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
		if len(h.projects) == 0 {
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

	cKeyString := "Undo Completion"
	if h.focus == FocusActive {
		cKeyString = "Complete Task"
	}

	refreshingLabel := ""
	if h.activeLoading || h.completedLoading {
		refreshingLabel = h.loadingSpinner.View() + " Refreshing..."
	}

	tabsView := h.tabs.View(width)

	content := fmt.Sprintf(
		"%s\n\n%s\n\nCompleted:\n%s\n\n%s",
		refreshingLabel,
		h.activeTaskTable.View(),
		h.completedTaskTable.View(),
		fmt.Sprintf("Controls: [Tab] Focus • [r] Refresh • [n] New Task • [x] Delete Task • [c] %s", cKeyString),
	)

	windowedContent := h.tabs.WrapContent(content, width)

	finalView := lipgloss.JoinVertical(lipgloss.Left, tabsView, windowedContent)

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		MaxHeight(height).
		MaxWidth(width).
		Render(finalView)
}
