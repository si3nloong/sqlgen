package sqlite

import (
	"fmt"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *sqliteDriver) ColumnDataTypes() map[string]*dialect.ColumnType {
	return map[string]*dialect.ColumnType{
		"rune": {
			DataType: s.columnDataType("CHAR(1)"),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "{{addrOfGoPath}}",
		},
		"byte": {
			DataType: s.columnDataType("CHAR(1)"),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "{{addrOfGoPath}}",
		},
		"string": {
			DataType: func(col dialect.GoColumn) string {
				size := 255
				if n := col.Size(); n > 0 {
					size = n
				}
				return fmt.Sprintf("VARCHAR(%d)", size)
			},
			Valuer:  "string({{goPath}})",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"bool": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "bool({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"int": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int8": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int16": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int32": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int64": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint8": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint16": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint32": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint64": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"float32": {
			DataType: s.columnDataType("REAL"),
			Valuer:   "float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"float64": {
			DataType: s.columnDataType("REAL"),
			Valuer:   "float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"time.Time": {
			DataType: s.columnDataType("TEXT"),
			Valuer:   "(time.Time)({{goPath}})",
			Scanner:  "(*time.Time)({{addrOfGoPath}})",
		},
	}
}

func (*sqliteDriver) columnDataType(dataType string) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		if !column.GoNullable() {
			dataType += " NOT NULL"
		}
		return dataType
	}
}
