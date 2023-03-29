package djafka

import (
	"fmt"
	"log"
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
	logger          *log.Logger
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

	connectionComponent := ConnectionComponent{
		Model:  connectionTable,
		config: config,
	}

	menu := Menu{
		Model: selectionTable,
	}

	resultComponent := ResultComponent{
		Model: resultTable,
	}

	*m = model{m.logger, connectionState, connectionComponent, resultComponent, menu, nil}

	return changeConnection(config.Connections[0])
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Printf("Received tea.Msg: %T\n", msg)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.connectionTable.Blur()
	m.selectionTable.Blur()
	m.resultComponent.Blur()

	switch m.state {
	case connectionState:
		m.connectionTable.Focus()
	case selectionState:
		m.selectionTable.Focus()
	case topicListState:
		m.resultComponent.Focus()
	default:
		panic("unhandled state")
	}

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
		m.resultComponent.SetHeight((msg.Height / 2) - 4)

	// Custom messages
	case ConnectionChangedMsg:
		cmd := m.changeConnection(Connection(msg))
		cmds = append(cmds, cmd)
	case TopicsSelectedMsg:
		cmd := m.loadTopics()
		cmds = append(cmds, cmd)
	case TopicsLoadedMsg:
		m.resultComponent.SetItems(msg)
	}

	m.connectionTable, cmd = m.connectionTable.Update(msg)
	cmds = append(cmds, cmd)
	m.selectionTable, cmd = m.selectionTable.Update(msg)
	cmds = append(cmds, cmd)
	m.resultComponent, cmd = m.resultComponent.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) changeConnection(conn Connection) tea.Cmd {
	return func() tea.Msg {
		if m.service != nil {
			m.service.Close()
		}

		service, err := NewService(conn)
		if err != nil {
			panic(fmt.Errorf("Failed to re-create service: %w", err))
		}

		m.service = service

		return ClientConnectedMsg{}
	}
}

func (m *model) loadTopics() tea.Cmd {
	return func() tea.Msg {
		topics, err := m.service.ListTopics()
		if err != nil {
			// TODO: Do properly :')
			panic(err)
		}

		return TopicsLoadedMsg(topics)
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

	rightPane := resultBorderStyle.Render(m.resultComponent.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, menuPane, rightPane)
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
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	logger := log.Default()
	logger.SetOutput(f)

	if _, err := tea.NewProgram(&model{logger: logger}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
