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
		if f.CustomMarshaler != "" {
			return Expr(f.CustomMarshaler+"(%s)").Format(impPkgs, name)
		} else if f.IsBinary {
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

func addrOf(impPkgs *Package) func(string, *templates.Field) string {
	return func(n string, f *templates.Field) string {
		name := "&" + n + "." + f.GoPath
		if f.CustomUnmarshaler != "" {
			return Expr(f.CustomUnmarshaler+"(%s)").Format(impPkgs, name)
		} else if f.IsBinary {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.BinaryUnmarshaler(%s)").Format(impPkgs, name)
		} else if types.Implements(newPointer(f.Type), sqlScanner) {
			return Expr("(database/sql.Scanner)(%s)").Format(impPkgs, name)
		} else if typ, ok := UnderlyingType(f.Type); ok {
			return typ.Decoder.Format(impPkgs, name)
		} else if f.IsTextUnmarshaler {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.TextUnmarshaler(%s)").Format(impPkgs, name)
		}
		return name
	}
}

func (g *Generator) createTableStmt() func(string, *templates.Model) string {
	return func(n string, model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)

		buf.WriteString(g.Quote("CREATE TABLE IF NOT EXISTS "))
		buf.WriteString("+ " + n + ".TableName() +" + g.QuoteStringBegin() + " (")
		for i, f := range model.Fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			dataType, isNull := inspectDataType(f)
			buf.WriteString(g.QuoteIdentifier(f.ColumnName) + " " + dataType)
			if !isNull {
				buf.WriteString(" NOT NULL")
			}
			if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
				buf.WriteString(" AUTO_INCREMENT")
			}
		}
		if model.PK != nil {
			buf.WriteString(",PRIMARY KEY (" + g.QuoteIdentifier(model.PK.Field.ColumnName) + ")")
		}
		buf.WriteString(");" + g.QuoteStringEnd())
		return buf.String()
	}
}

func (g *Generator) alterTableStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("ALTER TABLE " + g.QuoteIdentifier(model.TableName) + " ")
		for i, f := range model.Fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString("MODIFY ")
			dataType, isNull := inspectDataType(f)
			buf.WriteString(g.QuoteIdentifier(f.ColumnName) + " " + dataType)
			if !isNull {
				buf.WriteString(" NOT NULL")
			}
			if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
				buf.WriteString(" AUTO_INCREMENT")
			}
			if i > 0 {
				// buf.WriteString(" FIRST")
				buf.WriteString(" AFTER " + g.QuoteIdentifier(model.Fields[i-1].ColumnName))
			}
		}
		// if model.PK != nil {
		// 	buf.WriteString(",MODIFY PRIMARY KEY (" + model.PK.Field.ColumnName + ")")
		// }
		buf.WriteByte(';')
		return buf.String()
	}
}

func (g *Generator) insertOneStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("INSERT INTO " + g.QuoteIdentifier(model.TableName) + " (")
		var fields []*templates.Field
		if model.PK != nil && model.PK.IsAutoIncr {
			for i := range model.Fields {
				if model.Fields[i] != model.PK.Field {
					fields = append(fields, model.Fields[i])
				}
			}
		} else {
			fields = append(fields, model.Fields...)
		}
		for i := range fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.QuoteIdentifier(fields[i].ColumnName))
		}
		buf.WriteString(") VALUES (")
		for i := range fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.QuoteVar(i + 1))
		}
		buf.WriteString(");")
		return g.Quote(buf.String())
	}
}

func (g *Generator) findByPKStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("SELECT ")
		for i := range model.Fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.QuoteIdentifier(model.Fields[i].ColumnName))
		}
		buf.WriteString(" FROM " + g.QuoteIdentifier(model.TableName) + " WHERE ")
		buf.WriteString(g.QuoteIdentifier(model.PK.Field.ColumnName))
		buf.WriteString(" = " + g.QuoteVar(1) + " LIMIT 1;")
		return g.Quote(buf.String())
	}
}

func (g *Generator) updateByPKStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("UPDATE " + g.QuoteIdentifier(model.TableName) + " SET ")
		var pos int
		for i := range model.Fields {
			if model.PK != nil && model.PK.Field == model.Fields[i] {
				continue
			}
			if pos > 0 {
				buf.WriteByte(',')
			}
			pos++
			buf.WriteString(g.QuoteIdentifier(model.Fields[i].ColumnName) + " = " + g.QuoteVar(pos))
		}
		buf.WriteString(" WHERE ")
		buf.WriteString(g.QuoteIdentifier(model.PK.Field.ColumnName))
		buf.WriteString(" = " + g.QuoteVar(pos+1) + " LIMIT 1;")
		return g.Quote(buf.String())
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
			blr.WriteString(dialect.QuoteVar(i + 1))
		}
		blr.WriteByte(')')
		return blr.String()
	}
}

type FieldTypeValueResult struct {
	FuncName string
	Type     string
	Valuer   string
	Value    string
}

func getFieldTypeValue(impPkgs *Package, prefix string) func(*templates.Field) FieldTypeValueResult {
	return func(f *templates.Field) FieldTypeValueResult {
		typeStr := f.Type.String()
		if idx := strings.Index(typeStr, "."); idx > 0 {
			typeStr = Expr(typeStr).Format(impPkgs)
		}
		return FieldTypeValueResult{
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
