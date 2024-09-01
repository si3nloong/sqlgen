package mysql

import (
	"fmt"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *mysqlDriver) ColumnDataTypes() map[string]*dialect.ColumnType {
	return map[string]*dialect.ColumnType{
		"rune": {
			DataType: s.columnDataType("CHAR(1)"),
			Valuer:   "string({{goPath}})",
			Scanner:  "{{addrOfGoPath}}",
		},
		"string": {
			DataType: s.columnDataType("VARCHAR(255)"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"bool": {
			DataType: s.columnDataType("BOOL"),
			Valuer:   "bool({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"int": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int8": {
			DataType: s.columnDataType("TINYINT"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int16": {
			DataType: s.columnDataType("SMALLINT"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int32": {
			DataType: s.columnDataType("MEDIUMINT"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int64": {
			DataType: s.columnDataType("BIGINT"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint": {
			DataType: s.columnDataType("INTEGER UNSIGNED"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint8": {
			DataType: s.columnDataType("TINYINT UNSIGNED"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint16": {
			DataType: s.columnDataType("SMALLINT UNSIGNED"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint32": {
			DataType: s.columnDataType("MEDIUMINT UNSIGNED"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint64": {
			DataType: s.columnDataType("BIGINT UNSIGNED"),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"float32": {
			DataType: s.columnDataType("FLOAT"),
			Valuer:   "float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"float64": {
			DataType: s.columnDataType("FLOAT"),
			Valuer:   "float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"time.Time": {
			DataType: s.columnDataType("TIMESTAMP"),
			Valuer:   "time.Time({{goPath}})",
			Scanner:  "(*time.Time)({{addrOfGoPath}})",
		},
		"*string": {
			DataType: s.columnDataType("VARCHAR(255)"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.String({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfString({{addrOfGoPath}})",
		},
		"*[]byte": {
			DataType: s.columnDataType("BLOB"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.String({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfString({{addrOfGoPath}})"},
		"*bool": {
			DataType: s.columnDataType("BOOL"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Bool({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfBool({{addrOfGoPath}})"},
		"*uint": {
			DataType: s.columnDataType("INTEGER UNSIGNED"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*uint8": {
			DataType: s.columnDataType("TINYINT UNSIGNED"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*uint16": {
			DataType: s.columnDataType("SMALLINT UNSIGNED"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*uint32": {
			DataType: s.columnDataType("MEDIUMINT UNSIGNED"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*uint64": {
			DataType: s.columnDataType("BIGINT UNSIGNED"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*int": {
			DataType: s.columnDataType("INTEGER"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*int8": {
			DataType: s.columnDataType("TINYINT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*int16": {
			DataType: s.columnDataType("SMALLINT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*int32": {
			DataType: s.columnDataType("MEDIUMINT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*int64": {
			DataType: s.columnDataType("BIGINT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt({{addrOfGoPath}})",
		},
		"*float32": {
			DataType: s.columnDataType("FLOAT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfFloat({{addrOfGoPath}})",
		},
		"*float64": {
			DataType: s.columnDataType("FLOAT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfFloat({{addrOfGoPath}})",
		},
		"*time.Time": {
			DataType: s.columnDataType("TIMESTAMP"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Time({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfTime({{addrOfGoPath}})",
		},
		"sql.RawBytes": {
			DataType: s.columnDataType("TEXT"),
			Valuer:   "database/sql.RawBytes({{goPath}})",
			Scanner:  "(*database/sql.RawBytes)({{addrOfGoPath}})",
		},
		"encoding/json.RawMessage": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"[...]rune": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("VARCHAR(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.FixedSizeRunes({{goPath}}[:],{{len}})",
		},
		"[...]byte": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("CHAR(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.FixedSizeBytes({{goPath}}[:],{{len}})",
		},
		"[]byte": {
			DataType: s.columnDataType("BLOB"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"[]string": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalStringList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.StringList({{addrOfGoPath}})",
		},
		"[]bool": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalBoolList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.BoolList({{addrOfGoPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntList({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntList({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntList({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntList({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntList({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintList({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintList({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintList({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintList({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintList({{addrOfGoPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalFloatList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.FloatList({{addrOfGoPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalFloatList({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.FloatList({{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler({{addrOfGoPath}})",
		},
	}
}

func (*mysqlDriver) columnDataType(dataType string) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		if !column.Nullable() {
			dataType += " NOT NULL"
		}
		// if c.defaultValue != nil {
		// dataType += " DEFAULT " + format(c.defaultValue)
		// }
		// if c.extra != "" {
		// 	str += " " + c.extra
		// }
		return dataType
	}
}
