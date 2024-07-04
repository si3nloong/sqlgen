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

func dataType(f sequel.GoColumnSchema) *columnDefinition {
	var (
		ptrs = make([]types.Type, 0)
		col  = new(columnDefinition)
		t    = f.Type()
		prev types.Type
	)

	col.name = f.ColumnName()
	defer func() {
		if !col.nullable {
			col.nullable = len(ptrs) > 0
		}
	}()

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
			col.dataType = "CHAR"
			col.length = 1
			return col
		case "bool":
			col.dataType = "BOOL"
			col.defaultValue = false
			return col
		case "int8":
			col.dataType = "TINYINT"
			col.defaultValue = int64(0)
			return col
		case "int16":
			col.dataType = "SMALLINT"
			col.defaultValue = int64(0)
			return col
		case "int32":
			col.dataType = "MEDIUMINT"
			col.defaultValue = int64(0)
			return col
		case "int64":
			col.dataType = "BIGINT"
			col.defaultValue = int64(0)
			return col
		case "int":
			col.dataType = "INTEGER"
			col.defaultValue = int64(0)
			return col
		case "uint8":
			col.dataType = "TINYINT UNSIGNED"
			col.defaultValue = uint64(0)
			return col
		case "uint16":
			col.dataType = "SMALLINT UNSIGNED"
			col.defaultValue = uint64(0)
			return col
		case "uint32":
			col.dataType = "MEDIUMINT UNSIGNED"
			col.defaultValue = uint64(0)
			return col
		case "uint64":
			col.dataType = "BIGINT UNSIGNED"
			col.defaultValue = uint64(0)
			return col
		case "uint":
			col.dataType = "INTEGER UNSIGNED"
			col.defaultValue = uint64(0)
			return col
		case "float32":
			col.dataType = "FLOAT"
			col.defaultValue = float64(0.0)
			return col
		case "float64":
			col.dataType = "FLOAT"
			col.defaultValue = float64(0.0)
			return col
		case "cloud.google.com/go/civil.Time":
			col.dataType = "TIME"
			return col
		case "cloud.google.com/go/civil.Date":
			col.dataType = "DATE"
			return col
		case "cloud.google.com/go/civil.DateTime":
			col.dataType = "DATETIME"
			if size := f.Size(); size > 0 {
				col.length = size
				col.defaultValue = sql.RawBytes(fmt.Sprintf("CURRENT_TIMESTAMP(%d)", size))
				return col
			}
			col.defaultValue = sql.RawBytes(`CURRENT_TIMESTAMP`)
			return col
		case "time.Time":
			col.dataType = "TIMESTAMP"
			if size := f.Size(); size > 0 {
				col.length = size
				col.defaultValue = sql.RawBytes(fmt.Sprintf("CURRENT_TIMESTAMP(%d)", size))
				return col
			}
			col.defaultValue = sql.RawBytes(`CURRENT_TIMESTAMP`)
			return col
		case "string":
			col.dataType = "VARCHAR"
			col.defaultValue = ""
			col.length = int64(255)
			if size := f.Size(); size > 0 {
				col.length = size
			}
			return col
		case "[]byte":
			col.dataType = "BLOB"
			return col
		case "[16]byte":
			// if f.IsBinary {
			// 	return "BINARY(16)"
			// }
			col.dataType = "VARCHAR"
			col.defaultValue = sql.RawBytes(`UUID()`)
			col.length = 36
			return col
		case "encoding/json.RawMessage":
			col.dataType = "JSON"
			return col
		case "database/sql.NullBool":
			col.dataType = "BOOL"
			col.defaultValue = false
			col.nullable = true
			return col
		case "database/sql.NullString":
			col.dataType = "VARCHAR"
			col.defaultValue = ""
			col.length = 255
			col.nullable = true
			return col
		case "database/sql.NullInt16":
			col.dataType = "SMALLINT"
			col.defaultValue = int64(0)
			col.nullable = true
			return col
		case "database/sql.NullInt32":
			col.dataType = "MEDIUMINT"
			col.defaultValue = int64(0)
			col.nullable = true
			return col
		case "database/sql.NullInt64":
			col.dataType = "BIGINT"
			col.defaultValue = int64(0)
			col.nullable = true
			return col
		case "database/sql.NullFloat64":
			col.dataType = "FLOAT"
			col.defaultValue = float64(0.0)
			col.nullable = true
			return col
		case "database/sql.NullTime":
			col.dataType = "TIMESTAMP"
			col.nullable = true
			col.defaultValue = sql.RawBytes(`CURRENT_TIMESTAMP`)
			return col
		default:
			switch {
			case strings.HasPrefix(t.String(), "[]"):
				col.dataType = "JSON"
				return col
			case strings.HasPrefix(t.String(), "map"):
				col.dataType = "JSON"
				return col
			}
		}
		if prev == t {
			break
		}
		t = prev
	}
	col.dataType = "VARCHAR"
	if size := f.Size(); size == 0 {
		col.length = 255
	}
	return col
}

func format(v any) string {
	switch vi := v.(type) {
	case string:
		return "'" + vi + "'"
	case bool:
		return strconv.FormatBool(vi)
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
