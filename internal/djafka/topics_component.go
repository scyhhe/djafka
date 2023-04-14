package djafka

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TopicConfigLoader func(topic Topic) (TopicConfig, error)

type TopicsComponent struct {
	table.Model
	topics               map[string]Topic
	topicConfig          *TopicConfig
	topicConfigComponent table.Model
	topicConfigLoader    TopicConfigLoader
	needsReload          bool
}

func NewTopicsComponent(loader TopicConfigLoader) TopicsComponent {
	return TopicsComponent{
		Model: buildEmptyTable([]table.Column{
			{Title: TopicsLabel, Width: 30},
			{Title: "# of Partitions", Width: 30},
		}),
		topics:      map[string]Topic{},
		topicConfig: nil,
		topicConfigComponent: buildEmptyTable([]table.Column{
			{Title: "Topic Name", Width: 30},
			{Title: "Offset", Width: 20},
			{Title: "Partition", Width: 10},
		}),
		topicConfigLoader: loader,
	}
}

func (c TopicsComponent) Update(msg tea.Msg) (TopicsComponent, tea.Cmd) {
	cmds := []tea.Cmd{}

	if c.needsReload {
		cmds = append(cmds, c.loadTopicConfig(c.SelectedTopic()))
		c.needsReload = false
	}
	newTable, cmd := c.Model.Update(msg)
	c.Model = newTable
	cmds = append(cmds, cmd)

	// TODO: Handle TopicConfigLoaded message

	return c, tea.Batch(cmds...)
}

func (c *TopicsComponent) View() string {
	defocusTable(&c.Model)
	if c.Focused() {
		focusTable(&c.Model)
	}

	// TODO: Handle focus for sub-component

	resultPane := baseStyle.Render(c.Model.View())
	detailsPane := baseStyle.Render(c.topicConfigComponent.View())

	return lipgloss.JoinVertical(lipgloss.Right, resultPane, detailsPane)
}

func (c *TopicsComponent) SetTopics(topics []Topic) {
	c.needsReload = true

	topicMap := map[string]Topic{}
	for _, topic := range topics {
		topicMap[topic.Name] = topic
	}
	c.topics = topicMap

	rows := []table.Row{}
	for _, topic := range c.topics {
		rows = append(rows, table.Row{topic.Name, strconv.Itoa(topic.PartitionCount)})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })

	c.Model.SetRows(rows)
	c.Model.SetCursor(0)
}

func (c *TopicsComponent) SelectedTopic() *Topic {
	if len(c.Rows()) < 1 {
		return nil
	}

	topic := c.topics[c.SelectedRow()[0]]
	return &topic
}

func (c *TopicsComponent) loadTopicConfig(topic *Topic) tea.Cmd {
	if topic == nil {
		return nil
	}

	return func() tea.Msg {
		config, err := c.topicConfigLoader(*topic)
		if err != nil {
			return ErrorMsg(err)
		}

		return TopicSettingsLoadedMsg(config)
	}
}
