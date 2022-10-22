package survey

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	questions []*Question
	index     int
	response  *interface{}
	options   PromptConfig
}

type Option func(*Model)

func New(opts ...Option) Model {
	m := Model{}

	m.options = defaultAskOptions()

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.index < len(m.questions) {
		return m.questions[m.index].Prompt.Ask(&m)
	}
	return "Test"
}
