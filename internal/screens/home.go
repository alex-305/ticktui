package screens

import (
	"fmt"
	"time"

	"github.com/alex-305/ticktui/internal/components"
	"github.com/alex-305/ticktui/internal/context"
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

	taskTable          components.TaskTable
	completedTaskTable components.TaskTable
	paginator          paginator.Model
	focus              Focus

	activeLoaded     bool
	activeLoading    bool
	completedLoaded  bool
	completedLoading bool
	err              error
}

func NewHomeScreen(ctx context.AppContext) Screen {
	return &HomeScreen{
		ctx:              ctx,
		activeLoaded:     false,
		activeLoading:    false,
		completedLoaded:  false,
		completedLoading: false,
	}
}
func (h *HomeScreen) deleteTaskCmd(task types.Task) tea.Cmd {
	return func() tea.Msg {
		err := h.ctx.APIClient.DeleteTask(task.ProjectID, task.ID)
		return TaskDeletedMsg{err}
	}
}

func (h HomeScreen) fetchCompletedTasksCmd(projectIDs []string) tea.Cmd {
	return func() tea.Msg {
		now := time.Now()
		tasks, err := h.ctx.APIClient.ListCompletedTasks(projectIDs, types.TickTickTime(now.AddDate(0, 0, -4000)), types.TickTickTime(now))
		if err != nil {
			return err
		}
		return CompletedTaskListMsg(tasks)
	}
}

func (h *HomeScreen) fetchActiveTasksCmd(projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := h.ctx.APIClient.ListTasks(projectID)
		if err != nil {
			return err
		}
		return ActiveTaskListMsg(tasks)
	}
}

func (h *HomeScreen) fetchAllData(projectID string, allProjectIDs []string) tea.Cmd {
	return tea.Batch(
		h.fetchCompletedTasksCmd(allProjectIDs),
		h.fetchActiveTasksCmd(projectID),
	)
}

func (h *HomeScreen) fetchProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		projects, err := h.ctx.APIClient.ListProjects()
		if err != nil {
			return err
		}
		return projects
	}
}

func (h *HomeScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	var cmds []tea.Cmd

	if !h.activeLoaded && !h.activeLoading && h.err == nil {
		h.activeLoading = true
		return h, h.fetchProjectsCmd()
	}

	switch msg := msg.(type) {

	case ActiveTaskListMsg:
		h.taskTable = components.NewTaskTable(msg, width)
		h.activeLoading = false
		h.activeLoaded = true
		return h, nil
	case CompletedTaskListMsg:
		h.completedTaskTable = components.NewTaskTable(msg, width)
		h.completedLoading = false
		h.completedLoaded = true
		return h, nil

	case error:
		h.activeLoading = false
		h.completedLoading = false
		h.err = msg
		return h, nil

	case []types.Project:
		h.projects = msg
		ids := make([]string, len(msg))
		for i, p := range msg {
			ids[i] = p.ID
		}
		h.projectIDs = ids
		p := paginator.New()
		p.SetTotalPages(len(msg))
		h.paginator = p

		return h, h.fetchAllData(h.projects[h.activeProject].ID, ids)

	case GoBackScreenMsg:
		h.activeLoading = true
		h.completedLoading = true

		return h, h.fetchAllData(h.projects[h.activeProject].ID, h.projectIDs)
	case TaskDeletedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.activeLoading = false
			return h, nil
		}
		h.activeLoading = true
		return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if h.focus == FocusActive {
				h.focus = FocusCompleted
			} else {
				h.focus = FocusActive
			}
			return h, nil
		case "l":
			if h.activeProject < len(h.projects)-1 {
				h.activeProject++
				h.paginator.Page++
				h.activeLoading = true
				return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID)
			}
		case "h":
			if h.activeProject > 0 {
				h.activeProject--
				h.paginator.Page--
				h.activeLoading = true
				return h, h.fetchActiveTasksCmd(h.projects[h.activeProject].ID)
			}
		case "r":
			h.activeLoading = true
			h.completedLoading = true
			return h, h.fetchAllData(h.projects[h.activeProject].ID, h.projectIDs)
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

	if h.activeLoaded {
		cmd := h.taskTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	if h.completedLoaded {
		cmd := h.completedTaskTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return h, tea.Batch(cmds...)
}

func (h *HomeScreen) View(width, height int) string {
	if h.err != nil {
		return components.NewErrorBox(h.err, width, height).View()
	}

	if len(h.projects) == 0 {
		return "\n  Initializing projects..."
	}

	return fmt.Sprintf(
		"Project: %s\n\n%s\n\n%s\n\nCompleted:\n%s\n\n%s",
		h.projects[h.activeProject].Name,
		h.taskTable.View(),
		h.paginator.View(),
		h.completedTaskTable.View(),
		"Controls: [Tab] Focus • [r] Refresh",
	)
}
