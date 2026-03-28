package screens

import (
	"fmt"
	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

type HomeScreen struct {
	ctx       context.AppContext
	taskTable components.TaskTable
	loaded    bool
	loading   bool
	err       error
}

func NewHomeScreen(ctx context.AppContext) Screen {
	return HomeScreen{
		ctx:     ctx,
		loaded:  false,
		loading: false,
	}
}

func (h HomeScreen) fetchTasksCmd() tea.Cmd {
	return func() tea.Msg {
		tasks, err := h.ctx.TaskService.ListTasks("inbox")
		if err != nil {
			return err
		}
		return tasks
	}
}

func (h HomeScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	var cmds []tea.Cmd

	if !h.loaded && !h.loading && h.err == nil {
		h.loading = true
		return h, h.fetchTasksCmd()
	}

	switch msg := msg.(type) {
	case error:
		h.loading = false
		h.err = msg
		return h, nil

	case GoBackScreenMsg:
		h.loading = true
		return h, h.fetchTasksCmd()

	case []types.Task:
		h.taskTable = components.NewTaskTable(msg, width)
		h.loaded = true
		h.loading = false
		return h, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			h.loading = true
			return h, h.fetchTasksCmd()
		case "n":
			return h, func() tea.Msg {
				return ChangeScreenMsg{NewScreen: NewCreateTaskScreen(h.ctx)}
			}
		case "ctrl+c":
			return h, tea.Quit
		}
	}

	if h.loaded {
		cmd := h.taskTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return h, tea.Batch(cmds...)
}

func (h HomeScreen) View(width, height int) string {
	if h.err != nil {
		return fmt.Sprintf("\n Error: %v\n\n [q] Quit", h.err)
	}

	if !h.loaded {
		return "\n  Loading tasks..."
	}

	return fmt.Sprintf(
		"TickTUI - My Tasks\n\n%s\n\n[r] Refresh [n] New Task • [ctrl + c] Quit",
		h.taskTable.View(),
	)
}
