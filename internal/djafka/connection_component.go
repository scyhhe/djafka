package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func changeConnection(conn Connection) tea.Cmd {
	return func() tea.Msg {
		return ConnectionChangedMsg(conn)
	}
}

type ConnectionComponent struct {
	table.Model
	config *Config
}

func (c ConnectionComponent) Update(msg tea.Msg) (ConnectionComponent, tea.Cmd) {
	prevRow := c.SelectedRow()[0]
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable
	currentRow := c.SelectedRow()[0]

	_, isReset := msg.(ResetMsg)
	if prevRow != currentRow || isReset {
		conn, err := c.config.FindConnection(currentRow)
		if err != nil {
			return c, sendErrorCmd(err)
		}

		return c, tea.Batch(cmd, changeConnection(conn))
	}

	return c, cmd
}
