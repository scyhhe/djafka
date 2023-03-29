package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func selectTopics() tea.Cmd {
	return func() tea.Msg {
		return TopicsSelectedMsg{}
	}
}

func selectConsumers() tea.Cmd {
	return func() tea.Msg {
		return ConsumersSelectedMsg{}
	}
}

type Menu struct {
	table.Model
}

func (m Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	prevRow := m.SelectedRow()[0]
	newTable, cmd := m.Model.Update(msg)
	m.Model = newTable
	currentRow := m.SelectedRow()[0]

	hasRowChanged := prevRow != currentRow

	_, isClientConnected := msg.(ClientConnectedMsg)
	if hasRowChanged || isClientConnected {
		if currentRow == "Topics" {
			return m, tea.Batch(cmd, selectTopics())
		} else if currentRow == "Consumer Groups" {
			return m, tea.Batch(cmd, selectConsumers())
		}
	}

	return m, cmd
}
