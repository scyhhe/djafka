package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuComponent struct {
	table.Model
}

func (m MenuComponent) Update(msg tea.Msg) (MenuComponent, tea.Cmd) {
	newTable, cmd := m.Model.Update(msg)
	m.Model = newTable

	return m, cmd
}

func (m *MenuComponent) View() string {
	logger.Printf("Focused: %t\n", m.Focused())
	defocusTable(&m.Model)
	if m.Focused() {
		focusTable(&m.Model)
	}

	return m.Model.View()
}

func (m *MenuComponent) IsTopicsSelected() bool {
	return m.SelectedRow()[0] == TopicsLabel
}

func (m *MenuComponent) IsConsumerGroupssSelected() bool {
	return m.SelectedRow()[0] == ConsumerGroupsLabel
}

func (m *MenuComponent) IsInfoSelected() bool {
	return m.SelectedRow()[0] == InfoLabel
}
