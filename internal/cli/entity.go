package cli

type Entity struct {
	// Go package name
	Pkg string

	// Struct name
	Name string

	// Imported packages
	Imports []string

	// Struct field list
	FieldList []*Field
}

type Field struct {
	Name       string
	Column     string
	ActualType string
	Type       string
	IsScanner  bool
	IsValuer   bool
}
