package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TopicsListedMsg []string

func listTopics(s *Service) tea.Cmd {
	return func() tea.Msg {
		topics, err := s.ListTopics()
		if err != nil {
			// TODO: Do properly :')
			panic(err)
		}

		return TopicsListedMsg(topics)
	}
}

type Menu struct {
	table.Model
	service *Service
}

func (m *Menu) Update(msg tea.Msg) (Component, tea.Cmd) {
	prevRow := m.SelectedRow()[0]
	newTable, cmd := m.Model.Update(msg)
	m.Model = newTable
	currentRow := m.SelectedRow()[0]

	hasRowChanged := prevRow != currentRow

	if hasRowChanged {
		if currentRow == "Topics" {
			return m, tea.Batch(cmd, listTopics(m.service))
		}
	}

	return m, cmd
}
