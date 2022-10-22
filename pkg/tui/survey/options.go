package survey

// TODO: func Ask
// TODO: func AskOne

func WithQuestions(qs []*Question) Option {
	return func(m *Model) {
		m.questions = qs
	}
}

func WithResponse(response *interface{}) Option {
	return func(m *Model) {
		m.response = response
	}
}
