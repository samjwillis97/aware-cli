package table

import "github.com/charmbracelet/bubbles/help"

// WithColumns sets the table columns (headers).
func WithColumns(cols []Column) Option {
	return func(m *Model) {
		m.cols = cols
	}
}

// WithRows sets the table rows (data).
func WithRows(rows []Row) Option {
	return func(m *Model) {
		m.rows = rows
	}
}

// WithHeight sets the height of the table.
func WithHeight(h int) Option {
	return func(m *Model) {
		m.viewport.Height = h
	}
}

// WithWidth sets the width of the table.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.viewport.Width = w
	}
}

// WithFocused sets the focus state of the table.
func WithFocused(f bool) Option {
	return func(m *Model) {
		m.focus = f
	}
}

// WithStyles sets the table styles.
func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

// WithKeyMap sets the key map.
func WithKeyMap(km KeyMap) Option {
	return func(m *Model) {
		m.KeyMap = km
	}
}

// WithHelp sets whether to show help.
func WithHelp() Option {
	return func(m *Model) {
		m.helpEnabled = true
		m.help = help.New()
	}
}

// WithFullscreen sets whether to fullscreen the table.
func WithFullscreen(fullscreen bool) Option {
	return func(m *Model) {
		m.fullscreen = fullscreen
	}
}

// WithAutoWidth must be called after Columns and Rows have been set.
func WithAutoWidth(autowidth bool) Option {
	return func(m *Model) {
		m.autowidth = autowidth
	}
}

// WithRefresh sets the function to call when refreshing the table.
func WithRefresh(fn func() ([]Column, []Row)) Option {
	return func(m *Model) {
		m.refreshFunc = fn
	}
}

// WithAppending allows setting of a pointer to append from on command.
func WithAppending(row *Row) Option {
	return func(m *Model) {
		m.appendRow = row
	}
}

// WithStickyCursor when enabled will keep the cursor on the bottom row.
func WithStickyCursor(sticky bool) Option {
	return func(m *Model) {
		m.stickyCursor = sticky
	}
}

// WithCopyIndex sets which row to copy to the clipboard.
func WithCopyIndex(index int) Option {
	return func(m *Model) {
		m.copyIndex = index
	}
}
