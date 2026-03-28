package screens

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/types"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

type HomeScreen struct {
	ctx           context.AppContext
	projects      []types.Project
	activeProject int

	taskTable components.TaskTable
	paginator paginator.Model

	loaded  bool
	loading bool
	err     error
}

func NewHomeScreen(ctx context.AppContext) Screen {
	return HomeScreen{
		ctx:     ctx,
		loaded:  false,
		loading: false,
	}
}
func (h HomeScreen) deleteTaskCmd(task types.Task) tea.Cmd {
	return func() tea.Msg {
		err := h.ctx.TaskService.DeleteTask(task.ProjectID, task.ID)
		return TaskDeletedMsg{err}
	}
}

func (h HomeScreen) fetchTasksCmd(projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := h.ctx.TaskService.ListTasks(projectID)
		if err != nil {
			return err
		}
		return tasks
	}
}

func (h HomeScreen) fetchProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		projects, err := h.ctx.TaskService.ListProjects()
		if err != nil {
			return err
		}
		return projects
	}
}

func (h HomeScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	var cmds []tea.Cmd

	if !h.loaded && !h.loading && h.err == nil {
		h.loading = true
		return h, h.fetchProjectsCmd()
	}

	switch msg := msg.(type) {
	case error:
		h.loading = false
		h.err = msg
		return h, nil

	case []types.Project:
		h.projects = msg
		p := paginator.New()
		p.Type = paginator.Dots
		p.SetTotalPages(len(msg))
		h.paginator = p

		return h, h.fetchTasksCmd(h.projects[h.activeProject].ID)

	case GoBackScreenMsg:
		h.loading = true

		return h, h.fetchTasksCmd(h.projects[h.activeProject].ID)
	case TaskDeletedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.loading = false
			return h, nil
		}
		h.loading = true
		return h, h.fetchTasksCmd(h.projects[h.activeProject].ID)

	case []types.Task:
		h.taskTable = components.NewTaskTable(msg, width)
		h.loaded = true
		h.loading = false
		return h, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if h.activeProject < len(h.projects)-1 {
				h.activeProject++
				h.paginator.Page++
				h.loading = true
				return h, h.fetchTasksCmd(h.projects[h.activeProject].ID)
			}
		case "shift+tab":
			if h.activeProject > 0 {
				h.activeProject--
				h.paginator.Page--
				h.loading = true
				return h, h.fetchTasksCmd(h.projects[h.activeProject].ID)
			}
		case "r":
			h.loading = true
			return h, h.fetchTasksCmd(h.projects[h.activeProject].ID)
		case "n":
			return h, func() tea.Msg {
				return ChangeScreenMsg{NewScreen: NewCreateTaskScreen(h.ctx, h.projects[h.activeProject].ID)}
			}
		case "x":
			selectedTask, ok := h.taskTable.GetSelectedTask()
			if !ok {
				return h, nil
			}

			return h, h.deleteTaskCmd(selectedTask)
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
		"My Tasks in %s\n\n%s\n\n%s\n\n[Tab/Shift+Tab] Switch Projects • [r] Refresh • [n] New Task • [x] Delete Selected Task • [ctrl + c] Quit",
		h.projects[h.activeProject].Name,
		h.taskTable.View(),
		h.paginator.View(),
	)
}
