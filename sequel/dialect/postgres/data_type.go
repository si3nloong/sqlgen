package postgres

import (
	"database/sql"
	"fmt"
	"go/types"
	"strconv"
	"strings"
	"unsafe"

	"github.com/si3nloong/sqlgen/sequel"
)

func dataType(f sequel.GoColumnSchema) (dataType string) {
	var (
		ptrs = make([]types.Type, 0)
		t    = f.Type()
		prev types.Type
	)

	for t != nil {
		switch v := t.(type) {
		case *types.Pointer:
			prev = v.Elem()
			ptrs = append(ptrs, v)
		case *types.Basic:
			prev = t.Underlying()
		case *types.Named:
			prev = t.Underlying()
		default:
			break
		}

		switch t.String() {
		case "rune":
			return "CHAR(1)" + notNullDefault(ptrs)
		case "int8", "int16":
			return "INT2" + notNullDefault(ptrs, 0)
		case "int32", "int":
			return "INT" + notNullDefault(ptrs, 0)
		case "int64":
			return "INT8" + notNullDefault(ptrs, 0)
		case "bool":
			return "BOOL" + notNullDefault(ptrs, false)
		case "uint8", "uint16", "byte":
			return "INT2" + notNullDefault(ptrs, 0) + " CHECK(" + f.ColumnName() + " >= 0)"
		case "uint32", "uint":
			return "INT" + notNullDefault(ptrs, 0) + " CHECK(" + f.ColumnName() + " >= 0)"
		case "uint64":
			return "INT8" + notNullDefault(ptrs, 0) + " CHECK(" + f.ColumnName() + " >= 0)"
		case "float32":
			return "DOUBLE PRECISION" + notNullDefault(ptrs, 0.0)
		case "float64":
			return "DOUBLE PRECISION" + notNullDefault(ptrs, 0.0)
		case "cloud.google.com/go/civil.Time":
			return "TIME" + notNullDefault(ptrs, sql.RawBytes(`CURRENT_TIME`))
		case "cloud.google.com/go/civil.Date":
			return "DATE" + notNullDefault(ptrs, sql.RawBytes(`CURRENT_DATE`))
		case "cloud.google.com/go/civil.DateTime":
			if size := f.Size(); size > 0 {
				return fmt.Sprintf("TIMESTAMP(%d)", size) + notNullDefault(ptrs, sql.RawBytes(`NOW()`))
			}
			return "TIMESTAMP" + notNullDefault(ptrs, sql.RawBytes(`NOW()`))
		case "time.Time":
			if size := f.Size(); size > 0 {
				return fmt.Sprintf("TIMESTAMP(%d) WITH TIME ZONE", size) + notNullDefault(ptrs, sql.RawBytes(`NOW()`))
			}
			return "TIMESTAMP WITH TIME ZONE" + notNullDefault(ptrs, sql.RawBytes(`NOW()`))
		case "string":
			size := int64(255)
			if v := f.Size(); v > 0 {
				size = v
			}
			return fmt.Sprintf("VARCHAR(%d)", size) + notNullDefault(ptrs, "")
		case "[]rune":
			return "VARCHAR(255)" + notNullDefault(ptrs)
		case "[]byte":
			return "BYTEA" + notNullDefault(ptrs)
		case "[16]byte":
			// if f.IsBinary {
			// 	return "BIT(16)"
			// }
			return "VARBIT(36)"
		case "encoding/json.RawMessage":
			return "VARBIT" + notNullDefault(ptrs)
		default:
			if strings.HasPrefix(t.String(), "[]") {
				// if f.IsBinary {
				// 	return "JSONB" + notNullDefault(ptrs)
				// }
				return "JSON" + notNullDefault(ptrs)
			}
		}
		if prev == t {
			break
		}
		t = prev
	}
	return "VARCHAR(255)" + notNullDefault(ptrs)
}

func notNullDefault(ptrs []types.Type, defaultValue ...any) string {
	if len(ptrs) > 0 {
		return ""
	}
	if len(defaultValue) > 0 {
		return " NOT NULL DEFAULT " + format(defaultValue[0])
	}
	return " NOT NULL"
}

func format(v any) string {
	switch vi := v.(type) {
	case string:
		return "'" + vi + "'"
	case bool:
		return strconv.FormatBool(vi)
	case int:
		return strconv.Itoa(vi)
	case float32:
		return strconv.FormatFloat(float64(vi), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(vi, 'f', -1, 64)
	case sql.RawBytes:
		return unsafe.String(unsafe.SliceData(vi), len(vi))
	default:
		panic("unsupported data type")
	}
}
