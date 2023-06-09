package templates

import "go/types"

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
}

type Field struct {
	GoName string

	Name string

	Type types.Type
}
