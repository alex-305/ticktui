package createtaskscreen

import (
	"fmt"
	"regexp"

	"github.com/alex-305/ticktui/internal/context"
	"github.com/alex-305/ticktui/internal/screens"
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
	dueDate  string
	priority task.Priority
}

func NewCreateTaskScreen(ctx context.AppContext, projectID string) screens.Screen {
	ct := &CreateTaskScreen{ctx: ctx, projectID: projectID}

	ct.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Task Title").
				Value(&ct.title).
				Placeholder("What needs to be done?").
				Key("title"),

			huh.NewText().
				Title("Description").
				Value(&ct.desc).
				Placeholder("Add details...").
				Lines(5),

			huh.NewInput().
				Title("Due Date").
				Value(&ct.dueDate).
				Placeholder("YYYY-MM-DD (Optional)").
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					_, err := types.StringToTickTickTime(s)
					return err
				}),

			huh.NewSelect[task.Priority]().
				Title("Priority").
				Options(
					huh.NewOption("None", task.PriorityNone),
					huh.NewOption("Low", task.PriorityLow),
					huh.NewOption("Medium", task.PriorityMedium),
					huh.NewOption("High", task.PriorityHigh),
				).
				Value(&ct.priority),
		),
	)

	ct.form.Init()

	return ct
}

func (ct *CreateTaskScreen) Init() tea.Cmd {
	return nil
}

func (ct *CreateTaskScreen) Update(msg tea.Msg, width, height int) (screens.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case taskCreatedMsg:
		if msg.err != nil {
			ct.err = msg.err
			ct.loading = false
			return ct, nil
		}
		return ct, func() tea.Msg {
			return screens.GoBackScreenMsg{}
		}
	}

	form, cmd := ct.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		ct.form = f
	}

	if ct.form.State == huh.StateCompleted && !ct.loading {
		ct.loading = true
		return ct, func() tea.Msg {

			dueDate, err := types.StringToTickTickTime(ct.dueDate)
			if err != nil {
				return taskCreatedMsg{task: nil, err: err}
			}
			task, err := ct.ctx.APIClient.CreateTask(&types.Task{
				Title:     ct.title,
				Desc:      ct.desc,
				DueDate:   dueDate,
				Priority:  ct.priority,
				ProjectID: ct.projectID,
			})
			return taskCreatedMsg{task: task, err: err}
		}
	}

	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEsc {
		return ct, func() tea.Msg {
			return screens.GoBackScreenMsg{}
		}
	}

	return ct, cmd
}

func (ct *CreateTaskScreen) View(width, height int) string {
	if ct.loading {
		return "\n  ⏳ Creating task..."
	}

	var errMsg string
	if ct.err != nil {
		re := regexp.MustCompile(`"errorMessage"\s*:\s*"([^"]*)"`)
		matches := re.FindStringSubmatch(ct.err.Error())

		errMsg = "❌ Error: "

		if len(matches) > 0 {
			errMsg = errMsg + matches[1]
		} else {
			errMsg = errMsg + "Unknown"
		}
	}

	return fmt.Sprintf(
		" Create New Task\n\n%s\n%s\n [Esc] Cancel",
		ct.form.View(),
		errMsg,
	)
}
