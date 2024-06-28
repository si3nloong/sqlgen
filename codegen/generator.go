package codegen

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var (
	sqlFuncRegexp = regexp.MustCompile(`(?i)\s*(\w+\()(\w+\s*\,\s*)?(\{\})(\s*\,\s*\w+)?(\))\s*`)
)

type tableInfo struct {
	goName      string
	dbName      string
	tableName   string
	t           types.Type
	autoIncrKey sequel.ColumnSchema
	keys        []sequel.ColumnSchema
	columns     []sequel.ColumnSchema
}

// Mean table has only pk
func (b tableInfo) hasNoColsExceptPK() bool {
	return len(b.keys) == len(b.columns)
}

func (b tableInfo) GoName() string {
	return b.goName
}
func (b tableInfo) DatabaseName() string {
	return b.dbName
}
func (b tableInfo) TableName() string {
	return b.tableName
}
func (b tableInfo) AutoIncrKey() (sequel.ColumnSchema, bool) {
	return b.autoIncrKey, b.autoIncrKey != nil
}
func (b tableInfo) Keys() []sequel.ColumnSchema {
	return b.keys
}
func (b tableInfo) Columns() []sequel.ColumnSchema {
	return nil
}
func (b tableInfo) Implements(T *types.Interface) (wrongType bool) {
	types.MissingMethod(b.t, T, true)
	return
}

type columnInfo struct {
	goName  string
	goPath  string
	colName string
	colPos  int
	t       types.Type
	tag     tagOpts
	model   *config.Model
}

var (
	_ (sequel.ColumnSchema) = (*columnInfo)(nil)
)

func (i columnInfo) SQLValuer() sequel.SQLFunc {
	if i.model == nil {
		return nil
	}
	return func(placeholder string) string {
		return strings.Replace(i.model.SQLValuer, "{placeholder}", placeholder, 1)
	}
}

func (i columnInfo) SQLScanner() sequel.SQLFunc {
	if i.model == nil {
		return nil
	}
	return func(column string) string {
		return strings.Replace(i.model.SQLScanner, "{column}", column, 1)
	}
}

func (i columnInfo) GoName() string {
	return i.goName
}
func (i columnInfo) GoPath() string {
	return i.goPath
}
func (i columnInfo) Type() types.Type {
	return i.t
}
func (i columnInfo) ColumnName() string {
	return i.colName
}
func (i columnInfo) ColumnPos() int {
	return i.colPos
}
func (i *columnInfo) Implements(T *types.Interface) (wrongType bool) {
	types.MissingMethod(i.t, T, true)
	return
}

type Generator struct {
	*bytes.Buffer
	config    *config.Config
	dialect   sequel.Dialect
	quoteRune rune
	staticVar bool
}

func newGenerator(cfg *config.Config, dialect sequel.Dialect) *Generator {
	gen := new(Generator)
	gen.Buffer = new(bytes.Buffer)
	gen.config = cfg
	switch dialect.QuoteRune() {
	case '"':
		gen.quoteRune = '`'
	case '`':
		gen.quoteRune = '"'
	default:
		gen.quoteRune = '"'
	}
	gen.dialect = dialect
	gen.staticVar = dialect.QuoteVar(1) == dialect.QuoteVar(0)
	return gen
}

func (g *Generator) L(str ...any) error {
	for _, v := range str {
		g.WriteString(fmt.Sprintf("%v", v))
	}
	g.WriteByte('\n')
	return nil
}

func (g *Generator) Quote(str string) string {
	buf := make([]byte, 0, len(str))
	buf = append(buf, byte(g.quoteRune))
	for i := range str {
		buf = append(buf, byte(str[i]))
	}
	buf = append(buf, byte(g.quoteRune))
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

func (g *Generator) QuoteIdentifier(str string) string {
	if g.config.OmitQuoteIdentifier {
		return str
	}
	return g.dialect.QuoteIdentifier(str)
}

func (g *Generator) generate(pkg *packages.Package, dstDir string, typeInferred bool, schemas []*tableInfo) error {
	// rename := g.config.RenameFunc()

	importPkgs := NewPackage(pkg.PkgPath, pkg.Name)
	importPkgs.Import(types.NewPackage("strings", ""))
	importPkgs.Import(types.NewPackage("strconv", ""))
	importPkgs.Import(types.NewPackage("database/sql/driver", ""))
	importPkgs.Import(types.NewPackage("github.com/si3nloong/sqlgen/sequel", ""))

	for len(schemas) > 0 {
		t := schemas[0]

		if isImplemented(t.t, sqlDatabaser) {
			return fmt.Errorf(`sqlgen: struct %q has implements %q but wrong footprint`, sqlDatabaser, t.goName)
		} else {
			// g.L("func (" + t.goName + ") DatabaseName() string {")
			// g.L(`return ` + strconv.Quote(t.goName))
			// g.L("}")
		}
		// if t.Implements(sqlTabler) {
		// 	return fmt.Errorf(`sqlgen: struct %q has implements %q but wrong footprint`, sqlTabler, t.goName)
		// } else
		g.L()

		if !isImplemented(t.t, sqlTabler) {
			g.L("func (" + t.goName + ") TableName() string {")
			g.L(`return ` + g.Quote(t.tableName))
			g.L("}")
		}

		if len(t.keys) > 0 {
			g.L("func (" + t.goName + ") HasPK() {}")
			if t.autoIncrKey != nil {
				g.L("func (" + t.goName + ") IsAutoIncr() {}")
			}

			if len(t.keys) > 1 {
				g.buildCompositeKeys(importPkgs, t)
			} else {
				pk := t.keys[0]
				g.L("func (v " + t.goName + ") PK() (string, int, any) {")
				g.L(`return `, g.Quote(pk.ColumnName()), ", ", pk.ColumnPos(), ", ", g.valuer(importPkgs, "v."+pk.GoPath(), pk.Type()))
				g.L("}")
			}
		}

		if !isImplemented(t.t, sqlColumner) {
			g.L("func (" + t.goName + ") ColumnNames() []string {")
			g.WriteString("return []string{")
			for i, f := range t.columns {
				if i > 0 {
					g.WriteByte(',')
				}
				g.WriteString(g.Quote(f.ColumnName()))
			}
			g.WriteString("}\n")
			g.L("}")
		}
		if _, ok := lo.Find(t.columns, func(f sequel.ColumnSchema) bool {
			return f.SQLScanner() != nil
		}); ok {
			g.L("func (" + t.goName + ") SQLColumns() []string {")
			g.WriteString("return []string{")
			for i, f := range t.columns {
				if i > 0 {
					g.WriteByte(',')
				}
				if sqlScanner := f.SQLScanner(); sqlScanner != nil {
					matches := sqlFuncRegexp.FindStringSubmatch(sqlScanner("{}"))
					g.WriteString(g.Quote(matches[1] + matches[2] + f.ColumnName() + matches[4] + matches[5]))
				} else {
					g.WriteString(g.Quote(f.ColumnName()))
				}
			}
			g.WriteString("}\n")
			g.L("}")
		}
		if !isImplemented(t.t, sqlValuer) {
			g.buildValuer(importPkgs, t)
		}
		if !isImplemented(types.NewPointer(t.t), sqlScanner) {
			g.buildScanner(importPkgs, t)
		}

		if g.staticVar {
			g.L("func (" + t.goName + ") InsertPlaceholders(row int) string {")
			g.L(`return "(` + strings.Repeat(","+g.dialect.QuoteVar(0), len(t.columns))[1:] + `)"`)
			g.L("}")
		} else {
			g.L("func (" + t.goName + ") InsertPlaceholders(row int) string {")
			g.L(fmt.Sprintf("const noOfColumn = %d", len(t.columns)))
			g.WriteString(`return "("+`)
			for i, f := range t.columns {
				if i > 0 {
					g.WriteString(`+","+`)
				}
				if sqlValuer := f.SQLValuer(); sqlValuer != nil {
					matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
					g.WriteString(fmt.Sprintf("%q + strconv.Itoa((row * noOfColumn) + %d) +%q", matches[1]+string(g.dialect.VarRune()), f.ColumnPos(), matches[5]))
				} else {
					g.WriteString(fmt.Sprintf(`%q+ strconv.Itoa((row * noOfColumn) + %d)`, string(g.dialect.VarRune()), f.ColumnPos()))
				}
			}
			g.WriteString(`+")"` + "\n")
			g.L("}")
		}

		g.buildInsertOne(importPkgs, t)
		if len(t.keys) > 0 {
			g.buildFindByPK(importPkgs, t)
			if !t.hasNoColsExceptPK() {
				g.buildUpdateByPK(importPkgs, t)
			}
		}

		if !g.config.OmitGetters {
			for _, f := range t.columns {
				typeStr := f.Type().String()
				if idx := strings.Index(typeStr, "."); idx > 0 {
					typeStr = Expr(typeStr).Format(importPkgs)
				}
				var specificType string
				if !typeInferred {
					specificType = "[" + typeStr + "]"
				}
				if sqlValuer := f.SQLValuer(); sqlValuer != nil {
					matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
					// g.WriteString(fmt.Sprintf("%q + strconv.Itoa((row * noOfColumn) + %d) +%q", matches[1]+string(g.dialect.VarRune()), f.ColumnPos(), matches[5]))
					g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.SQLColumnValuer[", typeStr, "] {")
					g.L(`return sequel.SQLColumn`+specificType+`(`, g.Quote(f.ColumnName()), `, v.`, f.GoPath()+",", fmt.Sprintf(`func(placeholder string) string { return %q+ placeholder + %q}`, matches[1]+matches[2], matches[4]+matches[5]), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.Type()), ` })`)
					g.L("}")
				} else {
					g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
					g.L(`return sequel.Column`, specificType, `(`, g.Quote(f.ColumnName()), `, v.`, f.GoPath(), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.Type()), ` })`)
					g.L("}")
				}
			}
		}

		schemas = schemas[1:]
	}

	rmb := g.Buffer
	g.Buffer = new(bytes.Buffer)
	if !g.config.SkipHeader {
		g.L(fmt.Sprintf("// Code generated by sqlgen, version %s; DO NOT EDIT.", sqlgen.Version))
		g.L()
	}
	g.L("package " + pkg.Name)
	g.L()
	if len(importPkgs.imports) > 0 {
		g.L("import (")
		for _, pkg := range importPkgs.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				g.L(strconv.Quote(pkg.Path()))
			} else {
				g.L(pkg.Name() + " " + strconv.Quote(pkg.Path()))
			}
		}
		g.L(")")
	}

	g.Write(rmb.Bytes())

	formatted, err := imports.Process("", g.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}
	g.Reset()
	g.Write(formatted)
	// fset := token.NewFileSet()
	// fileAST, err := parser.ParseFile(fset, "", g.Bytes(), parser.ParseComments|parser.AllErrors)
	// if err != nil {
	// 	return err
	// }

	// ast.SortImports(fset, fileAST)
	// g.Reset()
	// if err := format.Node(g, fset, fileAST); err != nil {
	// 	return err
	// }
	// err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(g, fset, fileAST)

	os.MkdirAll(dstDir, fileMode)
	fileDest := filepath.Join(dstDir, g.config.Exec.Filename)
	f, err := os.OpenFile(fileDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, g); err != nil {
		return err
	}
	return f.Close()
}

func (g *Generator) buildCompositeKeys(importPkgs *Package, table *tableInfo) {
	g.L("func (v " + table.goName + ") CompositeKey() ([]string, []int, []any) {")
	g.WriteString("return []string{")
	for i, f := range table.keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.Quote(f.ColumnName()))
	}
	g.WriteString("}, []int{")
	for i, f := range table.keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(strconv.Itoa(f.ColumnPos()))
	}
	g.WriteString("}, []any{")
	for i, f := range table.keys {
		if i > 0 {
			g.WriteByte(',')
		}

		goPath := "v." + f.GoPath()
		g.WriteString(g.valuer(importPkgs, goPath, f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) buildValuer(importPkgs *Package, table *tableInfo) {
	g.L("func (v " + table.goName + ") Values() []any {")
	g.WriteString("return []any{")
	for i, f := range table.columns {
		if i > 0 {
			g.WriteByte(',')
		}

		goPath := "v." + f.GoPath()
		g.WriteString(g.valuer(importPkgs, goPath, f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) buildScanner(importPkgs *Package, table *tableInfo) {
	g.L("func (v *" + table.goName + ") Addrs() []any {")
	g.WriteString("return []any{")
	for i, f := range table.columns {
		if i > 0 {
			g.WriteByte(',')
		}

		goPath := "&v." + f.GoPath()
		g.WriteString(g.scanner(importPkgs, goPath, f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) valuer(importPkgs *Package, goPath string, t types.Type) string {

	if model, ok := g.config.Models[t.String()]; ok {
		// TODO: do it better
		return Expr(strings.Replace(model.Valuer, "{field}", goPath, -1)).Format(importPkgs)
	} else if _, wrong := types.MissingMethod(t, goSqlValuer, true); wrong {
		return Expr("(database/sql/driver.Valuer)(%s)").Format(importPkgs, goPath)
	} else if codec, _ := UnderlyingType(t); codec != nil {
		return codec.Encoder.Format(importPkgs, goPath)
	} else if isImplemented(t, textMarshaler) {
		return Expr("github.com/si3nloong/sqlgen/sequel/types.TextMarshaler(%s)").Format(importPkgs, goPath)
	} else if isImplemented(t, binaryMarshaler) {
		return Expr("github.com/si3nloong/sqlgen/sequel/types.BinaryMarshaler(%s)").Format(importPkgs, goPath)
	} else {
		return Expr("github.com/si3nloong/sqlgen/sequel/types.JSONMarshaler(%s)").Format(importPkgs, goPath)
	}
}

func (g *Generator) scanner(importPkgs *Package, goPath string, t types.Type) string {
	if model, ok := g.config.Models[t.String()]; ok {
		// TODO: do it better
		return Expr(strings.Replace(model.Scanner, "{field}", goPath, -1)).Format(importPkgs)
	} else if types.Implements(newPointer(t), goSqlScanner) {
		return Expr("(database/sql.Scanner)(%s)").Format(importPkgs, goPath)
	} else if codec, _ := UnderlyingType(t); codec != nil {
		return codec.Decoder.Format(importPkgs, goPath)
	} else if isImplemented(types.NewPointer(t), textMarshaler) {
		return Expr("github.com/si3nloong/sqlgen/sequel/types.TextUnmarshaler(%s)").Format(importPkgs, goPath)
	}
	return Expr("github.com/si3nloong/sqlgen/sequel/types.JSONUnmarshaler(%s)").Format(importPkgs, goPath)
}

func (g *Generator) buildFindByPK(importPkgs *Package, table *tableInfo) {
	buf := strpool.AcquireString()
	buf.WriteString("SELECT ")
	for i, f := range table.columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(f.ColumnName())
	}
	buf.WriteString(" FROM " + table.tableName + " WHERE ")
	for i, f := range table.keys {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(f.ColumnName() + " = " + g.dialect.QuoteVar(i+1))
	}
	buf.WriteString(" LIMIT 1;")
	g.L("func (v " + table.goName + ") FindByPKStmt() (string, []any) {")
	g.WriteString("return " + g.Quote(buf.String()) + ", []any{")
	strpool.ReleaseString(buf)
	for i, f := range table.keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) buildInsertOne(importPkgs *Package, table *tableInfo) {
	// Filter out auto increment key
	columns := lo.Filter(table.columns, func(v sequel.ColumnSchema, _ int) bool {
		return v != table.autoIncrKey
	})
	buf := strpool.AcquireString()
	buf.WriteString("INSERT INTO " + table.tableName + " (")
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(f.ColumnName())
	}
	buf.WriteString(") VALUES (")
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		if sqlValuer := f.SQLValuer(); sqlValuer != nil {
			matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
			// g.WriteString(fmt.Sprintf("%q + strconv.Itoa((row * noOfColumn) + %d) +%q", matches[1]+string(g.dialect.VarRune()), f.ColumnPos(), matches[5]))
			buf.WriteString(matches[1] + matches[2] + g.dialect.QuoteVar(f.ColumnPos()) + matches[4] + matches[5])
		} else {
			buf.WriteString(g.dialect.QuoteVar(f.ColumnPos()))
		}
	}
	buf.WriteByte(')')
	if g.config.Driver == config.Postgres {
		buf.WriteString(" RETURNING ")
		for i, f := range table.columns {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(f.ColumnName())
		}
	}
	buf.WriteByte(';')
	// If the columns and after filter columns is the same
	// mean it has no auto increment key
	g.L("func (v " + table.goName + ") InsertOneStmt() (string, []any) {")
	if len(columns) == len(table.columns) {
		g.L("return " + g.Quote(buf.String()) + ", v.Values()")
		strpool.ReleaseString(buf)
	} else {
		g.WriteString("return " + g.Quote(buf.String()) + ", []any{")
		strpool.ReleaseString(buf)
		for i, f := range columns {
			if i > 0 {
				g.WriteByte(',')
			}
			g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
		}
		g.WriteString("}\n")
	}
	g.L("}")
}

func (g *Generator) buildUpdateByPK(importPkgs *Package, table *tableInfo) {
	buf := strpool.AcquireString()
	buf.WriteString("UPDATE " + (table.tableName) + " SET ")
	columns := lo.Filter(table.columns, func(v sequel.ColumnSchema, _ int) bool {
		return !lo.Contains(table.keys, v)
	})
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(f.ColumnName() + " = " + g.dialect.QuoteVar(i+1))
	}
	buf.WriteString(" WHERE ")
	for i, k := range table.keys {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(k.ColumnName() + " = " + g.dialect.QuoteVar(i+len(columns)+1))
	}
	// if g.config.Driver == config.Postgres {
	// 	buf.WriteString(" RETURNING ")
	// 	for i, f := range table.columns {
	// 		if i > 0 {
	// 			buf.WriteByte(',')
	// 		}
	// 		buf.WriteString(f.ColumnName())
	// 	}
	// }
	buf.WriteString(" LIMIT 1;")
	g.L("func (v " + table.goName + ") UpdateByPKStmt() (string, []any) {")
	g.WriteString("return " + g.Quote(buf.String()) + ", []any{")
	strpool.ReleaseString(buf)
	for i, f := range append(columns, table.keys...) {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
	strpool.ReleaseString(buf)
}
