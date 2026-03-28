package services

import (
	"github.com/alex-305/ticktui/internal/api"
	"github.com/alex-305/ticktui/internal/types"
	"github.com/alex-305/ticktui/internal/types/task"
)

type TaskService struct {
	client *api.Client
}

func NewTaskService(client *api.Client) *TaskService {
	return &TaskService{client: client}
}

func (s *TaskService) CreateTask(title, description string, priority task.Priority) (*types.Task, error) {
	task := types.Task{
		Title:    title,
		Desc:     description,
		Priority: priority,
	}

	return s.client.CreateTask(&task)
}

func (s *TaskService) DeleteTask(projectID, taskID string) error {
	return s.client.DeleteTask(projectID, taskID)

}

func (s *TaskService) ListTasks(projectID string) ([]types.Task, error) {
	return s.client.ListTasks(projectID)
}

func (s *TaskService) ListProjects() ([]types.Project, error) {
	return s.client.ListProjects()
}
