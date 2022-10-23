package form

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	KeyMap KeyMap
	styles Styles
	help   help.Model

	sections []Section

	refreshFunc func() []Field
	cursor      int
}

type KeyMap struct {
	LineUp     key.Binding
	LineDown   key.Binding
	Exit       key.Binding
	ToggleHelp key.Binding
}

type Styles struct {
	FieldText lipgloss.Style
	ValueText lipgloss.Style
	Heading   lipgloss.Style
}

type Option func(*Model)

func DefaultKeyMap() KeyMap {
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Exit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "exit"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit, k.ToggleHelp}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown}, // First Column
		{k.ToggleHelp, k.Exit},
	}
}

func DefaultStyles() Styles {
	return Styles{
		FieldText: lipgloss.NewStyle().Padding(0, 1).Bold(true),
		ValueText: lipgloss.NewStyle().Padding(0, 1),
	}
}

func New(opts ...Option) Model {
	m := Model{
		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// Init is the Bubble Tea entrypoint.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is the Bubble Tea update loop.
// nolint:gocyclo // This requires refactoring to simplify.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ToggleHelp):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	sections := ""
	for _, section := range m.sections {
		fieldView := ""
		valueView := ""
		for _, field := range section.Fields {
			fieldView = lipgloss.JoinVertical(lipgloss.Top, fieldView, m.styles.FieldText.Render(field.Name))
			valueView = lipgloss.JoinVertical(lipgloss.Top, valueView, m.styles.ValueText.Render(field.Value))
		}
		fields := lipgloss.JoinHorizontal(lipgloss.Top, fieldView, valueView)

		sections = lipgloss.JoinVertical(
			lipgloss.Top,
			sections,
			lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true).Render(fields),
		)
	}
	return sections
}

func (m *Model) MoveDown(n int) {
}

func (m *Model) MoveUp(n int) {
}
