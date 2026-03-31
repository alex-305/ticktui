package homescreen

import (
	"time"

	types "github.com/alex-305/ticktui/pkg/tickticktypes"
	tea "github.com/charmbracelet/bubbletea"
)

func (h *HomeScreen) deleteTaskCmd(task *types.Task) tea.Cmd {
	return func() tea.Msg {
		err := h.ctx.APIClient.DeleteTask(task.ProjectID, task.ID)
		return ActionCompletedMsg{err}
	}
}

func (h *HomeScreen) completeTaskCmd(task *types.Task) tea.Cmd {
	return func() tea.Msg {
		err := h.ctx.APIClient.CompleteTask(task)
		return ActionCompletedMsg{err}
	}
}

func (h *HomeScreen) decompleteTaskCmd(task *types.Task) tea.Cmd {
	return func() tea.Msg {
		err := h.ctx.APIClient.DecompleteTask(task)
		return ActionCompletedMsg{err}
	}
}

func (h HomeScreen) fetchCompletedTasksCmd(projectIDs []string) tea.Cmd {
	return func() tea.Msg {
		now := time.Now()
		tasks, err := h.ctx.APIClient.ListCompletedTasks(projectIDs, types.TickTickTime(now.AddDate(0, 0, -4000)), types.TickTickTime(now))
		if err != nil {
			return CompletedTaskListMsg{tasks: tasks, err: err}
		}
		return CompletedTaskListMsg{tasks: tasks, err: nil}
	}
}

func (h *HomeScreen) fetchActiveTasksCmd(projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := h.ctx.APIClient.ListTasks(projectID)
		if err != nil {
			return ActiveTaskListMsg{tasks: tasks, err: err}
		}
		return ActiveTaskListMsg{tasks: tasks, err: nil}
	}
}

func (h *HomeScreen) fetchProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		projects, err := h.ctx.APIClient.ListProjects()
		if err != nil {
			return ProjectsLoadedMsg{projects: projects, err: err}
		}
		return ProjectsLoadedMsg{projects: projects, err: nil}
	}
}
