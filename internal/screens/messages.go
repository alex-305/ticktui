package screens

import "github.com/alex-305/ticktui/internal/types"

type ChangeScreenMsg struct {
	NewScreen Screen
}

type GoBackScreenMsg struct{}

type TaskDeletedMsg struct{ err error }

type ActiveTaskListMsg struct {
	tasks []types.Task
	err   error
}

type CompletedTaskListMsg struct {
	tasks []types.Task
	err   error
}

type ProjectsLoadedMsg struct {
	projects []types.Project
	err      error
}
