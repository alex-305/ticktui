package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TabStyles struct {
	InactiveTab lipgloss.Style
	ActiveTab   lipgloss.Style
	HintTab     lipgloss.Style
	Window      lipgloss.Style
}

func DefaultTabStyles() TabStyles {
	highlightColor := lipgloss.Color("#874BFD")
	hintColor := lipgloss.Color("#555555")

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")

	s := TabStyles{}
	s.InactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)

	s.ActiveTab = s.InactiveTab.
		Border(activeTabBorder, true)

	s.HintTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(hintColor).
		Foreground(hintColor).
		Padding(0, 1)

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
	LeftHint  string
	RightHint string
	Items     []string
	Active    int
	Styles    TabStyles
}

func NewTabs(leftHint, rightHint string) Tabs {
	return Tabs{
		LeftHint:  leftHint,
		RightHint: rightHint,
		Styles:    DefaultTabStyles(),
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
	if len(t.Items) == 0 && t.LeftHint == "" && t.RightHint == "" {
		return ""
	}

	type tabNode struct {
		text     string
		style    lipgloss.Style
		isActive bool
	}

	var nodes []tabNode
	if len(t.Items) > 1 && t.LeftHint != "" {
		nodes = append(nodes, tabNode{
			text:     "← " + t.LeftHint,
			style:    t.Styles.HintTab,
			isActive: false,
		})
	}

	for i, item := range t.Items {
		style := t.Styles.InactiveTab
		if i == t.Active {
			style = t.Styles.ActiveTab
		}
		nodes = append(nodes, tabNode{
			text:     item,
			style:    style,
			isActive: i == t.Active,
		})
	}

	if len(t.Items) > 1 && t.RightHint != "" {
		nodes = append(nodes, tabNode{
			text:     t.RightHint + " →",
			style:    t.Styles.HintTab,
			isActive: false,
		})
	}

	var renderedTabs []string

	for i, node := range nodes {
		isFirst := i == 0

		border := node.style.GetBorderStyle()

		if isFirst && node.isActive {
			border.BottomLeft = "│"
		} else if isFirst && !node.isActive {
			border.BottomLeft = "├"
		}

		style := node.style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(node.text))
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
