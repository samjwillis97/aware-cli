package table

// https://github.com/charmbracelet/bubbles/blob/master/table/table.go
// See above for example

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Model defines a state for the table widget.
type model struct {
	table table.Model
}

// New creates a new model for the table widget.
func New(data [][]string, opts ...table.Option) model {
    var cols []table.Column
    var rows []table.Row
    var widths []int

    maxWidth := 20 // Derive this from the viewport somehow
    minWidth := 10
    for i, row := range data {
        if i == 0 {
            for _, col := range row {
                cols = append(cols, table.Column{Title: col})
                width := minWidth
                if maxWidth > len(col) && len(col) > minWidth {
                    width = len(col) + 1
                }
                widths = append(widths, width)
            }
        } else {
            rows = append(rows, row)
            for i, col := range row {
                if len(col) > widths[i] {
                    if maxWidth > len(col) {
                        widths[i] = len(col)
                    } else {
                        widths[i] = maxWidth
                    }
                }
            }
        }
    }

    for i, width := range widths {
        cols[i].Width = width
    }

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

    for _, opt := range opts {
        opt(&t)
    }

	m := model{t}

    return m
}

func WithHelp() table.Option {
    return func(*table.Model) {

    }
}

func (m model) Init() tea.Cmd {
    return nil
}

// Update is the Bubble Tea update loop.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the component.
func (m model) View() string {
    return baseStyle.Render(m.table.View()) + "\n"
}
