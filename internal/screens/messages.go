package screens

import "github.com/alex-305/ticktui/internal/types"

type ChangeScreenMsg struct {
	NewScreen Screen
}

type GoBackScreenMsg struct{}

type TaskDeletedMsg struct{ err error }

type ActiveTaskListMsg []types.Task
type CompletedTaskListMsg []types.Task
