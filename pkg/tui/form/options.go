package form

func WithSections(sections []Section) Option {
	return func(m *Model) {
		m.sections = sections
	}
}

func WithRefreshFunc(fn func() []Field) Option {
	return func(m *Model) {
		m.refreshFunc = fn
	}
}
