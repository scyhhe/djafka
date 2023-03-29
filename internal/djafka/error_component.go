package djafka

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorComponent struct {
	Message string
	Width   int
	Height  int
}

func (c ErrorComponent) Update(msg tea.Msg) (ErrorComponent, tea.Cmd) {
	resizeMsg, isResized := msg.(tea.WindowSizeMsg)
	if isResized {
		c.Width = resizeMsg.Width
		c.Height = resizeMsg.Height
	}

	return c, nil
}

func (c ErrorComponent) View() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color("#ff0000")).
		Width(9).
		Align(lipgloss.Center).
		Render("Error")

	body := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Center).
		Render("Error")

	anyKey := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		Width(50).
		Align(lipgloss.Center).
		Render("Press any key to continue ...")

	box := lipgloss.JoinVertical(lipgloss.Center, header, body, anyKey)

	return lipgloss.Place(c.Width, c.Height, lipgloss.Center, lipgloss.Center, box)
}
