package form

type Section struct {
	Fields []*Field
}

type Field struct {
	Name  string
	Value string
}
