package djafka

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ResetOffsetPrompt struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	logger     *log.Logger
	groupName  string
	topicName  string
}

func InitialResetOffsetPrompt(log *log.Logger, group string, topic string) ResetOffsetPrompt {
	m := ResetOffsetPrompt{
		inputs:    make([]textinput.Model, 1),
		logger:    log,
		groupName: group,
		topicName: topic,
	}

	var t textinput.Model

	t = textinput.New()
	t.CursorStyle = cursorStyle
	t.CharLimit = 32
	t.Focus()
	t.Placeholder = "Offset"
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.CharLimit = 5
	t.Validate = validateTopic

	m.inputs[0] = t

	return m
}

func (m ResetOffsetPrompt) Init() tea.Cmd {
	return textinput.Blink
}

func (m ResetOffsetPrompt) Empty() ResetOffsetPrompt {
	return ResetOffsetPrompt{}
}

func (m ResetOffsetPrompt) Update(msg tea.Msg) (ResetOffsetPrompt, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			res := AddTopicCancel{}
			m.logger.Println("Submiting AddTopicCancel", res)
			return m, func() tea.Msg { return res }

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				parsed, err := strconv.ParseInt(m.inputs[0].Value(), 10, 64)
				if err != nil {
					return m, func() tea.Msg { return ErrorMsg(err) } //panic is annoying af
				}
				res := ResetOffsetMsg{
					m.groupName,
					m.topicName,
					parsed,
				}
				m.logger.Println("Submiting Values", res)
				return m, func() tea.Msg { return res }
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *ResetOffsetPrompt) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m ResetOffsetPrompt) View() string {
	var b strings.Builder

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\t%s\n\n", *button)

	return fmt.Sprintf(
		`
	 %s
	 %s
	 %s %s
	`,
		inputStyle.Width(30).Render("Offset"),
		m.inputs[0].View(),
		"",
		b.String(),
	) + "\n"

}
