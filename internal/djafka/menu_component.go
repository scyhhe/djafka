package djafka

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuComponent struct {
	table.Model
	selectionChanged bool
}

func NewMenuComponent() MenuComponent {
	menuColumns := []table.Column{{Title: MenuLabel, Width: 30}}
	menuRows := []table.Row{{TopicsLabel}, {ConsumerGroupsLabel}, {InfoLabel}}

	return MenuComponent{
		Model:            buildTable(menuColumns, menuRows),
		selectionChanged: false,
	}
}

func (c MenuComponent) Update(msg tea.Msg) (MenuComponent, tea.Cmd) {
	previousEntry := c.SelectedEntry()
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable

	c.selectionChanged = c.SelectedEntry() != previousEntry

	return c, cmd
}

func (c *MenuComponent) View() string {
	defocusTable(&c.Model)
	if c.Focused() {
		focusTable(&c.Model)
	}

	return c.Model.View()
}

func (c *MenuComponent) SelectedEntry() string {
	if len(c.Rows()) < 1 {
		return ""
	}

	return c.SelectedRow()[0]
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

func (c *MenuComponent) SelectionChanged() bool {
	return c.selectionChanged
}
