package djafka

import (
	"fmt"

	_ "embed"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

//go:embed info.md
var infoContent string

type InfoComponent struct {
	viewport.Model
}

const width = 68

func NewInfoComponent() (InfoComponent, error) {
	vp := viewport.New(width, 40)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return InfoComponent{}, fmt.Errorf("Failed to initialise new renderer: %w", err)
	}

	str, err := renderer.Render(infoContent)
	if err != nil {
		return InfoComponent{}, fmt.Errorf("Failed to reander content: %w", err)
	}

	vp.SetContent(str)

	return InfoComponent{vp}, nil
}

func (c InfoComponent) Update(msg tea.Msg) (InfoComponent, tea.Cmd) {
	newModel, cmd := c.Model.Update(msg)
	c.Model = newModel

	resizeMsg, isResized := msg.(tea.WindowSizeMsg)
	if isResized {
		c.Height = resizeMsg.Height - 1
	}

	return c, cmd
}

func (c InfoComponent) View() string {
	return c.Model.View()
}
