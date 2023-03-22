package cli

type Entity struct {
	Pkg       string
	Name      string
	FieldList []*Field
	// Fields map[string]
}

type Field struct {
	Name      string
	Column    string
	Type      string
	BaseType  string
	IsScanner bool
	IsValuer  bool
}
