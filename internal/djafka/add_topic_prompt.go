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
	"github.com/charmbracelet/lipgloss"
)

func validateTopic(s string) error {
	if len(s) > 64 {
		return fmt.Errorf("topic name is too long")
	}
	return nil
}

func validateInt(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

var (
	inputStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
	continueStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type AddTopicPrompt struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	logger     *log.Logger
}

func InitialAddTopicPrompt(log *log.Logger) AddTopicPrompt {
	m := AddTopicPrompt{
		inputs: make([]textinput.Model, 3),
		logger: log,
	}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Focus()
			t.Placeholder = ""
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
			t.Validate = validateTopic
		case 1:
			t.Placeholder = ""
			t.Validate = validateInt
		case 2:
			t.Placeholder = ""
			t.Validate = validateInt
		}

		m.inputs[i] = t
	}

	return m
}

func (m AddTopicPrompt) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddTopicPrompt) Update(msg tea.Msg) (AddTopicPrompt, tea.Cmd) {
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
				parsed1, err := strconv.ParseInt(m.inputs[1].Value(), 10, 64)
				if err != nil {
					return m, func() tea.Msg { return ErrorMsg(err) } //panic is annoying af
				}
				parsed2, err := strconv.ParseInt(m.inputs[2].Value(), 10, 64)
				if err != nil {
					return m, func() tea.Msg { return ErrorMsg(err) }
				}
				res := AddTopicSubmitMsg{
					m.inputs[0].Value(),
					int(parsed1),
					int(parsed2),
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

func (m *AddTopicPrompt) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AddTopicPrompt) View() string {
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
	 %s
	 %s
	 %s
	 %s
	 %s %s
	`,
		inputStyle.Width(30).Render("Topic Name"),
		m.inputs[0].View(),
		inputStyle.Width(12).Render("Partitions"),
		m.inputs[1].View(),
		inputStyle.Width(20).Render("Max Replication"),
		m.inputs[2].View(),
		"",
		b.String(),
	) + "\n"

}
