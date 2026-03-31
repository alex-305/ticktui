package taskeditscreen

import "github.com/alex-305/ticktui/internal/types"

type taskCreatedMsg struct {
	task *types.Task
	err  error
}

type taskUpdatedMsg struct {
	task *types.Task
	err  error
}
