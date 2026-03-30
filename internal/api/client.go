package api

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/config"
	"github.com/alex-305/ticktui/internal/types"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.ticktick.com/open/v1"
)

type Client struct {
	http *resty.Client
}

func NewClient(token string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Authorization", "Bearer "+token)

	return &Client{http: client}
}

func GetClient() (*Client, error) {
	token, err := config.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	return NewClient(token), nil
}

func (c *Client) ListProjects() ([]types.Project, error) {
	var projects []types.Project
	resp, err := c.http.R().
		SetResult(&projects).
		Get("/project")

	if err != nil {
		return nil, errors.Wrap(err, "listing projects")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list projects: %s", resp.String())
	}

	projects = append(projects, types.InboxProject)

	return projects, nil
}

func (c *Client) GetProject(id string) (types.Project, error) {
	if id == types.InboxProject.ID {
		return types.InboxProject, nil
	}
	var project types.Project
	resp, err := c.http.R().
		SetResult(&project).
		Get("/project/" + id)

	if err != nil {
		return types.NullProject, errors.Wrap(err, "getting project")
	}
	if resp.IsError() {
		return types.NullProject, fmt.Errorf("failed to get project: %s", resp.String())
	}
	if project == types.NullProject {
		return types.NullProject, fmt.Errorf("project not found: %s", id)
	}

	return project, nil
}

func (c *Client) GetTask(projectID string, taskID string) (*types.Task, error) {
	var task types.Task
	resp, err := c.http.R().
		SetResult(&task).
		Get(fmt.Sprintf("/project/%s/task/%s", projectID, taskID))

	if err != nil {
		return nil, errors.Wrap(err, "requesting task")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to list tasks: %s", resp.String())
	}

	return &task, nil
}

func (c *Client) ListTasks(projectID string) ([]types.Task, error) {
	var projectData struct {
		Tasks []types.Task `json:"tasks"`
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

func (c *Client) ListCompletedTasks(projectIDs []string, startDate, endDate types.TickTickTime) ([]types.Task, error) {
	var tasks []types.Task

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

func (c *Client) GetProjectWithTasks(projectID string) (*types.ProjectData, error) {
	var projectData types.ProjectData
	resp, err := c.http.R().
		SetResult(&projectData).
		Get(fmt.Sprintf("/project/%s/data", projectID))

	if err != nil {
		return nil, errors.Wrap(err, "getting project data")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to get project data: %s", resp.String())
	}

	return &projectData, nil
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

func (c *Client) UpdateProject(project types.Project) (types.Project, error) {
	resp, err := c.http.R().
		SetBody(project).
		SetResult(project).
		Post(fmt.Sprintf("/project/%s", project.ID))

	if err != nil {
		return types.NullProject, errors.Wrap(err, "updating project")
	}
	if resp.IsError() {
		return types.NullProject, fmt.Errorf("failed to update project: %s", resp.String())
	}

	return project, nil
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

func (c *Client) CompleteTask(projectID, taskID string) error {
	resp, err := c.http.R().
		Post(fmt.Sprintf("/project/%s/task/%s/complete", projectID, taskID))

	if err != nil {
		return errors.Wrap(err, "completing task")
	}
	if resp.IsError() {
		return fmt.Errorf("failed to complete task: %s", resp.String())
	}

	return nil
}

func (c *Client) CreateProject(project *types.Project) (*types.Project, error) {
	if project == nil {
		return nil, errors.New("project cannot be nil")
	}

	resp, err := c.http.R().
		SetBody(project).
		SetResult(project).
		Post("/project")

	if err != nil {
		return nil, errors.Wrap(err, "creating project")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to create project: %s", resp.String())
	}

	return project, nil
}

func (c *Client) DeleteProject(projectID string) error {
	resp, err := c.http.R().
		Delete(fmt.Sprintf("/project/%s", projectID))

	if err != nil {
		return errors.Wrap(err, "deleting project")
	}
	if resp.IsError() {
		return fmt.Errorf("failed to delete project: %s", resp.String())
	}

	return nil
}
