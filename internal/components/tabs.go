package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TabStyles struct {
	InactiveTab lipgloss.Style
	ActiveTab   lipgloss.Style
	Window      lipgloss.Style
}

func DefaultTabStyles() TabStyles {
	highlightColor := lipgloss.Color("#874BFD")
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")

	s := TabStyles{}
	s.InactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)

	s.ActiveTab = s.InactiveTab.
		Border(activeTabBorder, true)

	s.Window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(1, 1).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()

	return s
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

type Tabs struct {
	Items  []string
	Active int
	Styles TabStyles
}

func NewTabs() Tabs {
	return Tabs{
		Styles: DefaultTabStyles(),
	}
}

func (t *Tabs) SetItems(items []string) {
	t.Items = items
}

func (t *Tabs) SetActive(idx int) {
	if idx >= 0 && idx < len(t.Items) {
		t.Active = idx
	}
}

func (t Tabs) GetWindowWidth(screenWidth int) int {
	return screenWidth - t.Styles.Window.GetHorizontalFrameSize()
}

func (t Tabs) GetWindowHeight(screenHeight int) int {
	tabRowHeight := 3
	windowFrame := t.Styles.Window.GetVerticalFrameSize()
	return screenHeight - tabRowHeight - windowFrame
}

func (t Tabs) View(width int) string {
	if len(t.Items) == 0 {
		return ""
	}

	var renderedTabs []string

	for i, item := range t.Items {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(t.Items)-1, i == t.Active

		if isActive {
			style = t.Styles.ActiveTab
		} else {
			style = t.Styles.InactiveTab
		}

		border := style.GetBorderStyle()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(item))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	gap := max(width-lipgloss.Width(tabRow)-1, 0)

	filler := strings.Repeat("─", gap) + "┐"
	fillerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#874BFD"))

	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabRow, fillerStyle.Render(filler))
}

func (t Tabs) WrapContent(content string, width int) string {
	return t.Styles.Window.Width(width - 2).Render(content)
}
