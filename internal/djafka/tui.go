package djafka

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
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
	errorState
	addTopicState
	resetOffsetState
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	logger            *log.Logger
	state             sessionState
	previousState     sessionState
	errorComponent    ErrorComponent
	connectionTable   ConnectionComponent
	resultComponent   ResultComponent
	detailsComponent  DetailsComponent
	selectionTable    Menu
	service           *Service
	help              HelpComponent
	infoComponent     InfoComponent
	startupComponent  StartupComponent
	addTopicPrompt    AddTopicPrompt
	resetOffsetPrompt ResetOffsetPrompt
	selectedConsumer  *Consumer
	selectedTopic     *Topic
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

	help := buildHelp()

	help.FullHelpView(defaultKeys.FullHelp())

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

	addTopicPrompt := InitialAddTopicPrompt(m.logger)

	resetOffsetPrompt := InitialResetOffsetPrompt(m.logger)

	detailsComponent := DetailsComponent{
		Model: detailsTable,
	}

	helpComponent := HelpComponent{
		Model: help,
	}

	infoComponent, err := NewInfoComponent()
	if err != nil {
		panic(err)
	}

	startupComponent, cmd := NewStartupComponent()

	*m = model{
		logger:            m.logger,
		state:             connectionState,
		previousState:     connectionState,
		errorComponent:    ErrorComponent{},
		connectionTable:   connectionComponent,
		resultComponent:   resultComponent,
		detailsComponent:  detailsComponent,
		selectionTable:    menu,
		service:           nil,
		help:              helpComponent,
		infoComponent:     infoComponent,
		startupComponent:  startupComponent,
		addTopicPrompt:    addTopicPrompt,
		resetOffsetPrompt: resetOffsetPrompt,
		selectedConsumer:  nil,
		selectedTopic:     nil,
	}

	return tea.Batch(changeConnection(config.Connections[0]), cmd)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	_, isAddTopicPromptResult := msg.(AddTopicSubmitMsg)
	_, isResetOffsetPromptResult := msg.(ResetOffsetMsg)
	_, isAddTopicCancel := msg.(AddTopicCancel)
	_, isResetOffsetCancel := msg.(ResetOffsetCancel)

	if m.state == errorState {
		m.errorComponent, cmd = m.errorComponent.Update(msg)
		cmds = append(cmds, cmd)

		keyMsg, isKeyMsg := msg.(tea.KeyMsg)
		if isKeyMsg {
			switch keyMsg.String() {
			case QUIT, CANCEL:
				return m, tea.Quit
			default:
				m.restoreState()
				cmds = append(cmds, reset())
			}
		}

		return m, tea.Batch(cmds...)
	} else if isResetOffsetCancel {
		m.restoreState()
	} else if m.state == addTopicState && !isAddTopicPromptResult && !isAddTopicCancel {
		m.addTopicPrompt, cmd = m.addTopicPrompt.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	} else if m.state == resetOffsetState && !isResetOffsetPromptResult {
		m.resetOffsetPrompt, cmd = m.resetOffsetPrompt.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	m.connectionTable.Blur()
	m.selectionTable.Blur()
	m.resultComponent.Blur()
	m.detailsComponent.Blur()

	switch m.state {
	case connectionState:
		m.connectionTable.Focus()
	case selectionState:
		m.selectionTable.Focus()
	case resultState:
		m.resultComponent.Focus()
	case addTopicState:
	case resetOffsetState:
	case detailsState:
		m.detailsComponent.Focus()
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
				m.state = detailsState
			case detailsState:
				m.state = connectionState
			}
		case "ctrl+t":
			m.state = addTopicState
		case "ctrl+o":
			if m.state == resetOffsetState {
				m.restoreState()
			}
			m.resetOffsetPrompt = InitialResetOffsetPrompt(m.logger)
			m.state = resetOffsetState
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		}

	// Resizing
	case tea.WindowSizeMsg:
		m.connectionTable.SetHeight((msg.Height / 2) - 4)
		m.selectionTable.SetHeight((msg.Height / 2) - 4)
		m.resultComponent.SetHeight((msg.Height / 2) - 4)
		m.detailsComponent.SetHeight((msg.Height / 2) - 4)

	// Custom messages
	case ConnectionChangedMsg:
		cmd := m.changeConnection(Connection(msg))
		cmds = append(cmds, cmd)
	case TopicsSelectedMsg:
		m.resultComponent.SetRows([]table.Row{})
		m.resultComponent.SetColumns([]table.Column{
			{Title: TopicsLabel, Width: 30},
			{Title: "# of Partitions", Width: 30},
		})
		cmd := m.loadTopics()
		cmds = append(cmds, cmd)
	case TopicsLoadedMsg:
		m.resultComponent.SetTopics(msg)
	case TopicSelectedMsg:
		m.detailsComponent.SetRows([]table.Row{})
		m.detailsComponent.SetColumns([]table.Column{
			{Title: "Key", Width: 30},
			{Title: "Value", Width: 30},
		})
		cmd := m.loadTopicSettings(msg.Name)
		cmds = append(cmds, cmd)
		m.logger.Println("Saving selected topic with name: ", msg.Name)
		m.selectedTopic = &Topic{msg.Name, msg.PartitionCount}
	case TopicSettingsLoadedMsg:
		m.detailsComponent.SetTopicDetails(TopicConfig(msg))
	case ConsumersLoadedMsg:
		m.resultComponent.SetConsumers(msg)
	case ConsumersSelectedMsg:
		m.resultComponent.SetRows([]table.Row{})
		m.resultComponent.SetColumns([]table.Column{
			{Title: ConsumerIdLabel, Width: 30},
			{Title: GroupIdLabel, Width: 20},
			{Title: StateLabel, Width: 10},
		})
		cmd := m.loadConsumers()
		cmds = append(cmds, cmd)
	case ConsumerSelectedMsg:
		m.detailsComponent.SetRows([]table.Row{})
		m.detailsComponent.SetColumns([]table.Column{
			{Title: "Topic Name", Width: 30},
			{Title: "Offset", Width: 20},
			{Title: "Partition", Width: 10},
		})
		m.detailsComponent.SetConsumerDetails(Consumer(msg))
		m.selectedConsumer = &Consumer{msg.GroupId, msg.ConsumerId, msg.State, msg.TopicPartitions}
	case ErrorMsg:
		m.triggerErrorState(msg)
	case AddTopicCancel:
		m.logger.Println("Received AddTopicCancel")
		m.restoreState()
	case AddTopicSubmitMsg:
		m.logger.Println("Received AddTopicSubmitMsg with values: ", msg.name, msg.paritions, msg.replicationFactor)
		_, err := m.service.CreateTopic(msg.name, msg.paritions, msg.replicationFactor)
		if err != nil {
			cmds = append(cmds, sendErrorCmd(fmt.Errorf("Failed to create topig: %w", err)))
		}
		cmd := m.loadTopics()
		cmds = append(cmds, cmd)
		m.restoreState()
	case ResetOffsetMsg:
		m.logger.Println("Received ResetOffsetMsg with: ", msg.consumerGroup, msg.topicName, msg.offset)
		err := m.service.ResetConsumerOffsets(msg.consumerGroup, msg.topicName, msg.offset)
		if err != nil {
			cmds = append(cmds, sendErrorCmd(fmt.Errorf("Failed to reset offset: %w", err)))
		}
		m.restoreState()
	}

	m.errorComponent, cmd = m.errorComponent.Update(msg)
	cmds = append(cmds, cmd)
	m.connectionTable, cmd = m.connectionTable.Update(msg)
	cmds = append(cmds, cmd)
	m.selectionTable, cmd = m.selectionTable.Update(msg)
	cmds = append(cmds, cmd)
	m.resultComponent, cmd = m.resultComponent.Update(msg)
	cmds = append(cmds, cmd)
	m.detailsComponent, cmd = m.detailsComponent.Update(msg)
	cmds = append(cmds, cmd)
	m.infoComponent, cmd = m.infoComponent.Update(msg)
	cmds = append(cmds, cmd)
	m.startupComponent, cmd = m.startupComponent.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) triggerErrorState(err error) {
	m.previousState = m.state
	m.state = errorState
	m.errorComponent.Message = err.Error()
}

func (m *model) restoreState() {
	m.state = m.previousState

	m.addTopicPrompt = InitialAddTopicPrompt(m.logger) //reset prompt
}

func (m *model) changeConnection(conn Connection) tea.Cmd {
	return func() tea.Msg {
		if m.service != nil {
			m.service.Close()
		}

		service, err := NewService(conn, m.logger)
		if err != nil {
			return ErrorMsg(fmt.Errorf("Failed to re-create service: %w", err))
		}

		m.service = service

		return ClientConnectedMsg{}
	}
}

func (m *model) loadTopics() tea.Cmd {
	return func() tea.Msg {
		topics, err := m.service.ListTopics()
		if err != nil {
			return ErrorMsg(err)
		}

		return TopicsLoadedMsg(topics)
	}
}

func (m *model) loadTopicSettings(name string) tea.Cmd {
	return func() tea.Msg {
		config, err := m.service.GetTopicConfig(name)
		if err != nil {
			return ErrorMsg(err)
		}

		return TopicSettingsLoadedMsg(config)
	}
}

func (m *model) loadConsumers() tea.Cmd {
	return func() tea.Msg {
		consumerGroups, err := m.service.ListConsumerGroups()
		if err != nil {
			return ErrorMsg(err)
		}
		consumers, err := m.service.ListConsumers(consumerGroups)
		if err != nil {
			return ErrorMsg(err)
		}

		return ConsumersLoadedMsg(consumers)
	}
}

func sendErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

func reset() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{}
	}
}

func (m *model) View() string {
	if !m.startupComponent.Initialized() {
		return m.startupComponent.View()
	}

	if m.state == errorState {
		return m.errorComponent.View()
	} else if m.state == addTopicState {
		return m.addTopicPrompt.View()
	}

	connectionBorderStyle := defocusTable(&m.connectionTable.Model)
	selectionBorderStyle := defocusTable(&m.selectionTable.Model)
	resultBorderStyle := defocusTable(&m.resultComponent.Model)
	detailsBorderStyle := defocusTable(&m.detailsComponent.Model)

	helpView := m.help.View(defaultKeys)

	switch m.state {
	case connectionState:
		connectionBorderStyle = focusTable(&m.connectionTable.Model)
	case selectionState:
		selectionBorderStyle = focusTable(&m.selectionTable.Model)
	case resultState:
		resultBorderStyle = focusTable(&m.resultComponent.Model)
	case detailsState:
		detailsBorderStyle = focusTable(&m.detailsComponent.Model)
	}

	menuPane := lipgloss.JoinVertical(lipgloss.Left, connectionBorderStyle.Render(m.connectionTable.View()),
		selectionBorderStyle.Render(m.selectionTable.View()))

	resultPane := resultBorderStyle.Render(m.resultComponent.View())

	var resetPane string

	if m.state != resetOffsetState {
		resetPane = ""
	} else {
		resetPane = m.resetOffsetPrompt.View()
	}

	detailsPane := detailsBorderStyle.Render(m.detailsComponent.View())
	resultPane = lipgloss.JoinVertical(lipgloss.Right, resultPane, detailsPane, helpView)

	if m.selectionTable.IsInfoSelected() {
		resultPane = lipgloss.JoinVertical(lipgloss.Right, m.infoComponent.View(), helpView)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, menuPane, resultPane, resetPane)
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

func buildHelp() help.Model {
	h := help.New()

	return h
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
