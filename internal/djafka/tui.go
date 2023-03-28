package djafka

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Consumer struct {
	GroupId         string
	ConsumerId      string
	State           string
	TopicPartitions []ConsumerTopicPartition
}

type ConsumerTopicPartition struct {
	TopicName string
	Offset    int64
	Partition int32
}

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
	selectionTable  Table
	topicLIst       list.Model
	config          *Config
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
		{Title: "Topics", Width: 30},
	}
	selectionRows := []table.Row{{"Topics"}, {"Consumer Groups"}, {"Info"}}

	connectionTable := buildTable(connectionColumns, connectionRows)
	selectionTable := buildTable(selectionColumns, selectionRows)

	topicList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	service, err := NewService()
	if err != nil {
		panic(err)
	}

	connectionComponent := ConnectionComponent{
		Model:  connectionTable,
		config: config,
	}

	*m = model{connectionState, connectionComponent, Table{selectionTable}, topicList, config, service}

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
				m.state = connectionState
			}
		}
	// Resizing
	case tea.WindowSizeMsg:
		m.connectionTable.SetHeight((msg.Height / 2) - 4)
		m.selectionTable.SetHeight((msg.Height / 2) - 4)

	// Custom messages
	case ConnectionChangedMsg:
		panic(msg)
	}

	_, cmd = m.focusedComponent().Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) focusedComponent() Component {
	switch m.state {
	case connectionState:
		return &m.connectionTable
	case selectionState:
		return &m.selectionTable
	default:
		panic("unhandled state")
	}
}

func (m *model) View() string {
	connectionBorderStyle := defocusTable(&m.connectionTable.Model)
	selectionBorderStyle := defocusTable(&m.selectionTable.Model)

	switch m.state {
	case connectionState:
		connectionBorderStyle = focusTable(&m.connectionTable.Model)
	case selectionState:
		selectionBorderStyle = focusTable(&m.selectionTable.Model)
	}

	menuPane := lipgloss.JoinVertical(lipgloss.Left, connectionBorderStyle.Render(m.connectionTable.View()),
		selectionBorderStyle.Render(m.selectionTable.View()))

	return menuPane
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
