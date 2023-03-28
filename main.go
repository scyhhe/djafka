package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/scyhhe/djafka/internal/djafka"
)

type sessionState uint

const (
	connectionView sessionState = iota
	selectionView
	topicList
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	state           sessionState
	connectionTable table.Model
	selectionTable  table.Model
	topicLIst       list.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.selectionTable.Focused() {
				m.selectionTable.Blur()
			} else {
				m.selectionTable.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.state == connectionView {
				m.state = selectionView
			} else {
				m.state = connectionView
			}
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.focusedComponent().SelectedRow()[0]),
			)
		}
	}

	focusedViews := &m.connectionTable
	switch m.state {
	case connectionView:
		focusedViews = &m.connectionTable
	case selectionView:
		focusedViews = &m.selectionTable
	}

	*focusedViews, cmd = focusedViews.Update(msg)

	return m, tea.Batch(append(cmds, cmd)...)
}

func (m *model) focusedComponent() *table.Model {
	switch m.state {
	case connectionView:
		return &m.connectionTable
	case selectionView:
		return &m.selectionTable
	default:
		panic("unhandled state")
	}
}

func (m model) View() string {
	connectionViewStyle := baseStyle.Copy()
	selectionTableStyle := baseStyle.Copy()

	if m.state == connectionView {
		connectionViewStyle = makeFocused(connectionViewStyle)
	} else {
		selectionTableStyle = makeFocused(selectionTableStyle)
	}

	return lipgloss.JoinVertical(lipgloss.Left, connectionViewStyle.Render(m.connectionTable.View()), selectionTableStyle.Render(m.selectionTable.View()))
}

func makeFocused(s lipgloss.Style) lipgloss.Style {
	return s.BorderForeground(lipgloss.Color("69"))
}

func buildTable(cols []table.Column, rows []table.Row) table.Model {

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		MarginTop(0).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return t
}

func main() {
	service, err := djafka.NewService()
	if err != nil {
		panic(err)
	}
	defer service.Close()

	topColumns := []table.Column{
		{Title: "Connection", Width: 30},
	}
	middleColumns := []table.Column{
		{Title: "Topic", Width: 30},
	}

	connectionRows := []table.Row{{"pretty connection"}, {"lol connection"}, {"wonky connection"}}
	selectionRows := []table.Row{{"Topics"}, {"Consumer Groups"}, {"Info"}}

	connTable := buildTable(topColumns, connectionRows)
	selectionTable := buildTable(middleColumns, selectionRows)

	topicList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	m := model{connectionView, connTable, selectionTable, topicList}

	if _, err := tea.NewProgram(&m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
