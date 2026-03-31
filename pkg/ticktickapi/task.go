package ticktickapi

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/types"
	"github.com/alex-305/ticktui/internal/types/task"
	"github.com/pkg/errors"
)

func (c *Client) GetTask(projectID string, taskID string) (*types.Task, error) {
	var task *types.Task
	resp, err := c.http.R().
		SetResult(&task).
		Get(fmt.Sprintf("/project/%s/task/%s", projectID, taskID))

	if err != nil {
		return nil, errors.Wrap(err, "requesting task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to list tasks: %s", resp.String())
	}

	return task, nil
}

func (c *Client) ListTasks(projectID string) ([]*types.Task, error) {
	var projectData struct {
		Tasks []*types.Task `json:"tasks"`
	}
	resp, err := c.http.R().
		SetResult(&projectData).
		Get(fmt.Sprintf("/project/%s/data", projectID))

	if err != nil {
		return nil, errors.Wrap(err, "listing tasks")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list tasks: %s", resp.String())
	}

	return projectData.Tasks, nil
}

func (c *Client) ListCompletedTasks(projectIDs []string, startDate, endDate types.TickTickTime) ([]*types.Task, error) {
	var tasks []*types.Task

	stDate := startDate.ToMSFormat()
	eDate := endDate.ToMSFormat()

	resp, err := c.http.R().
		SetResult(&tasks).
		SetBody(map[string]any{
			"startDate": stDate,
			"endDate":   eDate},
		).
		Post("/task/completed")

	if err != nil {
		return nil, errors.Wrap(err, "listing tasks")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list tasks: %s", resp.String())
	}

	return tasks, nil
}

func (c *Client) CreateTask(task *types.Task) (*types.Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	resp, err := c.http.R().
		SetBody(task).
		SetResult(task).
		Post("/task")

	if err != nil {
		return nil, errors.Wrap(err, "creating task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to create task: %s", resp.String())
	}

	return task, nil
}

func (c *Client) UpdateTask(task *types.Task) (*types.Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	resp, err := c.http.R().
		SetBody(task).
		SetResult(task).
		Post(fmt.Sprintf("/task/%s", task.ID))

	if err != nil {
		return nil, errors.Wrap(err, "updating task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to update task: %s", resp.String())
	}

	return task, nil
}

func (c *Client) DeleteTask(projectID, taskID string) error {
	resp, err := c.http.R().
		Delete(fmt.Sprintf("/project/%s/task/%s", projectID, taskID))

	if err != nil {
		return errors.Wrap(err, "deleting task")
	}
	if resp.IsError() {
		return fmt.Errorf("failed to delete task: %s", resp.String())
	}

	return nil
}

func (c *Client) CompleteTask(t *types.Task) error {
	resp, err := c.http.R().
		Post(fmt.Sprintf("/project/%s/task/%s/complete", t.ProjectID, t.ID))

	if err != nil {
		return errors.Wrap(err, "completing task")
	}
	if resp.IsError() {
		return fmt.Errorf("failed to complete task: %s", resp.String())
	}

	return nil
}

func (c *Client) DecompleteTask(t *types.Task) error {
	t.Status = task.StatusNormal
	t.CompletedTime = types.TickTickTime{}
	_, err := c.UpdateTask(t)

	if err != nil {
		return err
	}

	return nil
}
