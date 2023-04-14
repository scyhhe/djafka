package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TableComponent struct {
	table.Model
	selectionChanged bool
}

func NewTableComponent(cols []table.Column) TableComponent {
	return TableComponent{
		Model:            buildEmptyTable(cols),
		selectionChanged: false,
	}
}

func (c TableComponent) Update(msg tea.Msg) (TableComponent, tea.Cmd) {
	previousRow := c.SelectedRow()
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable

	currentRow := c.SelectedRow()
	if currentRow != nil && previousRow != nil {
		c.selectionChanged = (*currentRow)[0] != (*previousRow)[0]
	} else {
		c.selectionChanged = false
	}

	return c, cmd
}

func (c *TableComponent) View() string {
	defocusTable(&c.Model)
	if c.Focused() {
		focusTable(&c.Model)
	}

	return c.Model.View()
}

func (c *TableComponent) SetRows(rows []table.Row) {
	c.Model.SetRows(rows)
	c.SetCursor(0)
}

func (c *TableComponent) SelectedRow() *table.Row {
	if len(c.Rows()) < 1 {
		return nil
	}

	row := c.Model.SelectedRow()
	return &row
}

func (c *TableComponent) SelectionChanged() bool {
	return c.selectionChanged
}
