package components

import (
	"errors"
	"time"

	types "github.com/alex-305/ticktui/pkg/tickticktypes"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskTable struct {
	Model table.Model
	Tasks []*types.Task
}

func NewTaskTable(tasks []*types.Task, width, height int) TaskTable {
	usableWidth := max(width, 20)
	columns := []table.Column{
		{Title: "Title", Width: int(float64(usableWidth) * 0.3)},
		{Title: "Description", Width: int(float64(usableWidth) * 0.4)},
		{Title: "Due Date", Width: int(float64(usableWidth) * 0.15)},
		{Title: "Priority", Width: int(float64(usableWidth) * 0.1)},
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

	usableHeight := max(height, 15)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(usableHeight),
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

	inactiveGray := lipgloss.Color("240")

	s.Cell = s.Cell.
		Foreground(inactiveGray).
		Background(lipgloss.NoColor{}).
		Bold(false)

	s.Selected = s.Cell

	s.Header = s.Header.
		Foreground(inactiveGray).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		BorderBottom(true).
		Bold(true)

	tt.Model.SetStyles(s)
}

func (tt *TaskTable) SetDimensions(width, height int) error {
	tt.Model.SetWidth(width)
	tt.Model.SetHeight(height)
	c := tt.Model.Columns()

	if len(c) == 0 {
		return errors.New("columns not yet initialized")
	}

	c[0].Width = int(float64(width) * 0.3)
	c[1].Width = int(float64(width) * 0.4)
	c[2].Width = int(float64(width) * 0.15)
	c[3].Width = int(float64(width) * 0.1)
	tt.Model.SetColumns(c)

	return nil
}

func (tt *TaskTable) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	tt.Model, cmd = tt.Model.Update(msg)
	return cmd
}

func (tt *TaskTable) GetSelectedTask() (*types.Task, bool) {
	if len(tt.Tasks) == 0 {
		return nil, false
	}

	currIndex := tt.Model.Cursor()
	if currIndex < 0 || currIndex >= len(tt.Tasks) {
		return nil, false
	}

	return tt.Tasks[currIndex], true
}

func (tt TaskTable) View() string {
	return tt.Model.View()
}
