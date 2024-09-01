package sqlite

import (
	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *sqliteDriver) ColumnDataTypes() map[string]*dialect.ColumnType {
	return map[string]*dialect.ColumnType{
		"string": {
			DataType: s.columnDataType("TEXT"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
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
	}
}

func (*sqliteDriver) columnDataType(dataType string) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		if !column.Nullable() {
			dataType += " NOT NULL"
		}
		return dataType
	}
}
