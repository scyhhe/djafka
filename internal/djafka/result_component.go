package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ResultComponent struct {
	table.Model
}

func (c ResultComponent) Update(msg tea.Msg) (ResultComponent, tea.Cmd) {
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable

	return c, cmd
}

func (c *ResultComponent) SetItems(items []string) {
	rows := []table.Row{}
	for _, item := range items {
		// TODO add current # of messages
		rows = append(rows, table.Row{item})
	}
	c.Model.SetRows(rows)
}
