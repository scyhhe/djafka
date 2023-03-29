package djafka

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func selectConsumer(c Consumer) tea.Cmd {
	return func() tea.Msg {
		return ConsumerSelectedMsg(c)
	}
}

func selectTopic(t Topic) tea.Cmd {
	return func() tea.Msg {
		return TopicSelectedMsg(t)
	}
}

type ResultComponent struct {
	table.Model
	consumers  map[string]Consumer
	isConsumer bool
}

func (c ResultComponent) Update(msg tea.Msg) (ResultComponent, tea.Cmd) {

	switch msg := msg.(type) {
	case ConsumersLoadedMsg:
		c.SetConsumers(msg)
		c.isConsumer = true

		consumers := map[string]Consumer{}

		for _, item := range msg {
			consumers[item.ConsumerId] = item
		}

		c.consumers = consumers
		return c, selectConsumer(msg[0])
	case TopicsLoadedMsg:
		c.SetTopics(msg)
		c.isConsumer = false
		return c, selectTopic(msg[0])

	default:
		if len(c.Rows()) > 0 {
			prevRow := c.SelectedRow()[0]
			newTable, cmd := c.Model.Update(msg)
			c.Model = newTable
			currentRow := c.SelectedRow()[0]

			if prevRow != currentRow {
				if c.isConsumer {
					return c, tea.Batch(cmd, selectConsumer(c.consumers[currentRow]))
				}
				return c, tea.Batch(cmd, selectTopic(Topic{Name: currentRow}))
			}
		}
	}

	return c, nil
}

func (c *ResultComponent) SetTopics(items []Topic) {
	rows := []table.Row{}
	for _, item := range items {
		rows = append(rows, table.Row{item.Name, strconv.Itoa(item.PartitionCount)})
	}
	c.Model.SetRows(rows)
	c.Model.SetCursor(0)
}

func (c *ResultComponent) SetConsumers(items []Consumer) {
	rows := []table.Row{}
	for _, item := range items {
		rows = append(rows, table.Row{item.ConsumerId, item.GroupId, item.State})
	}
	c.Model.SetRows(rows)
	c.Model.SetCursor(0)
}
