package djafka

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	New   key.Binding
	Reset key.Binding
	Quit  key.Binding
}

var defaultKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	New: key.NewBinding(
		key.WithKeys("new", "n"),
		key.WithHelp("n", "new topic"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset offset"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type HelpComponent struct {
	help.Model
	keyMap
}

func (c HelpComponent) Update(msg tea.Msg) (HelpComponent, tea.Cmd) {
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable

	return c, cmd
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.New, k.Reset},                // second column
		{k.Help, k.Quit},                // third column
	}
}
