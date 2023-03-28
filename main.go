package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/scyhhe/djafka/internal/djafka"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	connectionTable table.Model
	selectionTable  table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.selectionTable.SelectedRow()[0]),
			)
		}
	}
	m.selectionTable, cmd = m.selectionTable.Update(msg)
	return m, cmd
}

func (m model) View() string {
	leftPane := lipgloss.JoinVertical(lipgloss.Left, baseStyle.Render(m.connectionTable.View()), baseStyle.Render(m.selectionTable.View()))

	return leftPane
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

	topRows := []table.Row{{"pretty connection"}, {"lol connection"}, {"wonky connection"}}
	middleRows := []table.Row{{"Topics"}, {"Consumer Groups"}, {"Info"}}

	topTable := buildTable(topColumns, topRows)
	middleTable := buildTable(middleColumns, middleRows)

	m := model{topTable, middleTable}

	topicName, err := service.CreateTopic("tschusch")
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", topicName)

	topics, err := service.ListTopics()
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", topics)

	consumerGroups, err := service.ListConsumerGroups()
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", consumerGroups)

	consumers, err := service.ListConsumers(consumerGroups)
	if err != nil {
		panic(err)
	}

	fmt.Println("consumers:", consumers)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
