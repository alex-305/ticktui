package services

import "github.com/alex-305/ticktui/internal/api"
import "github.com/alex-305/ticktui/internal/types"

type TaskService struct {
	client *api.Client
}

func NewTaskService(client *api.Client) *TaskService {
	return &TaskService{client: client}
}

func (s *TaskService) CreateTask(title, description string) (*types.Task, error) {
	task := types.Task{
		Title: title,
		Desc:  description,
	}

	return s.client.CreateTask(&task)
}
