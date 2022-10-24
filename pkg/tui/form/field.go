package form

type Section struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name  string
	Value string
}
