package sqlite

import (
	"fmt"
	"go/types"
	"strings"

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
			return "CHAR(1)" + notNull(len(ptrs) > 0)
		case "bool", "int8":
			return "TINYINT" + notNull(len(ptrs) > 0)
		case "int16":
			return "SMALLINT" + notNull(len(ptrs) > 0)
		case "int32":
			return "MEDIUMINT" + notNull(len(ptrs) > 0)
		case "int64":
			return "BIGINT" + notNull(len(ptrs) > 0)
		case "int", "uint":
			return "INTEGER" + notNull(len(ptrs) > 0)
		case "uint8":
			return "TINYINT UNSIGNED" + notNull(len(ptrs) > 0)
		case "uint16":
			return "SMALLINT UNSIGNED" + notNull(len(ptrs) > 0)
		case "uint32":
			return "MEDIUMINT UNSIGNED" + notNull(len(ptrs) > 0)
		case "uint64":
			return "BIGINT UNSIGNED" + notNull(len(ptrs) > 0)
		case "float32":
			return "FLOAT" + notNull(len(ptrs) > 0)
		case "float64":
			return "FLOAT" + notNull(len(ptrs) > 0)
		case "cloud.google.com/go/civil.Date":
			return "DATE"
		case "time.Time":
			if size := f.Size(); size > 0 {
				return fmt.Sprintf("DATETIME(%d)", size) + notNull(len(ptrs) > 0)
			}
			return "DATETIME" + notNull(len(ptrs) > 0)
		case "string":
			size := 255
			if v := f.Size(); v > 0 {
				size = v
			}
			return fmt.Sprintf("VARCHAR(%d)", size) + notNull(len(ptrs) > 0)
		case "[]byte":
			return "BLOB" + notNull(len(ptrs) > 0)
		case "[16]byte":
			// if f.IsBinary {
			// 	return "BINARY(16)"
			// }
			return "VARCHAR(36)"
		case "encoding/json.RawMessage":
			return "JSON" + notNull(len(ptrs) > 0)
		default:
			switch {
			case strings.HasPrefix(t.String(), "[]"):
				return "JSON" + notNull(len(ptrs) > 0)
			case strings.HasPrefix(t.String(), "map"):
				return "JSON" + notNull(len(ptrs) > 0)
			}
		}
		if prev == t {
			break
		}
		t = prev
	}
	return "VARCHAR(255)" + notNull(len(ptrs) > 0)
}

func notNull(isNull bool) string {
	if isNull {
		return ""
	}
	return " NOT NULL"
}
