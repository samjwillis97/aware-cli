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

func WithData(data [][]string) Option {
    return func(m *Model) {
        var cols []Column
        var rows []Row
        for i, row := range data {
            if i == 0 {
                for _, col := range row {
                    cols = append(cols, Column{Title: col, Width: 10})
                }
            } else {
                rows = append(rows, row)
            }
        }
        m.cols = cols
        m.rows = rows
    }
}

func WithFullscreen(fullscreen bool) Option {
    return func(m *Model) {
        m.fullscreen = fullscreen
    }
}

// WithAutoWidth must be called after Columns and Rows have been set
func WithAutoWidth(autowidth bool) Option {
    return func(m *Model) {
        m.autowidth = autowidth
    }
}

func WithFooterText(text string) Option {
    return func(m *Model) {
        m.footerText = text
    }
}
