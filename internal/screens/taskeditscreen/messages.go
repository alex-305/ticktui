package taskeditscreen

import types "github.com/alex-305/ticktui/pkg/tickticktypes"

type taskCreatedMsg struct {
	task *types.Task
	err  error
}

type taskUpdatedMsg struct {
	task *types.Task
	err  error
}
