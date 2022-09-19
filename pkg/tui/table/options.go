package table

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

func WithHelp() Option {
	return func(m *Model) {
	}
}

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

func WithRefresh(fn func() ([]Column, []Row)) Option {
	return func(m *Model) {
		m.refreshFunc = fn
	}
}

func WithAppending(row *Row) Option {
	return func(m *Model) {
		m.appendRow = row
	}
}

func WithStickyCursor(sticky bool) Option {
	return func(m *Model) {
		m.stickyCursor = sticky
	}
}
