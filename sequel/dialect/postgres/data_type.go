package postgres

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/si3nloong/sqlgen/codegen/templates"
)

func dataType(f *templates.Field) (dataType string) {
	var (
		ptrs = make([]types.Type, 0)
		t    = f.Type
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
		case "int8", "int16":
			return "INT2" + notNull(len(ptrs) > 0)
		case "int32", "int":
			return "INT" + notNull(len(ptrs) > 0)
		case "int64":
			return "INT8" + notNull(len(ptrs) > 0)
		case "bool", "uint8", "uint16", "byte":
			return "INT2" + notNull(len(ptrs) > 0) + " CHECK(" + f.ColumnName + " > 0)"
		case "uint32", "uint":
			return "INT" + notNull(len(ptrs) > 0) + " CHECK(" + f.ColumnName + " > 0)"
		case "uint64":
			return "INT8" + notNull(len(ptrs) > 0) + " CHECK(" + f.ColumnName + " > 0)"
		case "float32":
			return "DOUBLE PRECISION" + notNull(len(ptrs) > 0)
		case "float64":
			return "DOUBLE PRECISION" + notNull(len(ptrs) > 0)
		case "cloud.google.com/go/civil.Date":
			return "DATE" + notNull(len(ptrs) > 0)
		case "time.Time":
			var size int
			if f.Size > 0 && f.Size < 7 {
				size = f.Size
			}
			if size > 0 {
				return fmt.Sprintf("TIMESTAMP(%d)", size) + notNull(len(ptrs) > 0)
			}
			return "TIMESTAMP" + notNull(len(ptrs) > 0)
		case "string":
			size := 255
			if f.Size > 0 {
				size = f.Size
			}
			return fmt.Sprintf("VARCHAR(%d)", size) + notNull(len(ptrs) > 0)
		case "[]rune":
			return "VARCHAR(255)" + notNull(len(ptrs) > 0)
		case "[]byte":
			return "BYTEA" + notNull(len(ptrs) > 0)
		case "[16]byte":
			if f.IsBinary {
				return "BIT(16)"
			}
			return "VARBIT(36)"
		case "encoding/json.RawMessage":
			return "VARBIT" + notNull(len(ptrs) > 0)
		default:
			if strings.HasPrefix(t.String(), "[]") {
				if f.IsBinary {
					return "JSONB" + notNull(len(ptrs) > 0)
				}
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
