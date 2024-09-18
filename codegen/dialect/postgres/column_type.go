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
	return map[string]*dialect.ColumnType{
		"rune": {
			DataType: s.columnDataType("char(1)"),
			Valuer:   "string({{goPath}})",
			Scanner:  "{{addrOfGoPath}}",
		},
		"string": {
			DataType: s.columnDataType("varchar(255)"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"bool": {
			DataType: s.columnDataType("boolean", false),
			Valuer:   "bool({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Bool({{addrOfGoPath}})",
		},
		"int8": {
			DataType: s.intDataType("smallint", int64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int16": {
			DataType: s.intDataType("smallint", int64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int32": {
			DataType: s.intDataType("integer", int64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int64": {
			DataType: s.columnDataType("bigint", int64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"int": {
			DataType: s.intDataType("integer", int64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint8": {
			DataType: s.intDataType("smallint", uint64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint16": {
			DataType: s.intDataType("smallint", uint64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint32": {
			DataType: s.intDataType("integer", uint64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint64": {
			DataType: s.intDataType("bigint", uint64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"uint": {
			DataType: s.intDataType("integer", uint64(0)),
			Valuer:   "int64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Integer({{addrOfGoPath}})",
		},
		"float32": {
			DataType: s.columnDataType("real"),
			Valuer:   "float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"float64": {
			DataType: s.columnDataType("double precision"),
			Valuer:   "float64({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Float({{addrOfGoPath}})",
		},
		"time.Time": {
			DataType: s.columnDataType("timestamp(6) with time zone"),
			Valuer:   "(time.Time)({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.Time({{addrOfGoPath}})",
		},
		"*time.Time": {
			DataType: s.columnDataType("timestamp(6) with time zone"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.Time({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.PtrOfTime({{addrOfGoPath}})",
		},
		"[...]string": {
			DataType: s.columnDataType("text[{{len}}]"),
			Valuer:   "github.com/lib/pq.StringArray({{goPath}}[:])",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}}[:])",
		},
		"[...]bool": {
			DataType: s.columnDataType("bool[{{len}}]"),
			Valuer:   "github.com/lib/pq.BoolArray({{goPath}}[:])",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}}[:])",
		},
		"[]string": {
			DataType: s.columnDataType("text[]"),
			Valuer:   "github.com/lib/pq.StringArray({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]byte": {
			DataType: s.columnDataType("blob"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"[]rune": {
			DataType: s.columnDataType("blob"),
			Valuer:   "string({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.String({{addrOfGoPath}})",
		},
		"[]bool": {
			DataType: s.columnDataType("bool[]"),
			Valuer:   "github.com/lib/pq.BoolArray({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[][]byte": {
			DataType: s.columnDataType("bytea"),
			Valuer:   "github.com/lib/pq.ByteaArray({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]float32": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "github.com/lib/pq.Float32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]float64": {
			DataType: s.columnDataType("double[]"),
			Valuer:   "github.com/lib/pq.Float64Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]int": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]int8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]int16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]int32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]int64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "github.com/lib/pq.Int64Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]uint": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]uint8": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]uint16": {
			DataType: s.columnDataType("int2[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]uint32": {
			DataType: s.columnDataType("int4[]"),
			Valuer:   "github.com/lib/pq.Int32Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array({{addrOfGoPath}})",
		},
		"[]uint64": {
			DataType: s.columnDataType("int8[]"),
			Valuer:   "github.com/lib/pq.Int64Array({{goPath}})",
			Scanner:  "github.com/lib/pq.Array{{addrOfGoPath}})",
		},
		"*": {
			DataType: s.columnDataType("json"),
			Valuer:   "github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler({{goPath}})",
			Scanner:  "github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler({{addrOfGoPath}})",
		},
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
		if !column.Nullable() {
			str += " NOT NULL"
		}
		// Auto increment cannot has default value
		if !column.AutoIncr() && len(defaultValue) > 0 {
			str += " DEFAULT " + format(defaultValue[0])
			switch any(defaultValue[0]).(type) {
			case uint64:
				str += " CHECK (" + s.QuoteIdentifier(column.Name()) + " >= 0)"
			}
		}

		// if c.extra != "" {
		// 	str += " " + c.extra
		// }
		return str
	}
}

func (*postgresDriver) columnDataType(dataType string, defaultValue ...any) func(dialect.GoColumn) string {
	return func(column dialect.GoColumn) string {
		str := dataType
		if !column.Nullable() {
			str += " NOT NULL"
		}
		if len(defaultValue) > 0 {
			// PRIMARY KEY cannot have default value
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