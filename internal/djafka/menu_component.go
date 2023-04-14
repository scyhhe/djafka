package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuComponent struct {
	TableComponent
}

func NewMenuComponent() MenuComponent {
	cols := []table.Column{{Title: MenuLabel, Width: 30}}
	rows := []table.Row{{TopicsLabel}, {ConsumerGroupsLabel}, {InfoLabel}}
	tableComponent := NewTableComponent(cols)
	tableComponent.SetRows(rows)

	return MenuComponent{
		TableComponent: tableComponent,
	}
}

func (c MenuComponent) Update(msg tea.Msg) (MenuComponent, tea.Cmd) {
	newTable, cmd := c.TableComponent.Update(msg)
	c.TableComponent = newTable

	return c, cmd
}

func (c *MenuComponent) SelectedEntry() string {
	row := c.SelectedRow()
	if row == nil {
		return ""
	}

	return (*row)[0]
}

func (c *MenuComponent) IsTopicsSelected() bool {
	return c.SelectedEntry() == TopicsLabel
}

func (c *MenuComponent) IsConsumerGroupssSelected() bool {
	return c.SelectedEntry() == ConsumerGroupsLabel
}

func (c *MenuComponent) IsInfoSelected() bool {
	return c.SelectedEntry() == InfoLabel
}
