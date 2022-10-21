package survey

// Icon holds the text and format to show for a particular icon
type Icon struct {
	Text   string
	Format string
}

// IconSet holds the icons to use for various prompts
type IconSet struct {
	HelpInput      Icon
	Error          Icon
	Help           Icon
	Question       Icon
	MarkedOption   Icon
	UnmarkedOption Icon
	SelectFocus    Icon
}

// Transformer is a function passed to a Question after a user has provided a response.
// The function can be used to implement a custom logic that will result to return
// a different representation of the given answer.
//
type Transformer func(ans interface{}) (newAns interface{})

// Validator is a function passed to a Question after a user has provided a response.
// If the function returns an error, then the user will be prompted again for another
// response.
type Validator func(ans interface{}) error

// Question is the core data structure for a survey questionnaire.
type Question struct {
	Name      string
	Prompt    Prompt
	Validate  Validator
	Transform Transformer
}

// PromptConfig holds the global configuration for a prompt
type PromptConfig struct {
	PageSize         int
	Icons            IconSet
	HelpInput        string
	SuggestInput     string
	Filter           func(filter string, option string, index int) bool
	KeepFilter       bool
	ShowCursor       bool
	RemoveSelectAll  bool
	RemoveSelectNone bool
}

// Prompt is the primary interface for the objects that can take user input
// and return a response.
type Prompt interface {
	Prompt(config *PromptConfig) (interface{}, error) // TODO: Implement
	Cleanup(*PromptConfig, interface{}) error         // TODO: Implement
	Error(*PromptConfig, error) error                 // TODO: Implemenet
}
