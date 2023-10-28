package templates

import (
	"go/types"
)

type ModelTmplParams struct {
	// Models
	Models []*Model
}

type PK struct {
	Field *Field
}

type Model struct {
	// Go struct name
	GoName string
	// Sql table name
	TableName string
	// Primary key
	PK *PK
	// Sql columns
	Fields []*Field
	// Is model implement `Tabler` interface
	HasTableName bool
	// Is model implement `Columner` interface
	HasColumn bool
	// // Is model implement `Rower` interface
	// HasRow bool
}

type Field struct {
	// Struct property name
	GoName string
	// Struct property path
	GoPath     string
	ColumnName string

	Type types.Type

	Tag   []string
	Index int
}
