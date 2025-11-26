package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/goutil"
	"github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype"
	"github.com/si3nloong/sqlgen/sequel/encoding"
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
			Valuer:   "{{goPath}}",
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
			Scanner: goutil.GenericFuncName(encoding.StringScanner[string, *string], "{{elemType}}", "{{addr}}"),
		},
		"bool": {
			DataType: s.columnDataType("bool", false),
			Valuer:   "(bool)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.BoolScanner[bool, *bool], "{{elemType}}", "{{addr}}"),
		},
		"int": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.IntScanner[int, *int], "{{elemType}}", "{{addr}}"),
		},
		"int8": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int8Scanner[int8, *int8], "{{elemType}}", "{{addr}}"),
		},
		"int16": {
			DataType: s.intDataType("int2", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int16Scanner[int16, *int16], "{{elemType}}", "{{addr}}"),
		},
		"int32": {
			DataType: s.intDataType("int4", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int32Scanner[int32, *int32], "{{elemType}}", "{{addr}}"),
		},
		"int64": {
			DataType: s.columnDataType("bigint", int64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int64Scanner[int64, *int64], "{{elemType}}", "{{addr}}"),
		},
		"uint": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.UintScanner[uint, *uint], "{{elemType}}", "{{addr}}"),
		},
		"uint8": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint8Scanner[uint8, *uint8], "{{elemType}}", "{{addr}}"),
		},
		"uint16": {
			DataType: s.intDataType("int2", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint16Scanner[uint16, *uint16], "{{elemType}}", "{{addr}}"),
		},
		"uint32": {
			DataType: s.intDataType("int4", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint32Scanner[uint32, *uint32], "{{elemType}}", "{{addr}}"),
		},
		"uint64": {
			DataType: s.intDataType("bigint", uint64(0)),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint64Scanner[uint64, *uint64], "{{elemType}}", "{{addr}}"),
		},
		"float32": {
			DataType: s.columnDataType("real"),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Float32Scanner[float32, *float32], "{{elemType}}", "{{addr}}"),
		},
		"float64": {
			DataType: s.columnDataType("double precision"),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Float64Scanner[float64, *float64], "{{elemType}}", "{{addr}}"),
		},
		"time.Time": {
			DataType: s.columnDataType("timestamptz(6)", sql.RawBytes(`NOW()`)),
			Valuer:   "(time.Time)({{goPath}})",
			Scanner:  goutil.GenericFunc(encoding.TimeScanner[*time.Time], "{{addr}}"),
		},
		"sql.RawBytes": {
			DataType: s.columnDataType("text"),
			Valuer:   "database/sql.RawBytes({{goPath}})",
			Scanner:  "(*database/sql.RawBytes)({{addrOfGoPath}})",
		},
		"encoding/json.RawMessage": {
			DataType: s.columnDataType("json"),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "({{addrOfGoPath}})",
		},
		"[...]rune": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("varchar(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: goutil.GenericFunc(encoding.RuneArrayScanner[rune], "{{goPath}}[:]", "{{len}}"),
		},
		"[...]byte": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("varchar(%d)", c.Size()))(c)
			},
			Valuer:  "{{goPath}}[:]",
			Scanner: goutil.GenericFunc(encoding.ByteArrayScanner[byte], "{{goPath}}[:]", "{{len}}"),
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
			Scanner:  "(*" + goutil.GetTypeName(pgtype.StringArray[string]{}) + "[{{elemType}}])({{addrOfGoPath}})",
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
			Valuer:   "(" + goutil.GetTypeName(pgtype.BoolArray[bool]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.BoolArray[bool]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[][]byte": {
			DataType: s.columnDataType("bytea"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.ByteArray({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgtype.ByteArray({{goPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.IntArray[int]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.IntArray[int]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Int8Array[int8]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Int8Array[int8]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Int16Array[int16]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Int16Array[int16]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Int32Array[int32]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Int32Array[int32]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Int64Array[int64]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Int64Array[int64]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.UintArray[uint]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.UintArray[uint]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Uint8Array[uint8]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Uint8Array[uint8]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Uint16Array[uint16]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Uint16Array[uint16]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Uint32Array[uint32]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Uint32Array[uint32]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Uint64Array[uint64]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Uint64Array[uint64]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Float32Array[float32]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Float32Array[float32]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "(" + goutil.GetTypeName(pgtype.Float64Array[float64]{}) + "[{{elemType}}])({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(pgtype.Float64Array[float64]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("json"),
			Valuer:   goutil.GenericFunc(encoding.JSONValue[any], "{{goPath}}"),
			Scanner:  goutil.GenericFunc(encoding.JSONScanner[any], "{{addr}}"),
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
