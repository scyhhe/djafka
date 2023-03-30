package djafka

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StartupComponent struct {
	progress.Model
	Width  int
	Height int
}

const (
	max = 500
	min = 100
)

func tick() tea.Cmd {
	return func() tea.Msg {
		rng := rand.Intn(max-min) + min
		time.Sleep(time.Millisecond * time.Duration(rng))
		return TickMsg{}
	}
}

func NewStartupComponent() (StartupComponent, tea.Cmd) {
	c := StartupComponent{Model: progress.New(progress.WithDefaultGradient())}
	return c, tick()
}

func (c StartupComponent) Update(msg tea.Msg) (StartupComponent, tea.Cmd) {
	if c.Initialized() {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.Width = msg.Width
		c.Height = msg.Height
	case progress.FrameMsg:
		newModel, cmd := c.Model.Update(msg)
		c.Model = newModel.(progress.Model)
		return c, cmd
	case TickMsg:
		cmd := c.Model.IncrPercent(0.1)
		return c, tea.Batch(tick(), cmd)
	}

	return c, nil
}

func (c StartupComponent) View() string {
	progress := c.Model.View()
	text := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		Width(50).
		Align(lipgloss.Center).
		Render("Initializing ...")

	box := lipgloss.JoinVertical(lipgloss.Center, progress, text)

	return lipgloss.Place(c.Width, c.Height, lipgloss.Center, lipgloss.Center, box)
}

func (c *StartupComponent) Initialized() bool {
	return c.Model.Percent() == 1.0
}
