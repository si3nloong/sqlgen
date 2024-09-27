package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *postgresDriver) ColumnDataTypes() map[string]*dialect.ColumnType {
	dataTypes := map[string]*dialect.ColumnType{
		"rune": {
			DataType: s.columnDataType("char(1)"),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "{{addrOfGoPath}}",
		},
		"byte": {
			DataType: s.columnDataType("char(1)"),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "{{addrOfGoPath}}",
		},
		"string": {
			DataType: func(col dialect.GoColumn) string {
				size := 255
				if n := col.Size(); n > 0 {
					size = n
				}
				return fmt.Sprintf("varchar(%d)", size)
			},
			Valuer:  "(string)({{goPath}})",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"bool": {
			DataType: s.columnDataType("bool", false),
			Valuer:   "(bool)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"int8": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int16": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int32": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int64": {
			DataType: s.columnDataType("bigint", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint8": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint16": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint32": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint64": {
			DataType: s.intDataType("bigint", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"float32": {
			DataType: s.columnDataType("real"),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"float64": {
			DataType: s.columnDataType("double precision"),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"time.Time": {
			DataType: s.columnDataType("timestamptz(6)", sql.RawBytes(`NOW()`)),
			Valuer:   "(time.Time)({{goPath}})",
			Scanner:  "(*time.Time)({{addrOfGoPath}})",
		},
		"*string": {
			DataType: func(col dialect.GoColumn) string {
				size := 255
				if n := col.Size(); n > 0 {
					size = n
				}
				return fmt.Sprintf("varchar(%d)", size)
			},
			Valuer:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"*bool": {
			DataType: s.columnDataType("bool", false),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"*float32": {
			DataType: s.columnDataType("real"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"*float64": {
			DataType: s.columnDataType("double precision"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"*time.Time": {
			DataType: s.columnDataType("timestamptz(6)"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Time({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Time({{addrOfGoPath}})",
		},
		"[...]rune": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("varchar(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.FixedSizeRunes(({{goPath}})[:],{{len}})",
		},
		"[...]byte": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("varchar(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: "github.com/si3nloong/sqlgen/sequel/types.FixedSizeBytes(({{goPath}})[:],{{len}})",
		},
		"[...]string": {
			DataType: s.columnDataType("text[{{len}}]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.StringArray[{{elemType}}]({{goPath}}[:])",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.StringArray[{{elemType}}]({{goPath}}[:])",
		},
		"[...]bool": {
			DataType: s.columnDataType("bool[{{len}}]"),
			Valuer:   "github.com/lib/pq.BoolArray({{goPath}}[:])",
			Scanner:  "github.com/lib/pq.Array(({{addrOfGoPath}})[:])",
		},
		"[]string": {
			DataType: s.columnDataType("text[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.StringArray[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.StringArrayScanner({{addrOfGoPath}})",
		},
		"[]rune": {
			DataType: s.columnDataType("text"),
			Valuer:   "string({{goPath}}[:])",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.FixedSizeRunes({{goPath}})",
		},
		"[]byte": {
			DataType: s.columnDataType("blob"),
			Valuer:   "{{goPath}}",
			Scanner:  "{{addrOfGoPath}}",
		},
		"[]bool": {
			DataType: s.columnDataType("bool[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.BoolArray[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.BoolArrayScanner({{addrOfGoPath}})",
		},
		"[][]byte": {
			DataType: s.columnDataType("bytea"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.ByteArray({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.ByteArray({{goPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.Float32Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.Float32Array[{{elemType}}]({{goPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.Float64Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.Float64Array[{{elemType}}]({{goPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayScanner({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayScanner({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.Int16Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayScanner({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.Int32Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayScanner({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.Int64Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.IntArrayScanner({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayScanner({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayScanner({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayScanner({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayScanner({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayValue({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types/pgtype.UintArrayScanner({{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("json"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler({{addrOfGoPath}})",
		},
	}
	s.mapIntegers(dataTypes)
	return dataTypes
}

func (s *postgresDriver) mapIntegers(dict map[string]*dialect.ColumnType) {
	types := [][2]string{
		{"int", "int4"}, {"int8", "int2"}, {"int16", "int2"}, {"int32", "int4"}, {"int64", "bigint"},
		{"uint", "int4"}, {"uint8", "int2"}, {"uint16", "int2"}, {"uint32", "int4"}, {"uint64", "bigint"},
	}
	for _, t := range types {
		dict[t[0]] = &dialect.ColumnType{
			DataType: s.intDataType(t[1], int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		}
	}
	types = [][2]string{
		{"*int", "int4"}, {"*int8", "int2"}, {"*int16", "int2"}, {"*int32", "int4"}, {"*int64", "bigint"},
		{"*uint", "int4"}, {"*uint8", "int2"}, {"*uint16", "int2"}, {"*uint32", "int4"}, {"*uint64", "bigint"},
	}
	for _, t := range types {
		dict[t[0]] = &dialect.ColumnType{
			DataType: s.intDataType(t[1], int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		}
	}
}

func (s *postgresDriver) intDataType(dataType string, defaultValue ...any) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		str := dataType
		if column.AutoIncr() {
			if strings.EqualFold(str, "integer") {
				str = strings.ReplaceAll(str, "integer", "serial")
			} else {
				str = strings.ReplaceAll(str, "int", "serial")
			}
		}
		// if !column.GoNullable() {
		// 	str += " NOT NULL"
		// }
		// // Auto increment cannot has default value
		// if !column.AutoIncr() && len(defaultValue) > 0 {
		// 	if !column.Key() {
		// 		str += " DEFAULT " + format(defaultValue[0])
		// 	}
		// 	switch any(defaultValue[0]).(type) {
		// 	case uint64:
		// 		str += " CHECK (" + s.QuoteIdentifier(column.ColumnName()) + " >= 0)"
		// 	}
		// }

		// if c.extra != "" {
		// 	str += " " + c.extra
		// }
		return str
	}
}

func (*postgresDriver) columnDataType(dataType string, defaultValue ...any) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		str := dataType
		// if !column.GoNullable() {
		// 	str += " NOT NULL"
		// }
		// // If it's not primary key or foreign key
		// if !column.Key() && len(defaultValue) > 0 {
		// 	// PRIMARY KEY cannot have default value
		// 	str += " DEFAULT " + format(defaultValue[0])
		// }
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
