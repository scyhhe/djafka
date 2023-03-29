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
	resultState
	detailsState
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	logger          *log.Logger
	state           sessionState
	connectionTable ConnectionComponent
	resultTable     ResultComponent
	detailsTable    DetailsComponent
	selectionTable  Menu
	service         *Service
}

func (m *model) Init() tea.Cmd {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	connectionColumns := []table.Column{
		{Title: ConnectionsLabel, Width: 30},
	}

	connectionRows := []table.Row{}
	for _, connection := range config.Connections {
		connectionRows = append(connectionRows, table.Row{connection.Name})
	}

	selectionColumns := []table.Column{
		{Title: MenuLabel, Width: 30},
	}
	selectionRows := []table.Row{{TopicsLabel}, {ConsumerGroupsLabel}, {InfoLabel}}

	resultColumns := []table.Column{
		{Title: ResultLabel, Width: 60},
	}

	detailColumns := []table.Column{
		{Title: DetailsLabel, Width: 60},
	}

	resultRows := []table.Row{}
	detailRows := []table.Row{}

	connectionTable := buildTable(connectionColumns, connectionRows)
	selectionTable := buildTable(selectionColumns, selectionRows)
	resultTable := buildTable(resultColumns, resultRows)
	detailsTable := buildTable(detailColumns, detailRows)

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

	detailsComponent := DetailsComponent{
		Model: detailsTable,
	}

	*m = model{m.logger, connectionState, connectionComponent, resultComponent, detailsComponent, menu, nil}

	return changeConnection(config.Connections[0])
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Printf("Received tea.Msg: %T\n", msg)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.connectionTable.Blur()
	m.selectionTable.Blur()
	m.resultTable.Blur()
	m.detailsTable.Blur()

	switch m.state {
	case connectionState:
		m.connectionTable.Focus()
	case selectionState:
		m.selectionTable.Focus()
	case resultState:
		m.resultTable.Focus()
	case detailsState:
		m.detailsTable.Focus()
	default:
		panic("unhandled state")
	}

	switch msg := msg.(type) {
	// Key presses
	case tea.KeyMsg:
		switch msg.String() {
		case ESC:
			if m.selectionTable.Focused() {
				m.selectionTable.Blur()
			} else {
				m.selectionTable.Focus()
			}
		case QUIT, CANCEL:
			return m, tea.Quit
		case TAB:
			switch m.state {
			case connectionState:
				m.state = selectionState
			case selectionState:
				m.state = resultState
			case resultState:
				m.state = connectionState
			case detailsState:
				m.state = detailsState
			}
		}
	// Resizing
	case tea.WindowSizeMsg:
		m.connectionTable.SetHeight((msg.Height / 2) - 4)
		m.selectionTable.SetHeight((msg.Height / 2) - 4)
		m.resultTable.SetHeight((msg.Height / 2) - 4)

	// Custom messages
	case ConnectionChangedMsg:
		cmd := m.changeConnection(Connection(msg))
		cmds = append(cmds, cmd)
	case TopicsSelectedMsg:
		m.resultTable.SetRows([]table.Row{})
		m.resultTable.SetColumns([]table.Column{
			{Title: TopicsLabel, Width: 60},
		})
		cmd := m.loadTopics()
		cmds = append(cmds, cmd)
	case TopicsLoadedMsg:
		m.resultTable.SetTopics(msg)
	case ConsumersLoadedMsg:
		m.resultTable.SetConsumers(msg)
	case ConsumersSelectedMsg:
		m.resultTable.SetRows([]table.Row{})
		m.resultTable.SetColumns([]table.Column{
			{Title: ConsumerIdLabel, Width: 30},
			{Title: GroupIdLabel, Width: 20},
			{Title: StateLabel, Width: 10},
		})
		cmd := m.loadConsumers()
		cmds = append(cmds, cmd)
	case ConsumerSelectedMsg:
		m.detailsTable.SetRows([]table.Row{})
		m.detailsTable.SetColumns([]table.Column{
			{Title: "Topic Name", Width: 30},
			{Title: "Offset", Width: 20},
			{Title: "Partition", Width: 10},
		})
		m.detailsTable.SetConsumerDetails(Consumer(msg))
	}

	m.connectionTable, cmd = m.connectionTable.Update(msg)
	cmds = append(cmds, cmd)
	m.selectionTable, cmd = m.selectionTable.Update(msg)
	cmds = append(cmds, cmd)
	m.resultTable, cmd = m.resultTable.Update(msg)
	cmds = append(cmds, cmd)
	m.detailsTable, cmd = m.detailsTable.Update(msg)
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

func (m *model) loadConsumers() tea.Cmd {
	return func() tea.Msg {
		consumerGroups, err := m.service.ListConsumerGroups()
		if err != nil {
			// TODO: Do properly :')
			panic(err)
		}
		consumers, err := m.service.ListConsumers(consumerGroups)
		if err != nil {
			// TODO: Do properly :')
			panic(err)
		}

		return ConsumersLoadedMsg(consumers)
	}
}

func (m *model) View() string {
	connectionBorderStyle := defocusTable(&m.connectionTable.Model)
	selectionBorderStyle := defocusTable(&m.selectionTable.Model)
	resultBorderStyle := defocusTable(&m.resultTable.Model)
	detailsBorderStyle := defocusTable(&m.detailsTable.Model)

	switch m.state {
	case connectionState:
		connectionBorderStyle = focusTable(&m.connectionTable.Model)
	case selectionState:
		selectionBorderStyle = focusTable(&m.selectionTable.Model)
	case resultState:
		resultBorderStyle = focusTable(&m.resultTable.Model)
	case detailsState:
		detailsBorderStyle = focusTable(&m.detailsTable.Model)
	}

	menuPane := lipgloss.JoinVertical(lipgloss.Left, connectionBorderStyle.Render(m.connectionTable.View()),
		selectionBorderStyle.Render(m.selectionTable.View()))

	resultPane := resultBorderStyle.Render(m.resultTable.View())
	detailsPane := detailsBorderStyle.Render(m.detailsTable.View())

	resultPane = lipgloss.JoinVertical(lipgloss.Right, resultPane, detailsPane)

	return lipgloss.JoinHorizontal(lipgloss.Top, menuPane, resultPane)
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
