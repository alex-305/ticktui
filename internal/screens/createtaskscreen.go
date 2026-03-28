package screens

import (
	"fmt"
	"regexp"

	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/types"
	"github.com/alex-305/ticktui/internal/types/task"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type taskCreatedMsg struct {
	task *types.Task
	err  error
}

type CreateTaskScreen struct {
	form    *huh.Form
	ctx     context.AppContext
	err     error
	loading bool

	projectID string

	title    string
	desc     string
	priority task.Priority
}

func NewCreateTaskScreen(ctx context.AppContext, projectID string) Screen {
	s := &CreateTaskScreen{ctx: ctx, projectID: projectID}

	s.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Task Title").
				Value(&s.title).
				Placeholder("What needs to be done?").
				Key("title"),

			huh.NewText().
				Title("Description").
				Value(&s.desc).
				Placeholder("Add details...").
				Lines(5),

			huh.NewSelect[task.Priority]().
				Title("Priority").
				Options(
					huh.NewOption("None", task.PriorityNone),
					huh.NewOption("Low", task.PriorityLow),
					huh.NewOption("Medium", task.PriorityMedium),
					huh.NewOption("High", task.PriorityHigh),
				).
				Value(&s.priority),
		),
	)

	s.form.Init()

	return s
}

func (h *CreateTaskScreen) Update(msg tea.Msg, width, height int) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case taskCreatedMsg:
		if msg.err != nil {
			h.err = msg.err
			h.loading = false
			return h, nil
		}
		return h, func() tea.Msg {
			return GoBackScreenMsg{}
		}
	}

	form, cmd := h.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		h.form = f
	}

	if h.form.State == huh.StateCompleted && !h.loading {
		h.loading = true
		return h, func() tea.Msg {

			task, err := h.ctx.APIClient.CreateTask(&types.Task{
				Title:     h.title,
				Desc:      h.desc,
				Priority:  h.priority,
				ProjectID: h.projectID,
			})
			return taskCreatedMsg{task: task, err: err}
		}
	}

	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEsc {
		return h, func() tea.Msg {
			return ChangeScreenMsg{NewScreen: NewHomeScreen(h.ctx)}
		}
	}

	return h, cmd
}

func (h *CreateTaskScreen) View(width, height int) string {
	if h.loading {
		return "\n  ⏳ Creating task..."
	}

	var errMsg string
	if h.err != nil {
		re := regexp.MustCompile(`"errorMessage"\s*:\s*"([^"]*)"`)
		matches := re.FindStringSubmatch(h.err.Error())

		errMsg = "❌ Error: "

		if len(matches) > 0 {
			errMsg = errMsg + matches[1]
		} else {
			errMsg = errMsg + "Unknown"
		}
	}

	return fmt.Sprintf(
		" Create New Task\n\n%s\n%s\n [Esc] Cancel",
		h.form.View(),
		errMsg,
	)
}
