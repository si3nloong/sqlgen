package templates

import (
	"go/types"
)

type ModelTmplParams struct {
	GoPkg string

	// Models
	Models []*Model
}

type Model struct {
	GoName string

	Name string

	// Primary key
	PK *Field

	Fields []*Field

	Func []*Func
}

type Field struct {
	GoName string

	Name string

	Type types.Type

	Tag []string

	Index uint
}

type Func struct {
	// Function name
	Name    string
	Recv    []string
	Returns []string
}
