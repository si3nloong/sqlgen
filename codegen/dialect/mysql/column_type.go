package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *mysqlDriver) ColumnDataTypes() map[string]*dialect.ColumnType {
	dataTypes := map[string]*dialect.ColumnType{
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
			DataType: s.columnDataType("VARCHAR(255)", ""),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"bool": {
			DataType: s.columnDataType("BOOL", false),
			Valuer:   "(bool)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"float32": {
			DataType: s.columnDataType("FLOAT", int64(0)),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float32({{addrOfGoPath}})",
		},
		"float64": {
			DataType: s.columnDataType("FLOAT", int64(0)),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float64({{addrOfGoPath}})",
		},
		"time.Time": {
			DataType: s.columnDataType("TIMESTAMP", sql.RawBytes("CURRENT_TIMESTAMP")),
			Valuer:   "(time.Time)({{goPath}})",
			Scanner:  "(*time.Time)({{addrOfGoPath}})",
		},
		"*string": {
			DataType: s.columnDataType("VARCHAR(255)"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.String({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"*[]byte": {
			DataType: s.columnDataType("BLOB"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.String({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"*bool": {
			DataType: s.columnDataType("BOOL"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Bool({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"*float32": {
			DataType: s.columnDataType("FLOAT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float32({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float32({{addrOfGoPath}})",
		},
		"*float64": {
			DataType: s.columnDataType("FLOAT"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float64({{addrOfGoPath}})",
		},
		"*time.Time": {
			DataType: s.columnDataType("TIMESTAMP"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Time({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Time({{addr}})",
		},
		"sql.RawBytes": {
			DataType: s.columnDataType("TEXT"),
			Valuer:   "database/sql.RawBytes({{goPath}})",
			Scanner:  "(*database/sql.RawBytes)({{addrOfGoPath}})",
		},
		"encoding/json.RawMessage": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "(string)({{goPath}})",
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
		"[]string": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalStringSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.StringSlice({{addrOfGoPath}})",
		},
		"[]byte": {
			DataType: s.columnDataType("BLOB"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"[]bool": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalBoolSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.BoolSlice({{addrOfGoPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalIntSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntSlice({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalIntSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntSlice({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalIntSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntSlice({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalIntSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntSlice({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalIntSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntSlice({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUintSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintSlice({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUintSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintSlice({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUintSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintSlice({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUintSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintSlice({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalUintSlice({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintSlice({{addrOfGoPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalFloatList({{goPath}},-1)",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float32Slice({{addrOfGoPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/encoding.MarshalFloatList({{goPath}},-1)",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float64Slice({{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler({{addrOfGoPath}})",
		},
	}
	s.mapIntegers(dataTypes)
	return dataTypes
}

func (s *mysqlDriver) mapIntegers(dict map[string]*dialect.ColumnType) {
	types := [][2]string{
		{"int", "INTEGER"}, {"int8", "TINYINT"}, {"int16", "SMALLINT"}, {"int32", "MEDIUMINT"}, {"int64", "BIGINT"},
		{"uint", "INTEGER UNSIGNED"}, {"uint8", "TINYINT UNSIGNED"}, {"uint16", "SMALLINT UNSIGNED"}, {"uint32", "MEDIUMINT UNSIGNED"}, {"uint64", "BIGINT UNSIGNED"},
	}
	for _, t := range types {
		dict[t[0]] = &dialect.ColumnType{
			DataType: s.columnDataType(t[1], int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		}
	}
	types = [][2]string{
		{"*int", "INTEGER"}, {"*int8", "TINYINT"}, {"*int16", "SMALLINT"}, {"*int32", "MEDIUMINT"}, {"*int64", "BIGINT"},
		{"*uint", "INTEGER UNSIGNED"}, {"*uint8", "TINYINT UNSIGNED"}, {"*uint16", "SMALLINT UNSIGNED"}, {"*uint32", "MEDIUMINT UNSIGNED"}, {"*uint64", "BIGINT UNSIGNED"},
	}
	for _, t := range types {
		dict[t[0]] = &dialect.ColumnType{
			DataType: s.columnDataType(t[1], int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		}
	}
}

func (*mysqlDriver) columnDataType(dataType string, defaultValue ...any) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		str := dataType
		if !column.GoNullable() {
			str += " NOT NULL"
		}
		if len(defaultValue) > 0 {
			str += " DEFAULT " + format(defaultValue[0])
		}
		// if c.extra != "" {
		// 	str += " " + c.extra
		// }
		return str
	}
}

func format(v any) string {
	switch vi := v.(type) {
	case string:
		return "'" + vi + "'"
	case bool:
		return strconv.FormatBool(vi)
	case int:
		return strconv.Itoa(vi)
	case int64:
		return strconv.FormatInt(vi, 10)
	case uint64:
		return strconv.FormatUint(vi, 10)
	case float32:
		return strconv.FormatFloat(float64(vi), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(vi, 'f', -1, 64)
	case sql.RawBytes:
		return unsafe.String(unsafe.SliceData(vi), len(vi))
	default:
		panic(fmt.Sprintf("unsupported data type %T", vi))
	}
}
