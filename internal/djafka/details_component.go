package djafka

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type DetailsComponent struct {
	table.Model
}

func (c DetailsComponent) Update(msg tea.Msg) (DetailsComponent, tea.Cmd) {
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable

	return c, cmd
}

func (c *DetailsComponent) SetConsumerDetails(item Consumer) {
	rows := []table.Row{}
	for _, item := range item.TopicPartitions {
		rows = append(rows, table.Row{item.TopicName, strconv.Itoa(int(item.Offset)), strconv.Itoa(int(item.Partition))})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	c.Model.SetRows(rows)
}

func (c *DetailsComponent) SetTopicDetails(item TopicConfig) {
	rows := []table.Row{}
	for key, value := range item.Settings {
		rows = append(rows, table.Row{key, value})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	c.Model.SetRows(rows)
}
