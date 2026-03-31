package taskeditscreen

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

type TaskEditScreen struct {
	form    *huh.Form
	ctx     context.AppContext
	err     error
	loading bool

	dueDateString string
	newTask       bool

	task *types.Task
}

func NewTaskEditScreen(ctx context.AppContext, projectID string, taskToEdit *types.Task) screens.Screen {

	tf := &TaskEditScreen{
		ctx:  ctx,
		task: taskToEdit,
	}

	if tf.task == nil {
		tf.newTask = true
		tf.task = &types.Task{}
	}

	tf.task.ProjectID = projectID

	tf.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Task Title").
				Value(&tf.task.Title).
				Placeholder("What needs to be done?").
				Key("title"),

			huh.NewText().
				Title("Description").
				Value(&tf.task.Desc).
				Placeholder("Add details...").
				Lines(5),

			huh.NewInput().
				Title("Due Date").
				Value(&tf.dueDateString).
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
				Value(&tf.task.Priority),
		),
	)

	tf.form.Init()

	return tf
}

func (tf *TaskEditScreen) Init() tea.Cmd {
	return nil
}

func (tf *TaskEditScreen) Update(msg tea.Msg, width, height int) (screens.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case taskCreatedMsg:
		if msg.err != nil {
			tf.err = msg.err
			tf.loading = false
			return tf, nil
		}
		return tf, func() tea.Msg {
			return screens.GoBackScreenMsg{}
		}
	}

	form, cmd := tf.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		tf.form = f
	}

	if tf.form.State == huh.StateCompleted && !tf.loading {
		tf.loading = true
		return tf, func() tea.Msg {

			if tf.dueDateString != "" {
				dueDate, err := types.StringToTickTickTime(tf.dueDateString)
				if err != nil {
					return taskCreatedMsg{task: nil, err: err}
				}
				tf.task.DueDate = dueDate
			}
			task, err := tf.ctx.APIClient.CreateTask(tf.task)
			return taskCreatedMsg{task: task, err: err}
		}
	}

	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEsc {
		return tf, func() tea.Msg {
			return screens.GoBackScreenMsg{}
		}
	}

	return tf, cmd
}

func (tf *TaskEditScreen) View(width, height int) string {
	if tf.loading {
		if tf.newTask {
			return "\n Creating task..."
		} else {
			return "\n Updating task..."
		}
	}

	var errMsg string
	if tf.err != nil {
		re := regexp.MustCompile(`"errorMessage"\s*:\s*"([^"]*)"`)
		matches := re.FindStringSubmatch(tf.err.Error())

		errMsg = "❌ Error: "

		if len(matches) > 0 {
			errMsg = errMsg + matches[1]
		} else {
			errMsg = errMsg + "Unknown"
		}
	}

	var headerString string

	if tf.newTask {
		headerString = "Create New Task"
	} else {
		headerString = "Update Task"
	}

	return fmt.Sprintf(
		" %s\n\n%s\n%s\n [Esc] Cancel",
		headerString,
		tf.form.View(),
		errMsg,
	)
}
