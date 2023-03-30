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

func selectInfo() tea.Cmd {
	return func() tea.Msg {
		return InfoSelectedMsg{}
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
		if currentRow == TopicsLabel {
			return m, tea.Batch(cmd, selectTopics())
		} else if currentRow == ConsumerGroupsLabel {
			return m, tea.Batch(cmd, selectConsumers())
		} else if currentRow == InfoLabel {
			return m, tea.Batch(cmd, selectInfo())
		}
	}

	return m, cmd
}

func (m *Menu) IsInfoSelected() bool {
	return m.SelectedRow()[0] == InfoLabel
}
