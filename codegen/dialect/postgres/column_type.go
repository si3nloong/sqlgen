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
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int8({{addrOfGoPath}})",
		},
		"int16": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int16({{addrOfGoPath}})",
		},
		"int32": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int32({{addrOfGoPath}})",
		},
		"int64": {
			DataType: s.columnDataType("bigint", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int64({{addrOfGoPath}})",
		},
		"int": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint8": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint8({{addrOfGoPath}})",
		},
		"uint16": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint16({{addrOfGoPath}})",
		},
		"uint32": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint32({{addrOfGoPath}})",
		},
		"uint64": {
			DataType: s.intDataType("bigint", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint64({{addrOfGoPath}})",
		},
		"uint": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint({{addrOfGoPath}})",
		},
		"float32": {
			DataType: s.columnDataType("real"),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float32({{addrOfGoPath}})",
		},
		"float64": {
			DataType: s.columnDataType("double precision"),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float64({{addrOfGoPath}})",
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
		"*int8": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int8Scanner({{addr}})",
		},
		"*int16": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int16Scanner({{addr}})",
		},
		"*int32": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int32Scanner({{addr}})",
		},
		"*int64": {
			DataType: s.columnDataType("bigint", int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Int64Scanner({{addr}})",
		},
		"*int": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.IntScanner({{addr}})",
		},
		"*uint8": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint8Scanner({{addr}})",
		},
		"*uint16": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint16Scanner({{addr}})",
		},
		"*uint32": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint32Scanner({{addr}})",
		},
		"*uint64": {
			DataType: s.intDataType("bigint", uint64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Uint64Scanner({{addr}})",
		},
		"*uint": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Integer({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.UintScanner({{addr}})",
		},
		"*float32": {
			DataType: s.columnDataType("real"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float32({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float32({{addrOfGoPath}})",
		},
		"*float64": {
			DataType: s.columnDataType("double precision"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Float64({{addrOfGoPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float64({{addrOfGoPath}})",
		},
		"*time.Time": {
			DataType: s.columnDataType("timestamptz(6)"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Time({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Time({{addr}})",
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
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.StringArray[{{elemType}}]({{goPath}}[:])",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.StringArray[{{elemType}}]({{goPath}}[:])",
		},
		"[...]bool": {
			DataType: s.columnDataType("bool[{{len}}]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.BoolArray[{{elemType}}]({{goPath}}[:])",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.BoolArray[{{elemType}}])({{addrOfGoPath}}[:])",
		},
		"[]string": {
			DataType: s.columnDataType("text[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.StringArray[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.StringArray[{{elemType}}])({{addrOfGoPath}})",
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
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.BoolArray[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.BoolArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[][]byte": {
			DataType: s.columnDataType("bytea"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.ByteArray({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.ByteArray({{goPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Float32Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Float32Array[{{elemType}}]({{goPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Float64Array[{{elemType}}]({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Float64Array[{{elemType}}]({{goPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.IntArray[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.IntArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Int8Array[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.IntArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Int16Array[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.IntArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Int32Array[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.IntArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Int64Array[{{elemType}}]({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.IntArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.UintArray[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.UintArray[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.UintArray[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint8Array[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint16Array[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint16Array[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint32Array[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint32Array[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint64Array[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.Uint64Array[{{elemType}}])({{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("json"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler({{addrOfGoPath}})",
		},
	}
	return dataTypes
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
