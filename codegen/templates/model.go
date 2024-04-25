package templates

import (
	"go/types"
)

type ModelTmplParams struct {
	Models      []*Model
	OmitGetters bool
}

type PK struct {
	IsAutoIncr bool
	Field      *Field
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

func (m Model) HasNotOnlyPK() bool {
	for i := range m.Fields {
		if m.PK != nil && m.Fields[i] != m.PK.Field {
			return true
		}
	}
	return false
}

type Field struct {
	// Struct property name
	GoName string
	// Struct property path
	GoPath     string
	ColumnName string

	Type types.Type

	CustomMarshaler   string
	CustomUnmarshaler string

	IsBinary          bool
	IsTextMarshaler   bool
	IsTextUnmarshaler bool
	Size              int
	Index             int
}
