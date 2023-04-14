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
	TableComponent
	topics               map[string]Topic
	topicConfigComponent TableComponent
	topicConfigLoader    TopicConfigLoader
	needsReload          bool
}

func NewTopicsComponent(loader TopicConfigLoader) TopicsComponent {
	return TopicsComponent{
		TableComponent: NewTableComponent([]table.Column{
			{Title: TopicsLabel, Width: 30},
			{Title: "# of Partitions", Width: 30},
		}),
		topics: map[string]Topic{},
		topicConfigComponent: NewTableComponent([]table.Column{
			{Title: "Topic Name", Width: 30},
			{Title: "Offset", Width: 20},
			{Title: "Partition", Width: 10},
		}),
		topicConfigLoader: loader,
	}
}

func (c TopicsComponent) Update(msg tea.Msg) (TopicsComponent, tea.Cmd) {
	cmds := []tea.Cmd{}

	newTable, cmd := c.TableComponent.Update(msg)
	c.TableComponent = newTable
	cmds = append(cmds, cmd)

	if c.SelectionChanged() {
		c.needsReload = true
	}

	if c.needsReload {
		cmds = append(cmds, c.loadTopicConfig(c.SelectedTopic()))
		c.needsReload = false
	}

	switch msg := msg.(type) {
	case TopicSettingsLoadedMsg:
		c.setTopicConfig(TopicConfig(msg))
	}

	return c, tea.Batch(cmds...)
}

func (c *TopicsComponent) View() string {
	// TODO: Handle focus for sub-component
	c.topicConfigComponent.Blur()

	resultPane := baseStyle.Render(c.TableComponent.View())
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

	c.TableComponent.SetRows(rows)
	c.TableComponent.SetCursor(0)
}

func (c *TopicsComponent) SelectedTopic() *Topic {
	row := c.SelectedRow()
	if row == nil {
		return nil
	}

	topic := c.topics[(*row)[0]]
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

func (c *TopicsComponent) setTopicConfig(config TopicConfig) {
	rows := []table.Row{}
	for key, value := range config.Settings {
		rows = append(rows, table.Row{key, value})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	c.topicConfigComponent.SetRows(rows)
}
