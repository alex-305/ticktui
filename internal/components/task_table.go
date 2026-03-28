package components

import (
	"time"

	"github.com/alex-305/ticktui/internal/types"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskTable struct {
	Model table.Model
	Tasks []types.Task
}

func NewTaskTable(tasks []types.Task, width int) TaskTable {
	columns := []table.Column{
		{Title: "Title", Width: 25},
		{Title: "Description", Width: 35},
		{Title: "Due Date", Width: 15},
		{Title: "Priority", Width: 10},
	}

	rows := make([]table.Row, len(tasks))
	for i, t := range tasks {
		tm := time.Time(t.DueDate)
		dueDateStr := "None"
		if !tm.IsZero() {
			dueDateStr = tm.Format("2006-01-02")
		}

		rows[i] = table.Row{
			t.Title,
			t.Desc,
			dueDateStr,
			renderPriority(int(t.Priority)),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	tt := TaskTable{Model: t, Tasks: tasks}
	tt.ApplyActiveStyle()

	return tt
}

func renderPriority(p int) string {
	switch p {
	case 5:
		return "!!!"
	case 3:
		return "!!"
	case 1:
		return "!"
	default:
		return "None"
	}
}

func (tt *TaskTable) ApplyActiveStyle() {
	s := table.DefaultStyles()

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	s.Header = s.Header.
		Foreground(lipgloss.Color("205")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)

	tt.Model.SetStyles(s)
}

func (tt *TaskTable) ApplyInactiveStyle() {
	s := table.DefaultStyles()

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235")).
		Bold(false)

	s.Header = s.Header.
		Foreground(lipgloss.Color("240")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		BorderBottom(true).
		Bold(false)

	tt.Model.SetStyles(s)
}

func (tt *TaskTable) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	tt.Model, cmd = tt.Model.Update(msg)
	return cmd
}

func (tt *TaskTable) GetSelectedTask() (types.Task, bool) {

	if len(tt.Tasks) == 0 {
		return types.Task{}, false
	}

	currIndex := tt.Model.Cursor()

	if currIndex < 0 || currIndex >= len(tt.Tasks) {
		return types.Task{}, false
	}

	return tt.Tasks[currIndex], true

}

func (tt TaskTable) View() string {
	return tt.Model.View()
}
