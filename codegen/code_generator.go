package codegen

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen"
	"github.com/si3nloong/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type Generator struct {
	*bytes.Buffer
	config             *Config
	dialect            dialect.Dialect
	quoteRune          rune
	staticVar          bool
	columnTypes        map[string]*dialect.ColumnType
	defaultColumnTypes map[string]*dialect.ColumnType
	errs               []error
}

func newGenerator(cfg *Config, d dialect.Dialect) *Generator {
	gen := new(Generator)
	gen.Buffer = new(bytes.Buffer)
	gen.config = cfg
	switch d.QuoteRune() {
	case '"':
		gen.quoteRune = '`'
	case '`':
		gen.quoteRune = '"'
	default:
		gen.quoteRune = '"'
	}
	gen.dialect = d
	gen.staticVar = d.QuoteVar(1) == d.QuoteVar(0)
	gen.defaultColumnTypes = d.ColumnDataTypes()
	gen.columnTypes = make(map[string]*dialect.ColumnType)
	for k, decl := range cfg.DataTypes {
		gen.columnTypes[k] = &dialect.ColumnType{
			DataType: func(c dialect.GoColumn) string {
				return decl.DataType
			},
			Scanner:    decl.Scanner,
			Valuer:     decl.Valuer,
			SQLScanner: decl.SQLScanner,
			SQLValuer:  decl.SQLValuer,
		}
	}
	return gen
}

func (g *Generator) LogError(err error) {
	// When it's strict mode, the program will stop once
	// it encounter any error
	if g.config.Strict != nil && *g.config.Strict {
		panic(err)
	}
	g.errs = append(g.errs, err)
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
	if !g.config.QuoteIdentifier {
		return str
	}
	return g.dialect.QuoteIdentifier(str)
}

func (g *Generator) generate(pkg *packages.Package, dstDir string, typeInferred bool, schemas []*tableInfo) error {
	defer g.Reset()
	importPkgs := NewPackage(pkg.PkgPath, pkg.Name)
	importPkgs.Import(types.NewPackage("strings", ""))
	importPkgs.Import(types.NewPackage("strconv", ""))
	importPkgs.Import(types.NewPackage("database/sql/driver", ""))
	importPkgs.Import(types.NewPackage("github.com/si3nloong/sqlgen/sequel", ""))

	for len(schemas) > 0 {
		t := schemas[0]

		g.L()

		if method, wrongType := t.Implements(sqlDatabaser); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "DatabaseName" but wrong footprint`, t.goName))
		} else if method != nil && !wrongType && t.dbName != "" {
			g.L("func (" + t.goName + ") DatabaseName() string {")
			g.L(`return ` + g.Quote(t.dbName))
			g.L("}")
		}

		// Build the "TableName" function which return the table name
		if method, wrongType := t.Implements(sqlTabler); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.goName))
		} else if method != nil && !wrongType {
			g.L("func (" + t.goName + ") TableName() string {")
			g.L(`return ` + g.Quote(t.tableName))
			g.L("}")
		} else {
			// TODO: we need to do something when table name is declare by user
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

		// Build the "Columns" function which return the column names
		if method, wrongType := t.Implements(sqlColumner); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Columns" but wrong footprint`, t.goName))
		} else if method != nil && !wrongType {
			g.L("func (" + t.goName + ") Columns() []string {")
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

		// Build the "SQLColumns" function which return the column SQL query
		if method, wrongType := t.Implements(sqlQueryColumner); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "SQLColumns" but wrong footprint`, t.goName))
		} else if method != nil && !wrongType {
			if _, ok := lo.Find(t.columns, func(v *columnInfo) bool {
				_, exists := v.sqlScanner()
				return exists
			}); ok {
				g.L("func (" + t.goName + ") SQLColumns() []string {")
				g.WriteString("return []string{")
				for i, f := range t.columns {
					if i > 0 {
						g.WriteByte(',')
					}
					g.WriteString(g.Quote(g.sqlScanner(f)))
				}
				g.WriteString("}\n")
				g.L("}")
			}
		}

		// Build the "Values" function which return the column values
		if method, wrongType := t.Implements(sqlValuer); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Values" but wrong footprint`, t.goName))
		} else if method != nil && !wrongType {
			g.buildValuer(importPkgs, t)
		}

		// Build the "Addrs" function which return the column addressable values
		if method, wrongType := t.PtrImplements(sqlScanner); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Addrs" but wrong footprint`, t.goName))
		} else if method != nil && !wrongType {
			g.buildScanner(importPkgs, t)
		}

		if !t.hasNoColsExceptAutoPK() {
			if g.staticVar {
				g.L("func (" + t.goName + ") InsertPlaceholders(row int) string {")
				g.L(`return "(` + strings.Repeat(","+g.dialect.QuoteVar(0), len(t.colsWithoutAutoIncrPK()))[1:] + `)"`)
				g.L("}")
			} else {
				cols := t.colsWithoutAutoIncrPK()
				g.L("func (" + t.goName + ") InsertPlaceholders(row int) string {")
				g.L(fmt.Sprintf("const noOfColumn = %d", len(cols)))
				g.WriteString(`return "("+`)
				for i, f := range cols {
					if i > 0 {
						g.WriteString(`+","+`)
					}
					if sqlValuer, ok := f.sqlValuer(); ok {
						matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
						g.WriteString(fmt.Sprintf("%q + strconv.Itoa((row * noOfColumn) + %d) +%q", matches[1]+string(g.dialect.VarRune()), i+1, matches[5]))
					} else {
						g.WriteString(fmt.Sprintf(`%q+ strconv.Itoa((row * noOfColumn) + %d)`, string(g.dialect.VarRune()), i+1))
					}
				}
				g.WriteString(`+")"` + "\n")
				g.L("}")
			}
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

				if sqlValuer, ok := f.sqlValuer(); ok {
					matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
					if len(matches) > 4 {
						g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.SQLColumnValuer[", typeStr, "] {")
						g.L(`return sequel.SQLColumn`+specificType+`(`, g.Quote(f.ColumnName()), `, v.`, f.GoPath()+",", fmt.Sprintf(`func(placeholder string) string { return %q+ placeholder + %q}`, matches[1]+matches[2], matches[4]+matches[5]), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.Type()), ` })`)
						g.L("}")
					} else {
						g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
						g.L(`return sequel.Column`, specificType, `(`, g.Quote(f.ColumnName()), `, v.`, f.GoPath(), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.Type()), ` })`)
						g.L("}")
					}
				} else {
					g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
					g.L("return sequel.Column", specificType, "(", g.Quote(f.ColumnName()), ", v.", f.GoPath(), ", func(val ", typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.Type()), ` })`)
					g.L("}")
				}
			}
		}

		schemas = schemas[1:]
	}

	rmb := g.Buffer
	g.Buffer = new(bytes.Buffer)
	g.buildHeader()
	g.L("package " + pkg.Name)
	g.L()

	if len(importPkgs.imports) > 0 {
		g.L("import (")
		for _, pkg := range importPkgs.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				g.L(strconv.Quote(pkg.Path()))
			} else {
				// If the import is alias import path
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
	// panic("stop now!")

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

	// Copy the buffer to the output file
	if _, err := io.Copy(f, g); err != nil {
		return err
	}
	return f.Close()
}

func (g *Generator) buildHeader() {
	if !g.config.SkipHeader {
		g.L(fmt.Sprintf("// Code generated by sqlgen, version %s; DO NOT EDIT.", sqlgen.Version))
		g.L()
	}
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
		g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
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
		g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
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
		g.WriteString(g.scanner(importPkgs, "&v."+f.GoPath(), f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) buildFindByPK(importPkgs *Package, t *tableInfo) {
	buf := strpool.AcquireString()
	buf.WriteString("SELECT ")
	for i, f := range t.columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(g.sqlScanner(f))
	}
	buf.WriteString(" FROM ")
	var query string
	if method, wrongType := t.Implements(sqlTabler); wrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.goName))
	} else if method != nil {
		buf.WriteString(g.QuoteIdentifier(t.tableName))
	} else {
		query = g.Quote(buf.String()) + "+ v.TableName() +"
		buf.Reset()
	}
	buf.WriteString(" WHERE ")
	if len(t.keys) > 1 {
		// Composite primary key
		keyCols := lo.Map(t.keys, func(v *columnInfo, _ int) string {
			return v.ColumnName()
		})
		buf.WriteString("(" + strings.Join(keyCols, ",") + ")" + " = ")
		buf.WriteByte('(')
		for i, k := range t.keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.sqlValuer(k, i+1))
		}
		buf.WriteByte(')')
	} else {
		for i, f := range t.keys {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(f.ColumnName() + " = " + g.dialect.QuoteVar(i+1))
		}
	}
	buf.WriteString(" LIMIT 1;")
	g.L("func (v " + t.goName + ") FindOneByPKStmt() (string, []any) {")
	g.WriteString("return " + query + g.Quote(buf.String()) + ", []any{")
	strpool.ReleaseString(buf)
	for i, f := range t.keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) buildInsertOne(importPkgs *Package, t *tableInfo) {
	// Filter out auto increment key
	columns := lo.Filter(t.columns, func(col *columnInfo, _ int) bool {
		return col != t.autoIncrKey
	})
	var query string
	buf := strpool.AcquireString()
	if method, wrongType := t.Implements(sqlTabler); wrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.goName))
	} else if method != nil {
		buf.WriteString("INSERT INTO " + t.tableName)
	} else {
		query = g.Quote("INSERT INTO ") + "+ v.TableName() +"
	}
	buf.WriteString(" (")
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
		buf.WriteString(g.sqlValuer(f, i))
	}
	buf.WriteByte(')')
	if g.config.Driver == Postgres {
		buf.WriteString(" RETURNING ")
		for i, f := range t.columns {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.sqlScanner(f))
		}
	}
	buf.WriteByte(';')
	// If the columns and after filter columns is the same
	// mean it has no auto increment key
	g.L("func (v " + t.goName + ") InsertOneStmt() (string, []any) {")
	if len(columns) == len(t.columns) {
		g.L("return " + query + g.Quote(buf.String()) + ", v.Values()")
		strpool.ReleaseString(buf)
	} else {
		g.WriteString("return " + query + g.Quote(buf.String()) + ", []any{")
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

func (g *Generator) buildUpdateByPK(importPkgs *Package, t *tableInfo) {
	buf := strpool.AcquireString()
	var query string
	if method, wrongType := t.Implements(sqlTabler); wrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.goName))
	} else if method != nil {
		buf.WriteString("UPDATE " + g.QuoteIdentifier(t.tableName))
	} else {
		query = g.Quote("UPDATE ") + "+ v.TableName() +"
	}
	buf.WriteString(" SET ")
	columns := lo.Filter(t.columns, func(col *columnInfo, _ int) bool {
		return !lo.Contains(t.keys, col)
	})
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(f.ColumnName() + " = " + g.sqlValuer(f, i))
	}
	buf.WriteString(" WHERE ")
	if len(t.keys) > 1 {
		keyCols := lo.Map(t.keys, func(v *columnInfo, _ int) string {
			return v.ColumnName()
		})
		buf.WriteString("(" + strings.Join(keyCols, ",") + ")" + " = ")
		buf.WriteByte('(')
		for i, k := range t.keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.sqlValuer(k, i+len(columns)))
		}
		buf.WriteByte(')')
	} else {
		for i, k := range t.keys {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(k.ColumnName() + " = " + g.sqlValuer(k, i+len(columns)))
		}
	}
	buf.WriteByte(';')
	g.L("func (v " + t.goName + ") UpdateOneByPKStmt() (string, []any) {")
	g.WriteString("return " + query + g.Quote(buf.String()) + ", []any{")
	strpool.ReleaseString(buf)
	for i, f := range append(columns, t.keys...) {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.valuer(importPkgs, "v."+f.GoPath(), f.Type()))
	}
	g.WriteString("}\n")
	g.L("}")
	strpool.ReleaseString(buf)
}

func (g *Generator) valuer(importPkgs *Package, goPath string, t types.Type) string {
	if columnType, ok := g.columnTypes[t.String()]; ok {
		return Expr(columnType.Valuer).Format(importPkgs, ExprParams{GoPath: goPath})
	} else if _, wrong := types.MissingMethod(t, goSqlValuer, true); wrong {
		return Expr("(database/sql/driver.Valuer)({{goPath}})").Format(importPkgs, ExprParams{GoPath: goPath})
	} else if columnType, ok := g.columnDataType(t); ok && columnType.Valuer != "" {
		return Expr(columnType.Valuer).Format(importPkgs, ExprParams{GoPath: goPath, Len: arraySize(t)})
	} else if isImplemented(t, textMarshaler) {
		return Expr("github.com/si3nloong/sqlgen/sequel/types.TextMarshaler({{goPath}})").Format(importPkgs, ExprParams{GoPath: goPath})
	}
	return Expr(g.defaultColumnTypes["*"].Valuer).Format(importPkgs, ExprParams{GoPath: goPath})
}

func (g *Generator) scanner(importPkgs *Package, goPath string, t types.Type) string {
	if columnType, ok := g.columnTypes[t.String()]; ok && columnType.Scanner != "" {
		return Expr(columnType.Scanner).Format(importPkgs, ExprParams{GoPath: goPath})
	} else if types.Implements(newPointer(t), goSqlScanner) {
		return Expr("(database/sql.Scanner)({{addrOfGoPath}})").Format(importPkgs, ExprParams{GoPath: goPath})
	} else if columnType, ok := g.columnDataType(t); ok && columnType.Scanner != "" {
		return Expr(columnType.Scanner).Format(importPkgs, ExprParams{GoPath: goPath, Len: arraySize(t)})
	} else if isImplemented(newPointer(t), textUnmarshaler) {
		return Expr("github.com/si3nloong/sqlgen/sequel/types.TextUnmarshaler({{addrOfGoPath}})").Format(importPkgs, ExprParams{GoPath: goPath})
	}
	return Expr(g.defaultColumnTypes["*"].Scanner).Format(importPkgs, ExprParams{GoPath: goPath})
}

func (g *Generator) sqlScanner(f *columnInfo) string {
	if sqlScanner, ok := f.sqlScanner(); ok {
		matches := sqlFuncRegexp.FindStringSubmatch(sqlScanner("{}"))
		if len(matches) > 4 {
			return matches[1] + matches[2] + f.ColumnName() + matches[4] + matches[5]
		} else {
			return f.ColumnName()
		}
	}
	return f.ColumnName()
}

func (g *Generator) sqlValuer(f *columnInfo, idx int) string {
	if sqlValuer, ok := f.sqlValuer(); ok {
		matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
		// g.WriteString(fmt.Sprintf("%q + strconv.Itoa((row * noOfColumn) + %d) +%q", matches[1]+string(g.dialect.VarRune()), f.ColumnPos(), matches[5]))
		if len(matches) > 4 {
			return matches[1] + matches[2] + g.dialect.QuoteVar(idx+1) + matches[4] + matches[5]
		}
		return g.dialect.QuoteVar(idx + 1)
	}
	return g.dialect.QuoteVar(idx + 1)
}
