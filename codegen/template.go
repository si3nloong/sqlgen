package codegen

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func reserveImport(impPkgs *Package) func(pkgPath string, aliases ...string) string {
	return func(pkgPath string, aliases ...string) string {
		name := filepath.Base(pkgPath)
		if len(aliases) > 0 {
			name = aliases[0]
		}
		impPkgs.Import(types.NewPackage(pkgPath, name))
		return ""
	}
}

func castAs(impPkgs *Package) func(*templates.Field, ...string) string {
	return func(f *templates.Field, n ...string) string {
		var name string
		if len(n) > 0 {
			name = n[0]
		} else {
			name = "v." + f.GoPath
		}
		if f.IsBinary {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.BinaryMarshaler(%s)").Format(impPkgs, name)
		} else if _, wrong := types.MissingMethod(f.Type, sqlValuer, true); wrong {
			return Expr("(database/sql/driver.Valuer)(%s)").Format(impPkgs, name)
		} else if typ, ok := UnderlyingType(f.Type); ok {
			return typ.Encoder.Format(impPkgs, name)
		} else if f.IsTextMarshaler {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.TextMarshaler(%s)").Format(impPkgs, name)
		}
		return name
	}
}

func addrOf(impPkgs *Package) func(n string, f *templates.Field) string {
	return func(n string, f *templates.Field) string {
		v := "&" + n + "." + f.GoPath
		if f.IsBinary {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.BinaryUnmarshaler(%s)").Format(impPkgs, v)
		} else if types.Implements(newPointer(f.Type), sqlScanner) {
			return Expr("(database/sql.Scanner)(%s)").Format(impPkgs, v)
		} else if typ, ok := UnderlyingType(f.Type); ok {
			return typ.Decoder.Format(impPkgs, v)
		} else if f.IsTextUnmarshaler {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.TextUnmarshaler(%s)").Format(impPkgs, v)
		}
		return v
	}
}

func createTableStmt(dialect sequel.Dialect) func(string, *templates.Model) string {
	return func(n string, model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)

		buf.WriteString(`"CREATE TABLE IF NOT EXISTS "+ `)
		buf.WriteString(n + `.TableName() +" (`)
		for i, f := range model.Fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			dataType, isNull := inspectDataType(f)
			buf.WriteString(dialect.Wrap(f.ColumnName) + " " + dataType)
			if !isNull {
				buf.WriteString(" NOT NULL")
			}
			if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
				buf.WriteString(" AUTO_INCREMENT")
			}
		}
		if model.PK != nil {
			buf.WriteString(",PRIMARY KEY (" + dialect.Wrap(model.PK.Field.ColumnName) + ")")
		}
		buf.WriteString(`);"`)
		return buf.String()
	}
}

func alterTableStmt(dialect sequel.Dialect) func(model *templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("ALTER TABLE " + dialect.Wrap(model.TableName) + " ")
		for i, f := range model.Fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString("MODIFY ")
			dataType, isNull := inspectDataType(f)
			buf.WriteString(dialect.Wrap(f.ColumnName) + " " + dataType)
			if !isNull {
				buf.WriteString(" NOT NULL")
			}
			if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
				buf.WriteString(" AUTO_INCREMENT")
			}
			if i > 0 {
				// buf.WriteString(" FIRST")
				buf.WriteString(" AFTER " + dialect.Wrap(model.Fields[i-1].ColumnName))
			}
		}
		// if model.PK != nil {
		// 	buf.WriteString(",MODIFY PRIMARY KEY (" + model.PK.Field.ColumnName + ")")
		// }
		buf.WriteByte(';')
		return buf.String()
	}
}

func varStmt(dialect sequel.Dialect) func(*templates.Model) string {
	return func(model *templates.Model) string {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		noOfCols := len(model.Fields)
		if model.PK != nil && model.PK.IsAutoIncr {
			noOfCols--
		}
		blr.WriteByte('(')
		for i := 0; i < noOfCols; i++ {
			if i > 0 {
				blr.WriteByte(',')
			}
			blr.WriteString(dialect.Var(i))
		}
		blr.WriteByte(')')
		return blr.String()
	}
}

type FieldTypeValue struct {
	FuncName string
	Type     string
	Valuer   string
	Value    string
}

func getFieldTypeValue(impPkgs *Package, prefix string) func(*templates.Field) FieldTypeValue {
	return func(f *templates.Field) FieldTypeValue {
		typeStr := f.Type.String()
		if idx := strings.Index(typeStr, "."); idx > 0 {
			typeStr = Expr(typeStr).Format(impPkgs)
		}
		return FieldTypeValue{
			FuncName: prefix + f.GoName,
			Type:     typeStr,
			Valuer:   castAs(impPkgs)(f),
			Value:    "v." + f.GoName,
		}
	}
}

func inspectDataType(f *templates.Field) (dataType string, null bool) {
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
			return "CHAR(1)", len(ptrs) > 0
		case "bool":
			return "TINYINT", len(ptrs) > 0
		case "int8":
			return "TINYINT", len(ptrs) > 0
		case "int16":
			return "SMALLINT", len(ptrs) > 0
		case "int32":
			return "MEDIUMINT", len(ptrs) > 0
		case "int64":
			return "BIGINT", len(ptrs) > 0
		case "int":
			return "INTEGER", len(ptrs) > 0
		case "uint8":
			return "TINYINT UNSIGNED", len(ptrs) > 0
		case "uint16":
			return "SMALLINT UNSIGNED", len(ptrs) > 0
		case "uint32":
			return "MEDIUMINT UNSIGNED", len(ptrs) > 0
		case "uint64":
			return "BIGINT UNSIGNED", len(ptrs) > 0
		case "uint":
			return "INTEGER UNSIGNED", len(ptrs) > 0
		case "float32":
			return "FLOAT", len(ptrs) > 0
		case "float64":
			return "FLOAT", len(ptrs) > 0
		case "cloud.google.com/go/civil.Date":
			return "DATE", len(ptrs) > 0
		case "time.Time":
			var size int
			if f.Size > 0 && f.Size < 7 {
				size = f.Size
			}
			if size > 0 {
				return fmt.Sprintf("DATETIME(%d)", size), len(ptrs) > 0
			}
			return "DATETIME", len(ptrs) > 0
		case "string":
			size := 255
			if f.Size > 0 {
				size = f.Size
			}
			return fmt.Sprintf("VARCHAR(%d)", size), len(ptrs) > 0
		case "[]byte":
			return "BLOB", true
		case "[16]byte":
			if f.IsBinary {
				return "BINARY(16)", len(ptrs) > 0
			}
			return "VARCHAR(36)", len(ptrs) > 0
		default:
			if strings.HasPrefix(t.String(), "[]") {
				return "JSON", len(ptrs) > 0
			}
		}
		if prev == t {
			break
		}
		t = prev
	}
	return "VARCHAR(255)", len(ptrs) > 0
}

func newPointer(t types.Type) *types.Pointer {
	v, ok := t.(*types.Pointer)
	if ok {
		return v
	}
	return types.NewPointer(t)
}
