package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Table struct {
	table.Model
}

func (t *Table) Update(msg tea.Msg) (Component, tea.Cmd) {
	newTable, cmd := t.Model.Update(msg)
	t.Model = newTable
	return t, cmd
}
