package codegen

import (
	"go/types"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/templates"
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
		} else {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler(%s)").Format(impPkgs, name)
		}
		// return name
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
		} else {
			return Expr("github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler(%s)").Format(impPkgs, name)
		}
	}
}

func (g *Generator) insertOneStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("INSERT INTO " + g.QuoteIdentifier(model.TableName) + " (")
		fields := make([]*templates.Field, 0)
		if model.IsAutoIncr {
			fields = lo.Filter(model.Fields, func(f *templates.Field, _ int) bool {
				return f != model.Keys[0]
			})
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
		for i, k := range model.Keys {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(g.QuoteIdentifier(k.ColumnName) + " = " + g.QuoteVar(i+1))
		}
		buf.WriteString(" LIMIT 1;")
		return g.Quote(buf.String())
	}
}

func (g *Generator) updateByPKStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		buf.WriteString("UPDATE " + g.QuoteIdentifier(model.TableName) + " SET ")
		var pos int
		for _, f := range model.Fields {
			if lo.Contains(model.Keys, f) {
				continue
			}
			if pos > 0 {
				buf.WriteByte(',')
			}
			pos++
			buf.WriteString(g.QuoteIdentifier(f.ColumnName) + " = " + g.QuoteVar(pos))
		}
		buf.WriteString(" WHERE ")
		for i, k := range model.Keys {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(g.QuoteIdentifier(k.ColumnName) + " = " + g.QuoteVar(pos+i+1))
		}
		buf.WriteString(" LIMIT 1;")
		return g.Quote(buf.String())
	}
}

func (g *Generator) varStmt() func(*templates.Model) string {
	return func(model *templates.Model) string {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		noOfCols := len(model.Fields)
		if len(model.Keys) > 0 && model.IsAutoIncr {
			noOfCols--
		}
		blr.WriteByte('(')
		for i := 0; i < noOfCols; i++ {
			if i > 0 {
				blr.WriteByte(',')
			}
			blr.WriteString(g.QuoteVar(i + 1))
		}
		blr.WriteByte(')')
		return blr.String()
	}
}

func (g *Generator) varRune() string {
	return string(g.dialect.VarRune())
}

func (g *Generator) typeConstraint(typeInferred bool) func(FieldTypeValueResult) string {
	return func(result FieldTypeValueResult) string {
		if typeInferred {
			return ""
		}
		return "[" + result.Type + "]"
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
