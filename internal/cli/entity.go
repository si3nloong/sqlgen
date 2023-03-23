package cli

type ImportPkg struct {
	Name  string
	Alias *string
	Path  string
}

type Entity struct {
	// Go package name
	GoPkg string

	// Go struct name
	GoName string

	// Table name
	Name string

	// Table primary key
	PrimaryKey *Field

	// Imported packages
	Imports []*ImportPkg

	// Struct field list
	Fields []*Field
}

type Field struct {
	GoName     string
	Name       string
	ActualType string
	Type       string
	IsScanner  bool
	IsValuer   bool
}
