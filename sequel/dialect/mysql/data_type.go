package mysql

import (
	"database/sql"
	"fmt"
	"go/types"
	"strconv"
	"strings"
	"unsafe"

	"github.com/si3nloong/sqlgen/sequel"
)

func dataType(f sequel.ColumnSchema) (dataType string) {
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
		case "bool":
			return "BOOL" + notNullDefault(ptrs, false)
		case "int8":
			return "TINYINT" + notNullDefault(ptrs, 0)
		case "int16":
			return "SMALLINT" + notNullDefault(ptrs, 0)
		case "int32":
			return "MEDIUMINT" + notNullDefault(ptrs, 0)
		case "int64":
			return "BIGINT" + notNullDefault(ptrs, 0)
		case "int", "uint":
			return "INTEGER" + notNullDefault(ptrs, 0)
		case "uint8":
			return "TINYINT UNSIGNED" + notNullDefault(ptrs, 0)
		case "uint16":
			return "SMALLINT UNSIGNED" + notNullDefault(ptrs, 0)
		case "uint32":
			return "MEDIUMINT UNSIGNED" + notNullDefault(ptrs, 0)
		case "uint64":
			return "BIGINT UNSIGNED" + notNullDefault(ptrs, 0)
		case "float32":
			return "FLOAT" + notNullDefault(ptrs, 0.0)
		case "float64":
			return "FLOAT" + notNullDefault(ptrs, 0.0)
		case "cloud.google.com/go/civil.Time":
			return "TIME" + notNullDefault(ptrs)
		case "cloud.google.com/go/civil.Date":
			return "DATE" + notNullDefault(ptrs)
		case "cloud.google.com/go/civil.DateTime":
			if size := f.Size(); size > 0 {
				return fmt.Sprintf("DATETIME(%d)", size) + notNullDefault(ptrs, sql.RawBytes(fmt.Sprintf("CURRENT_TIMESTAMP(%d)", size)))
			}
			return "DATETIME" + notNullDefault(ptrs, sql.RawBytes(`CURRENT_TIMESTAMP`))
		case "time.Time":
			if size := f.Size(); size > 0 {
				return fmt.Sprintf("TIMESTAMP(%d)", size) + notNullDefault(ptrs, sql.RawBytes(fmt.Sprintf("CURRENT_TIMESTAMP(%d)", size)))
			}
			return "TIMESTAMP" + notNullDefault(ptrs, sql.RawBytes(`CURRENT_TIMESTAMP`))
		case "string":
			size := 255
			if v := f.Size(); v > 0 {
				size = v
			}
			return fmt.Sprintf("VARCHAR(%d)", size) + notNullDefault(ptrs, "")
		case "[]byte":
			return "BLOB" + notNullDefault(ptrs)
		case "[16]byte":
			// if f.IsBinary {
			// 	return "BINARY(16)"
			// }
			return "VARCHAR(36)" + notNullDefault(ptrs, sql.RawBytes(`UUID()`))
		case "encoding/json.RawMessage":
			return "JSON" + notNullDefault(ptrs)
		default:
			switch {
			case strings.HasPrefix(t.String(), "[]"):
				return "JSON" + notNullDefault(ptrs)
			case strings.HasPrefix(t.String(), "map"):
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
