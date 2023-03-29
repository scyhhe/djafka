package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func selectConsumer() tea.Cmd {
	return func() tea.Msg {
		return ConsumerSelectedMsg{}
	}
}

type ResultComponent struct {
	table.Model
}

func (c ResultComponent) Update(msg tea.Msg) (ResultComponent, tea.Cmd) {
	// prevRow := c.SelectedRow()[0]
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable
	// currentRow := c.SelectedRow()[0]

	// hasRowChanged := prevRow != currentRow

	// if hasRowChanged {
	// 	selectConsumer()
	// }

	return c, cmd
}

func (c *ResultComponent) SetTopics(items []string) {
	rows := []table.Row{}
	for _, item := range items {
		// TODO add current # of messages
		rows = append(rows, table.Row{item})
	}
	c.Model.SetRows(rows)
}

func (c *ResultComponent) SetConsumers(items []Consumer) {
	rows := []table.Row{}
	for _, item := range items {
		rows = append(rows, table.Row{item.ConsumerId, item.GroupId, item.State})
	}
	c.Model.SetRows(rows)
}
