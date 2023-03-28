package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionChangedMsg string

func changeConnection(url string) tea.Cmd {
	return func() tea.Msg {
		return ConnectionChangedMsg(url)
	}
}

type ConnectionComponent struct {
	table.Model
	config *Config
}

func (c *ConnectionComponent) Update(msg tea.Msg) (Component, tea.Cmd) {
	prevRow := c.SelectedRow()[0]
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable
	currentRow := c.SelectedRow()[0]

	if prevRow != currentRow {
		url, err := c.config.FindConnection(currentRow)
		if err != nil {
			// TODO: Properly handle error
			panic(err)
		}

		return c, tea.Batch(cmd, changeConnection(url))
	}

	return c, cmd
}
