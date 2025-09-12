package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/si3nloong/sqlgen/cmd/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/internal/goutil"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/sqltype"
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
			Scanner:  goutil.GenericFuncName(encoding.StringScanner[string, *string], "{{elemType}}", "{{addr}}"),
		},
		"bool": {
			DataType: s.columnDataType("BOOL", false),
			Valuer:   "(bool)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.BoolScanner[bool, *bool], "{{elemType}}", "{{addr}}"),
		},
		"int": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.IntScanner[int, *int], "{{elemType}}", "{{addr}}"),
		},
		"int8": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int8Scanner[int8, *int8], "{{elemType}}", "{{addr}}"),
		},
		"int16": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int16Scanner[int16, *int16], "{{elemType}}", "{{addr}}"),
		},
		"int32": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int32Scanner[int32, *int32], "{{elemType}}", "{{addr}}"),
		},
		"int64": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Int64Scanner[int64, *int64], "{{elemType}}", "{{addr}}"),
		},
		"uint": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.UintScanner[uint, *uint], "{{elemType}}", "{{addr}}"),
		},
		"uint8": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint8Scanner[uint8, *uint8], "{{elemType}}", "{{addr}}"),
		},
		"uint16": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint16Scanner[uint16, *uint16], "{{elemType}}", "{{addr}}"),
		},
		"uint32": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint32Scanner[uint32, *uint32], "{{elemType}}", "{{addr}}"),
		},
		"uint64": {
			DataType: s.columnDataType("INTEGER", false),
			Valuer:   "(int64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Uint64Scanner[uint64, *uint64], "{{elemType}}", "{{addr}}"),
		},
		"float32": {
			DataType: s.columnDataType("FLOAT", int64(0)),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Float32Scanner[float32, *float32], "{{elemType}}", "{{addr}}"),
		},
		"float64": {
			DataType: s.columnDataType("FLOAT", int64(0)),
			Valuer:   "(float64)({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.Float64Scanner[float64, *float64], "{{elemType}}", "{{addr}}"),
		},
		"time.Time": {
			DataType: s.columnDataType("TIMESTAMP", sql.RawBytes("CURRENT_TIMESTAMP")),
			Valuer:   "(time.Time)({{goPath}})",
			Scanner:  goutil.GenericFunc(encoding.TimeScanner[*time.Time], "{{addr}}"),
		},
		"sql.RawBytes": {
			DataType: s.columnDataType("TEXT"),
			Valuer:   "database/sql.RawBytes({{goPath}})",
			Scanner:  "(*database/sql.RawBytes)({{addrOfGoPath}})",
		},
		"encoding/json.RawMessage": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "(string)({{goPath}})",
			Scanner:  "({{addrOfGoPath}})",
		},
		"[...]rune": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("VARCHAR(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: goutil.GenericFunc(encoding.RuneArrayScanner[rune], "{{goPath}}[:]", "{{len}}"),
		},
		"[...]byte": {
			DataType: func(c dialect.GoColumn) string {
				return s.columnDataType(fmt.Sprintf("CHAR(%d)", c.Size()))(c)
			},
			Valuer:  "string({{goPath}}[:])",
			Scanner: goutil.GenericFunc(encoding.ByteArrayScanner[byte], "{{goPath}}[:]", "{{len}}"),
		},
		"[]string": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.StringSlice[string]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.StringSlice[string]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]byte": {
			DataType: s.columnDataType("BLOB"),
			Valuer:   "string({{goPath}})",
			Scanner:  goutil.GenericFuncName(encoding.StringScanner[[]byte, *[]byte], "{{baseType}}", "{{addr}}"),
		},
		"[]bool": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.BoolSlice[bool]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.BoolSlice[bool]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.IntSlice[int]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.IntSlice[int]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Int8Slice[int8]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Int8Slice[int8]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Int16Slice[int16]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Int16Slice[int16]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Int32Slice[int32]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Int32Slice[int32]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Int64Slice[int64]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Int64Slice[int64]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.UintSlice[uint]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.UintSlice[uint]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Uint8Slice[uint8]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Uint8Slice[uint8]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Uint16Slice[uint16]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Uint16Slice[uint16]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Uint32Slice[uint32]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Uint32Slice[uint32]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GetTypeName(sqltype.Uint64Slice[uint64]{}) + "[{{elemType}}]({{goPath}})",
			Scanner:  "(*" + goutil.GetTypeName(sqltype.Uint64Slice[uint64]{}) + "[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/sqltype.Float32Slice[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/sqltype.Float32Slice[{{elemType}}])({{addrOfGoPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("JSON"),
			Valuer:   "(github.com/si3nloong/sqlgen/sequel/sqltype.Float64Slice[{{elemType}}])({{goPath}})",
			Scanner:  "(*github.com/si3nloong/sqlgen/sequel/sqltype.Float64Slice[{{elemType}}])({{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("JSON"),
			Valuer:   goutil.GenericFunc(encoding.JSONValue[any], "{{goPath}}"),
			Scanner:  goutil.GenericFunc(encoding.JSONScanner[any], "{{addr}}"),
		},
	}
	return dataTypes
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
