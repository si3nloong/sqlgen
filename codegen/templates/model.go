package templates

import (
	"go/types"

	"github.com/samber/lo"
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
	// Sql database name
	DatabaseName *string

	// Sql table name
	TableName string

	IsAutoIncr bool

	Keys []*Field
	// Sql columns
	Fields          []*Field
	HasDatabaseName bool
	// Is model implement `Tabler` interface
	HasTableName bool
	// Is model implement `Columner` interface
	HasColumn bool
	// // Is model implement `Rower` interface
	// HasRow bool
}

func (m Model) IsCompositeKey() bool {
	return len(m.Keys) > 1
}

func (m Model) HasNotOnlyPK() bool {
	for i := range m.Fields {
		if !lo.Contains(m.Keys, m.Fields[i]) {
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
