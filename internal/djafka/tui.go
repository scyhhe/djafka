package djafka

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DataProvider interface {
	ListTopics() ([]string, error)
	ListConsumerGroups() ([]string, error)
	ListConsumers(groupIds []string) ([]Consumer, error)
}

type Component interface {
	Update(tea.Msg) (Component, tea.Cmd)
}

type sessionState uint

const (
	connectionState sessionState = iota
	selectionState
	topicListState
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	state           sessionState
	connectionTable ConnectionComponent
	resultComponent ResultComponent
	selectionTable  Menu
	service         *Service
}

func (m *model) Init() tea.Cmd {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	connectionColumns := []table.Column{
		{Title: "Connections", Width: 30},
	}

	connectionRows := []table.Row{}
	for _, connection := range config.Connections {
		connectionRows = append(connectionRows, table.Row{connection.Name})
	}

	selectionColumns := []table.Column{
		{Title: "Menu", Width: 30},
	}
	selectionRows := []table.Row{{"Topics"}, {"Consumer Groups"}, {"Info"}}

	topicResultColumns := []table.Column{
		{Title: "Topics", Width: 70},
	}

	topicResultRows := []table.Row{}

	connectionTable := buildTable(connectionColumns, connectionRows)
	selectionTable := buildTable(selectionColumns, selectionRows)
	resultTable := buildTable(topicResultColumns, topicResultRows)

	// topicList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	service, err := NewService(config.Connections[0])
	if err != nil {
		panic(err)
	}

	connectionComponent := ConnectionComponent{
		Model:  connectionTable,
		config: config,
	}

	menu := Menu{
		Model:   selectionTable,
		service: service,
	}

	resultComponent := ResultComponent{
		Model: resultTable,
	}

	*m = model{connectionState, connectionComponent, resultComponent, menu, service}

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	// Key presses
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
			switch m.state {
			case connectionState:
				m.state = selectionState
			case selectionState:
				m.state = topicListState
			case topicListState:
				m.state = connectionState
			}
		}
	// Resizing
	case tea.WindowSizeMsg:
		m.connectionTable.SetHeight((msg.Height / 2) - 4)
		m.selectionTable.SetHeight((msg.Height / 2) - 4)

	// Custom messages
	case ConnectionChangedMsg:
		m.changeConnection(Connection(msg))
	case TopicsListedMsg:
		m.resultComponent.SetItems(msg)
	}

	_, cmd = m.focusedComponent().Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) changeConnection(conn Connection) {
	if m.service != nil {
		m.service.Close()
	}

	service, err := NewService(conn)
	if err != nil {
		panic(fmt.Errorf("Failed to re-create service: %w", err))
	}

	m.service = service
}

func (m *model) focusedComponent() Component {
	switch m.state {
	case connectionState:
		return &m.connectionTable
	case selectionState:
		return &m.selectionTable
	case topicListState:
		return &m.resultComponent
	default:
		panic("unhandled state")
	}
}

func (m *model) View() string {
	connectionBorderStyle := defocusTable(&m.connectionTable.Model)
	selectionBorderStyle := defocusTable(&m.selectionTable.Model)
	resultBorderStyle := defocusTable(&m.resultComponent.Model)

	switch m.state {
	case connectionState:
		connectionBorderStyle = focusTable(&m.connectionTable.Model)
	case selectionState:
		selectionBorderStyle = focusTable(&m.selectionTable.Model)
	case topicListState:
		resultBorderStyle = focusTable(&m.resultComponent.Model)
	}

	menuPane := lipgloss.JoinVertical(lipgloss.Left, connectionBorderStyle.Render(m.connectionTable.View()),
		selectionBorderStyle.Render(m.selectionTable.View()))

	return lipgloss.JoinHorizontal(lipgloss.Top, menuPane, resultBorderStyle.Render(m.resultComponent.View()))
}

func makeFocused(s lipgloss.Style) lipgloss.Style {
	return s.BorderForeground(lipgloss.Color("69"))
}

func defocusTable(t *table.Model) lipgloss.Style {
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		MarginTop(0).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("240")).
		Bold(false)

	t.SetStyles(s)
	return baseStyle.Copy()
}

func focusTable(t *table.Model) lipgloss.Style {
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
	return makeFocused(baseStyle.Copy())
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

func Run() {
	if _, err := tea.NewProgram(&model{}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
