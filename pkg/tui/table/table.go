// Package table contains a bubbletea interface for an interactive table.
package table

// https://github.com/charmbracelet/bubbles/blob/master/table/table.go
// See above for example.
import (
	"fmt"
	"strings"

	"ampaware.com/cli/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"golang.design/x/clipboard"
)

// Model defines a state for the table widget.
type Model struct {
	KeyMap KeyMap

	cols         []Column
	rows         []Row
	cursor       int
	focus        bool
	styles       Styles
	fullscreen   bool
	autowidth    bool
	stickyCursor bool
	helpEnabled  bool
	help         help.Model

	refreshFunc func() ([]Column, []Row)
	copyIndex   int

	appendRow *Row

	viewport viewport.Model
}

// CommandMessage is a type to efficiently contain messages to send to the table from an external source.
type CommandMessage byte

const (
	// AppendReady is used to let the table know a new row is ready to append.
	AppendReady CommandMessage = iota
)

// Row represents one line in the table.
type Row []string

// Column defines the table structure.
type Column struct {
	Title string
	Width int
}

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu menu.
type KeyMap struct {
	LineUp       key.Binding
	LineDown     key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding
	ToggleFocus  key.Binding
	Exit         key.Binding
	Execute      key.Binding
	Refresh      key.Binding
	ToClipboard  key.Binding
	Copy         key.Binding
	Paste        key.Binding
	ToggleHelp   key.Binding
}

// Styles contains style definitions for this list component. By default, these
// values are generated by DefaultStyles.
type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
	Footer   lipgloss.Style
}

// Option is used to set options in New.
type Option func(*Model)

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d", "½ page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		ToggleFocus: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("ESC", "toggle focus"),
		),
		Exit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "exit"),
		),
		Execute: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open item"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "R"),
			key.WithHelp("r", "refresh"),
		),
		ToClipboard: key.NewBinding(
			key.WithKeys("c", "C"),
			key.WithHelp("c", "to clipboard"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit, k.ToggleHelp}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown}, // First Column
		{k.PageUp, k.PageDown},
		{k.GotoTop, k.GotoBottom},
		{k.Execute, k.Refresh},
		{k.ToggleHelp, k.Exit},
	}
}

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")),
		Header: lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderTop(true).
			BorderBottom(true).
			Bold(false),
		Cell: lipgloss.NewStyle().Padding(0, 1),
		Footer: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderTop(true).
			BorderBottom(true).
			Bold(false),
	}
}

// SetStyles sets the table styles.
func (m *Model) SetStyles(s Styles) {
	m.styles = s
	m.UpdateViewport()
}

// New creates a new model for the table widget.
func New(opts ...Option) Model {
	// TODO: Add a nice Header (Optional)
	// TODO: Add an open/enter/execute function (Optional)
	// TODO: Add a delete function (Optional)
	// TODO: Add a filter (Optional)
	// TODO: Better Footers
	// TODO: Add an option for columns to overflow
	// TODO: Show help
	// TODO: Value to Clipboard
	// TODO: Method to Append for Telemetry Generation!
	// TODO: Better status of table - to use in footer

	// Maybe use channels for comms

	// TODO: WithAppend() could take a channel and a pointer to a []Row
	// When a new messages comes on the channel we read from the pointer

	m := Model{
		cursor:   0,
		viewport: viewport.New(0, 20),

		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.UpdateViewport()

	return m
}

// Init is the Bubble Tea entrypoint.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is the Bubble Tea update loop.
// nolint:gocyclo // This requires refactoring to simplify.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CommandMessage:
		{
			if msg == AppendReady {
				m.AppendRow()
			}
		}
	case tea.WindowSizeMsg:
		requiredPadding := 0
		requiredPadding += 3 // Headers
		requiredPadding += 3 // Footer
		if m.helpEnabled {
			requiredPadding++ // Short Help
		}
		m.SetHeight(msg.Height - requiredPadding)
		m.SetWidth(msg.Width)
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		case key.Matches(msg, m.KeyMap.ToggleFocus):
			if m.Focused() {
				m.Blur()
			} else {
				m.Focus()
			}
		case key.Matches(msg, m.KeyMap.Refresh):
			m.Refresh() // Channel Might work, need to make sure states are set correctly though
		case key.Matches(msg, m.KeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ToClipboard):
			m.CopyToClipboard()
		case key.Matches(msg, m.KeyMap.Execute):
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.SelectedRow()[1]),
			)
		case key.Matches(msg, m.KeyMap.ToggleHelp):
			if m.helpEnabled {
				if m.help.ShowAll {
					m.SetHeight(m.viewport.Height + 1)
				} else {
					m.SetHeight(m.viewport.Height - 1)
				}
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// Focused returns the focus state of the table.
func (m Model) Focused() bool {
	return m.focus
}

// Focus focusses the table, allowing the user to move around the rows and
// interact.
func (m *Model) Focus() {
	m.focus = true
	m.UpdateViewport()
}

// Blur blurs the table, preventing selection or movement.
func (m *Model) Blur() {
	m.focus = false
	m.UpdateViewport()
}

// View renders the component.
func (m Model) View() string {
	view := m.headersView()
	view += "\n" + m.viewport.View()
	view += "\n" + m.footersView()
	if m.helpEnabled {
		view += "\n" + m.help.View(m.KeyMap)
	}
	return view
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *Model) UpdateViewport() {
	m.UpdateColumnWidths()

	// TODO: Stick Cursor
	m.stickCursor()
	stickCursorToBottom := m.stickyCursor && (m.cursor == len(m.rows)-2)
	if stickCursorToBottom {
		m.cursor = len(m.rows) - 1
	}

	renderedRows := make([]string, 0, len(m.rows))
	for i := range m.rows {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)

	if stickCursorToBottom {
		m.viewport.GotoBottom()
	}
}

func (m *Model) stickCursor() {
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m Model) SelectedRow() Row {
	return m.rows[m.cursor]
}

// SetRows set a new rows state.
func (m *Model) SetRows(r []Row) {
	m.rows = r
	m.UpdateViewport()
}

// SetWidth sets the width of the viewport of the table.
func (m *Model) SetWidth(w int) {
	m.viewport.Width = w
	m.UpdateViewport()
}

// SetHeight sets the height of the viewport of the table.
func (m *Model) SetHeight(h int) {
	m.viewport.Height = h
	m.UpdateViewport()
}

// Height returns the viewport height of the table.
func (m Model) Height() int {
	return m.viewport.Height
}

// Width returns the viewport width of the table.
func (m Model) Width() int {
	return m.viewport.Width
}

// Cursor returns the index of the selected row.
func (m Model) Cursor() int {
	return m.cursor
}

// SetCursor sets the cursor position in the table.
func (m *Model) SetCursor(n int) {
	m.cursor = clamp(n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// MoveUp moves the selection up by any number of row.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, len(m.rows)-1)
	m.UpdateViewport()

	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	}
}

// Refresh executes the given refresh function and re sets the data.
func (m *Model) Refresh() {
	m.cols, m.rows = m.refreshFunc()
	m.UpdateViewport()
}

// AppendRow gets the row from the appendRow pointer and adds it to
// the existing data.
func (m *Model) AppendRow() {
	m.rows = append(m.rows, *m.appendRow)
	m.UpdateViewport()
}

// UpdateColumnWidths will automatically set the column Widths.
// nolint:gocyclo // This requires refactoring to simplify.
func (m *Model) UpdateColumnWidths() {
	if m.autowidth && m.viewport.Width > 0 {
		// The padding may be what is causing this to be out of whack
		// TODO Address
		availableWidth := m.viewport.Width - 10
		evenWidth := m.viewport.Width / len(m.cols)
		minWidth := 10
		var widths []int

		// should find ideal widths first time
		for _, col := range m.cols {
			widths = append(widths, len(col.Title)+1)
		}
		for _, row := range m.rows {
			for i, cell := range row {
				width := len(cell)
				if width > widths[i] {
					widths[i] = width
				}
			}
		}

		usedWidth := 0
		// Determine the total width used and set any less than minimum to minimum
		for i, width := range widths {
			if width < minWidth {
				widths[i] = minWidth
			}
			usedWidth += widths[i]
		}

		// Optimize the table widths until the usedWidth matches totalWidth
		// This will do for now
		maxLoops := 1000
		for iters := 0; usedWidth != availableWidth; iters++ {
			newWidth := 0
			for _, width := range widths {
				newWidth += width
			}
			if newWidth == availableWidth || iters > maxLoops {
				break
			}
			if newWidth > availableWidth {
				// Need to Collapse
				for i, width := range widths {
					// Iterate over the current widths
					// Find all the widths greater than desired
					// Use some smart maths to size them accordingly
					// to the rest of the room left
					if width > evenWidth {
						widths[i] = evenWidth
					}
				}
			}

			if newWidth < availableWidth {
				// Need to Expand
				// This branch works could probably do with some optimization
				addToEach := (availableWidth - newWidth) / len(widths)
				newUsedWidth := (addToEach * len(widths)) + newWidth

				if newUsedWidth < availableWidth {
					for i := 0; i < (availableWidth - newUsedWidth); i++ {
						widths[i]++
					}
				} else if newUsedWidth > availableWidth {
					for i := 0; i < (newUsedWidth - availableWidth); i++ {
						widths[i]--
					}
				}

				// TODO: Check Go array access, can I assign to width?
				// rather than use the index
				for i, width := range widths {
					widths[i] = width + addToEach
				}
			}
		}

		for i, width := range widths {
			m.cols[i].Width = width
		}
	}
}

// MoveDown moves the selection down by any number of row.
// It can not go below the last row.
func (m *Model) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, len(m.rows)-1)
	m.UpdateViewport()

	if m.cursor > (m.viewport.YOffset + (m.viewport.Height - 1)) {
		m.viewport.SetYOffset(m.cursor - (m.viewport.Height - 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *Model) GotoTop() {
	m.MoveUp(m.cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model) GotoBottom() {
	m.MoveDown(len(m.rows))
}

// CopyToClipboard will copy the value at the copyIndex for the current row the cursor is on.
func (m *Model) CopyToClipboard() {
	err := clipboard.Init()
	utils.ExitIfError(err)

	clipboard.Write(clipboard.FmtText, []byte(m.rows[m.cursor][m.copyIndex]))
}

// FromValues create the table rows from a simple string. It uses `\n` by
// default for getting all the rows and the given separator for the fields on
// each row.
func (m *Model) FromValues(value, separator string) {
	rows := []Row{}
	for _, line := range strings.Split(value, "\n") {
		r := Row{}
		for _, field := range strings.Split(line, separator) {
			r = append(r, field)
		}
		rows = append(rows, r)
	}

	m.SetRows(rows)
}

func (m Model) headersView() string {
	s := make([]string, 0, len(m.cols))
	for _, col := range m.cols {
		style := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		renderedCell := style.Render(runewidth.Truncate(col.Title, col.Width, "…"))
		s = append(s, m.styles.Header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m Model) footersView() string {
	// TODO: Put line above this
	style := lipgloss.NewStyle().Width(m.viewport.Width).MaxWidth(m.viewport.Width)
	rendered := style.Render(fmt.Sprintf("Showing %d entries", len(m.rows)))
	return m.styles.Footer.Render(rendered)
	// return fmt.Sprintf("Showing %d entries", len(m.rows))
}

func (m *Model) renderRow(rowID int) string {
	s := make([]string, 0, len(m.cols))
	for i, value := range m.rows[rowID] {
		style := lipgloss.NewStyle().Width(m.cols[i].Width).MaxWidth(m.cols[i].Width).Inline(true)
		renderedCell := m.styles.Cell.Render(style.Render(runewidth.Truncate(value, m.cols[i].Width, "…")))
		s = append(s, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, s...)

	if rowID == m.cursor {
		return m.styles.Selected.Render(row)
	}

	return row
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
