package api

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/types"
	"github.com/pkg/errors"
)

func (c *Client) ListProjects() ([]*types.Project, error) {
	var projects []*types.Project
	resp, err := c.http.R().
		SetResult(&projects).
		Get("/project")

	if err != nil {
		return nil, errors.Wrap(err, "listing projects")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list projects: %s", resp.String())
	}

	projects = append(projects, &types.InboxProject)

	return projects, nil
}

func (c *Client) GetProject(id string) (*types.Project, error) {
	if id == types.InboxProject.ID {
		return &types.InboxProject, nil
	}
	var project *types.Project
	resp, err := c.http.R().
		SetResult(&project).
		Get("/project/" + id)

	if err != nil {
		return &types.NullProject, errors.Wrap(err, "getting project")
	}
	if resp.IsError() {
		return &types.NullProject, fmt.Errorf("failed to get project: %s", resp.String())
	}
	if project == &types.NullProject {
		return &types.NullProject, fmt.Errorf("project not found: %s", id)
	}

	return project, nil
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

func (c *Client) UpdateProject(project *types.Project) (*types.Project, error) {
	resp, err := c.http.R().
		SetBody(project).
		SetResult(project).
		Post(fmt.Sprintf("/project/%s", project.ID))

	if err != nil {
		return &types.NullProject, errors.Wrap(err, "updating project")
	}
	if resp.IsError() {
		return &types.NullProject, fmt.Errorf("failed to update project: %s", resp.String())
	}

	return project, nil
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
