package codegen

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen"
	"github.com/si3nloong/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/internal/compiler"
	"github.com/si3nloong/sqlgen/internal/goutil"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/imports"
)

// The minimum go type we need to map
var goTypes = []string{
	"byte",
	"rune",
	"bool",
	"string",
	"float32", "float64",
	"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"time.Time",
	// "any", "sql.RawBytes", "json.RawMessage"
}

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

func newGenerator(cfg *Config, d dialect.Dialect) (*Generator, error) {
	gen := new(Generator)
	gen.Buffer = new(bytes.Buffer)
	gen.config = cfg
	switch d.QuoteRune() {
	case '"':
		gen.quoteRune = '`'
	case '`':
		gen.quoteRune = '"'
	default:
		return nil, fmt.Errorf(`sqlgen: invalid quote character %q for string`, d.QuoteRune())
	}
	gen.dialect = d
	gen.staticVar = d.QuoteVar(1) == d.QuoteVar(0)
	gen.defaultColumnTypes = d.ColumnDataTypes()
	// Check the dialect cover the basic go types
	for _, t := range goTypes {
		if _, ok := gen.defaultColumnTypes[t]; !ok {
			return nil, fmt.Errorf(`sqlgen: SQL dialect %q missing column type mapping for type %q`, d.Driver(), t)
		}
	}
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
	return gen, nil
}

func (g *Generator) LogError(err error) {
	// When it's strict mode, the program will stop once
	// it encounter any error
	if g.config.Strict != nil && *g.config.Strict {
		panic(err)
	}
	g.errs = append(g.errs, err)
}

func (g *Generator) L(values ...any) error {
	for i := range values {
		switch vi := values[i].(type) {
		case string:
			g.WriteString(vi)
		default:
			g.WriteString(fmt.Sprintf("%v", vi))
		}
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

// Generate model functions
func (g *Generator) generateModels(
	dstDir string,
	schema *compiler.Package,
) error {
	defer g.Reset()
	importPkgs := NewPackage(schema.Pkg.PkgPath, schema.Pkg.Name)
	importPkgs.Import(types.NewPackage("strings", ""))
	importPkgs.Import(types.NewPackage("strconv", ""))
	importPkgs.Import(types.NewPackage("database/sql/driver", ""))
	importPkgs.Import(types.NewPackage("github.com/si3nloong/sqlgen/sequel", ""))

	for len(schema.Tables) > 0 {
		t := schema.Tables[0]

		g.L()

		// if method, wrongType := t.Implements(sqlDatabaser); wrongType {
		// 	g.LogError(fmt.Errorf(`sqlgen: struct %q has function "DatabaseName" but wrong footprint`, t.Name))
		// } else if method != nil && !wrongType && t.dbName != "" {
		// 	g.L("func (" + t.GoName + ") DatabaseName() string {")
		// 	g.L(`return ` + g.Quote(g.QuoteIdentifier(t.dbName)))
		// 	g.L("}")
		// }
		var readonly bool
		if method, wrongType := t.PtrImplements(locker); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.Name))
		} else if method == nil && !wrongType {
			readonly = true
		}

		// Build the "TableName" function which return the table name
		if method, wrongType := t.Implements(sqlTabler); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.Name))
		} else if method != nil && !wrongType {
			g.L("func ("+t.GoName+") ", sqlTabler.Method(0).Name(), "() string {")
			g.L(`return ` + g.Quote(g.QuoteIdentifier(t.Name)))
			g.L("}")
		} else {
			// TODO: we need to do something when table name is declare by user
		}

		if t.HasPK() {
			g.L("func (" + t.GoName + ") HasPK() {}")
			pk, ok := t.AutoIncrKey()
			if ok {
				g.L("func (" + t.GoName + ") IsAutoIncr() {}")
				g.L("func (v *" + t.GoName + ") ScanAutoIncr(val int64) error {")
				g.L("v." + pk.GoName + " = " + pk.Type.String() + "(val)")
				g.L("return nil")
				g.L("}")
			} else if len(t.Keys) == 1 {
				pk = t.Keys[0]
			}
			if pk != nil {
				g.L("func (v " + t.GoName + ") PK() (string, int, any) {")
				g.L(`return `, g.Quote(g.QuoteIdentifier(pk.Name)), ", ", pk.Pos, ", ", g.getOrValue(importPkgs, "v", pk))
				g.L("}")
			} else {
				g.buildCompositeKeys(importPkgs, t)
			}
		}

		// Build the "SQLColumns" function which return the column SQL query
		if method, wrongType := t.Implements(sqlQueryColumner); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "SQLColumns" but wrong footprint`, t.Name))
		} else if method != nil && !wrongType {
			// if _, ok := lo.Find(t.columns, func(v *columnInfo) bool {
			// 	_, exists := v.sqlScanner()
			// 	return exists
			// }); ok {
			// 	g.L("func (" + t.goName + ") SQLColumns() []string {")
			// 	g.WriteString("return []string{")
			// 	for i, f := range t.columns {
			// 		if i > 0 {
			// 			g.WriteByte(',')
			// 		}
			// 		g.WriteString(g.Quote(g.sqlScanner(f)))
			// 	}
			// 	g.L("}")
			// 	g.L("}")
			// }
		}

		// Build the "Columns" function which return the column names
		if method, wrongType := t.Implements(sqlColumner); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Columns" but wrong footprint`, t.Name))
		} else if method != nil && !wrongType {
			g.L("func ("+t.GoName+") ", sqlColumner.Method(0).Name(), "() []string {")
			g.WriteString("return []string{")
			for i, col := range t.Columns {
				if i > 0 {
					g.WriteByte(',')
				}
				g.WriteString(g.Quote(g.QuoteIdentifier(col.Name)))
			}
			g.L(fmt.Sprintf("} // %d", len(t.Columns)))
			g.L("}")
		}

		// Build the "Values" function which return the column values
		if method, wrongType := t.Implements(sqlValuer); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Values" but wrong footprint`, t.Name))
		} else if method != nil && !wrongType {
			g.buildValuer(importPkgs, t)
		}

		// Build the "Addrs" function which return the column addressable values
		if method, wrongType := t.PtrImplements(sqlScanner); wrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Addrs" but wrong footprint`, t.Name))
		} else if method != nil && !wrongType {
			g.buildScanner(importPkgs, t)
		}

		if !readonly {
			insertColumns := t.InsertColumns()
			if len(insertColumns) > 0 && len(insertColumns) != len(t.Columns) {
				g.L("func (" + t.GoName + ") InsertColumns() []string {")
				g.WriteString("return []string{")
				for i, col := range insertColumns {
					if i > 0 {
						g.WriteByte(',')
					}
					g.WriteString(g.Quote(g.QuoteIdentifier(col.Name)))
				}
				g.L(fmt.Sprintf("} // %d", len(insertColumns)))
				g.L("}")
			}

			if len(insertColumns) > 0 {
				if g.staticVar {
					g.L("func (" + t.GoName + ") InsertPlaceholders(row int) string {")
					g.L(`return "(` + strings.Repeat(","+g.dialect.QuoteVar(0), len(insertColumns))[1:] + `)"` + fmt.Sprintf(" // %d", len(insertColumns)))
					g.L("}")
				} else {
					g.L("func (" + t.GoName + ") InsertPlaceholders(row int) string {")
					g.L(fmt.Sprintf("const noOfColumn = %d", len(insertColumns)))
					g.WriteString(`return "("+`)
					for i := range insertColumns {
						if i > 0 {
							g.WriteString(`+","+`)
						}
						g.WriteString(fmt.Sprintf(`%q+ strconv.Itoa((row * noOfColumn) + %d)`, string(g.dialect.VarRune()), i+1))
					}
					g.WriteString(`+")"`)
					g.L("}")
				}

				g.buildInsertOne(importPkgs, t)
			}
		}

		if t.HasPK() {
			g.buildFindByPK(importPkgs, t)
			if !readonly && len(t.ColumnsWithoutPK()) > 0 {
				g.buildUpdateByPK(importPkgs, t)
			}
		}

		// Build getter
		for _, f := range t.Columns {
			g.L("func (v "+t.GoName+") ", valueFunc(f), " driver.Value {")
			queue := []string{}
			// Find all the possible pointer paths
			ptrPaths := f.GoPtrPaths()
			for _, p := range ptrPaths {
				g.L("if v" + p.GoPath + " != nil {")
				queue = append(queue, "}")
			}

			if f.IsPtr() {
				// Deference the pointer value and return it
				g.L("return ", g.valuer(importPkgs, "*v"+f.GoPath, assertAsPtr[types.Pointer](f.Type).Elem()))
			} else {
				g.L("return " + g.valuer(importPkgs, "v"+f.GoPath, f.Type))
			}
			for len(queue) > 0 {
				g.L(queue[0])
				queue = queue[1:]
			}
			if len(ptrPaths) > 0 {
				g.L("return nil")
			}
			g.L("}")
			// 	// 		if idx := strings.Index(typeStr, "."); idx > 0 {
			// 	// 			typeStr = Expr(typeStr).Format(importPkgs)
			// 	// 		}
			// 	// 		var specificType string
			// 	// 		if !typeInferred {
			// 	// 			specificType = "[" + typeStr + "]"
			// 	// 		}

			// 	// 		if sqlValuer, ok := f.sqlValuer(); ok {
			// 	// 			matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
			// 	// 			if len(matches) > 4 {
			// 	// 				g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.SQLColumnValuer[", typeStr, "] {")
			// 	// 				g.L(`return sequel.SQLColumn`+specificType+`(`, g.Quote(g.QuoteIdentifier(f.ColumnName())), `, v.`, f.GoPath()+",", fmt.Sprintf(`func(placeholder string) string { return %q+ placeholder + %q}`, matches[1]+matches[2], matches[4]+matches[5]), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.t), ` })`)
			// 	// 				g.L("}")
			// 	// 			} else {
			// 	// 				g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
			// 	// 				g.L(`return sequel.Column`, specificType, `(`, g.Quote(g.QuoteIdentifier(f.ColumnName())), `, v.`, f.GoPath(), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.t), ` })`)
			// 	// 				g.L("}")
			// 	// 			}
			// 	// 		} else {
			// 	// 			g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
			// 	// 			g.L("return sequel.Column", specificType, "(", g.Quote(g.QuoteIdentifier(f.ColumnName())), ", v.", f.GoPath(), ", func(val ", typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.t), ` })`)
			// 	// 			g.L("}")
			// 	// 		}
			// }
		}
		for _, f := range t.Columns {
			typeStr := f.Type.String()
			if idx := strings.Index(typeStr, "."); idx > 0 {
				typeStr = Expr(typeStr).Format(importPkgs)
			}
			g.L("func (v "+t.GoName+") ", g.config.Getter.Prefix+f.GoName, "() sequel.ColumnValuer[", typeStr, "] {")
			g.L("return sequel.Column(", g.Quote(g.QuoteIdentifier(f.Name)), ", v"+f.GoPath, ", func(val ", typeStr, `) driver.Value {`)
			// g.L("return ", g.valuer(importPkgs, "val", f.Type))

			if f.IsPtr() {
				g.L("if val != nil {")
				// Deference the pointer value and return it
				g.L("return ", g.valuer(importPkgs, "*val", assertAsPtr[types.Pointer](f.Type).Elem()))
				g.L("}")
				g.L("return nil")
			} else {
				g.L("return " + g.valuer(importPkgs, "val", f.Type))
			}
			g.L("})")
			g.L("}")
		}

		schema.Tables = schema.Tables[1:]
	}

	rmb := g.Buffer
	g.Buffer = new(bytes.Buffer)
	g.buildHeader()
	g.L("package " + schema.Pkg.Name)
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

	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return err
	}
	fileDest := filepath.Join(dstDir, g.config.Exec.Filename)
	f, err := os.OpenFile(fileDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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

func (g *Generator) buildCompositeKeys(importPkgs *Package, table *compiler.Table) {
	g.L("func (v " + table.GoName + ") CompositeKey() ([]string, []int, []any) {")
	g.WriteString("return []string{")
	for i, f := range table.Keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.Quote(f.Name))
	}
	g.WriteString("}, []int{")
	for i, f := range table.Keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(strconv.Itoa(f.Pos))
	}
	g.WriteString("}, []any{")
	for i, k := range table.Keys {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.getOrValue(importPkgs, "v", k))
	}
	g.L("}")
	g.L("}")
}

func (g *Generator) buildValuer(importPkgs *Package, table *compiler.Table) {
	columns := table.InsertColumns()
	if len(columns) > 0 {
		g.L("func (v "+table.GoName+") ", sqlValuer.Method(0).Name(), "() []any {")
		g.WriteString("return []any{")
		tmpl := " // %" + strwidth(len(columns)) + "d - %s"
		for _, f := range columns {
			g.WriteString("\n" + g.getOrValue(importPkgs, "v", f) + ",")
			g.WriteString(fmt.Sprintf(tmpl, f.Pos, f.Name))
		}
		g.L("\n}")
		g.L("}")
	}
}

func (g *Generator) buildScanner(importPkgs *Package, table *compiler.Table) {
	g.L("func (v *"+table.GoName+") ", sqlScanner.Method(0).Name(), "() []any {")
	for _, f := range table.GoPtrPaths() {
		path := "v" + f.GoPath
		g.L("if " + path + " == nil {")
		g.L(path+" = new(", Expr(strings.TrimPrefix(f.Type.String(), "*")).Format(importPkgs, ExprParams{}), ")")
		g.L("}")
	}
	g.WriteString("return []any{")
	tmpl := " // %" + strwidth(len(table.Columns)) + "d - %s"
	for _, f := range table.Columns {
		g.WriteString("\n" + g.scanner(importPkgs, "&v"+f.GoPath, f.Type) + "," + fmt.Sprintf(tmpl, f.Pos, f.Name))
	}
	g.WriteString("\n}\n")
	g.L("}")
}

func (g *Generator) buildFindByPK(importPkgs *Package, t *compiler.Table) {
	buf := strpool.AcquireString()
	buf.WriteString("SELECT ")
	for i, f := range t.Columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(g.sqlScanner(f))
	}
	buf.WriteString(" FROM ")
	var query string
	if method, wrongType := t.Implements(sqlTabler); wrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		buf.WriteString(g.QuoteIdentifier(t.Name))
	} else {
		query = g.Quote(buf.String()) + "+ v.TableName() +"
		buf.Reset()
	}
	buf.WriteString(" WHERE ")
	if pk, ok := t.AutoIncrKey(); ok {
		buf.WriteString(g.QuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(1))
	} else if len(t.Keys) == 1 {
		pk := t.Keys[0]
		buf.WriteString(g.QuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(1))
	} else {
		keyNames := lo.Map(t.Keys, func(v *compiler.Column, _ int) string {
			return g.QuoteIdentifier(v.Name)
		})
		buf.WriteString("(" + strings.Join(keyNames, ",") + ")" + " = ")
		buf.WriteByte('(')
		for i, k := range t.Keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.sqlValuer(k, i))
		}
		buf.WriteByte(')')
	}
	buf.WriteString(" LIMIT 1;")
	g.L("func (v " + t.GoName + ") FindOneByPKStmt() (string, []any) {")
	g.WriteString("return " + query + g.Quote(buf.String()) + ", []any{")
	strpool.ReleaseString(buf)
	if pk, ok := t.AutoIncrKey(); ok {
		g.WriteString(g.getOrValue(importPkgs, "v", pk))
	} else {
		for i, f := range t.Keys {
			if i > 0 {
				g.WriteByte(',')
			}
			g.WriteString(g.getOrValue(importPkgs, "v", f))
		}
	}
	g.WriteString("}\n")
	g.L("}")
}

func (g *Generator) buildInsertOne(importPkgs *Package, t *compiler.Table) {
	var query string
	buf := strpool.AcquireString()
	if method, wrongType := t.Implements(sqlTabler); wrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		buf.WriteString("INSERT INTO " + g.QuoteIdentifier(t.Name))
	} else {
		query = g.Quote("INSERT INTO ") + "+ v.TableName() +"
	}
	buf.WriteString(" (")
	columns := t.InsertColumns()
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(g.QuoteIdentifier(f.Name))
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
		for i, f := range t.Columns {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.sqlScanner(f))
		}
	}
	buf.WriteByte(';')
	// If the columns and after filter columns is the same
	// mean it has no auto increment key
	g.L("func (v " + t.GoName + ") InsertOneStmt() (string, []any) {")
	if len(columns) == len(t.Columns) {
		g.L("return " + query + g.Quote(buf.String()) + ", v.Values()")
		strpool.ReleaseString(buf)
	} else {
		g.WriteString("return " + query + g.Quote(buf.String()) + ", []any{")
		strpool.ReleaseString(buf)
		for i, f := range columns {
			if i > 0 {
				g.WriteByte(',')
			}
			g.WriteString(g.getOrValue(importPkgs, "v", f))
		}
		g.WriteString("}\n")
	}
	g.L("}")
}

func (g *Generator) buildUpdateByPK(importPkgs *Package, t *compiler.Table) {
	buf := strpool.AcquireString()
	var query string
	if method, wrongType := t.Implements(sqlTabler); wrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		buf.WriteString("UPDATE " + g.QuoteIdentifier(t.Name))
	} else {
		query = g.Quote("UPDATE ") + "+ v.TableName() +"
	}
	buf.WriteString(" SET ")
	columns := t.ColumnsWithoutPK()
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(g.QuoteIdentifier(f.Name) + " = " + g.sqlValuer(f, i))
	}
	buf.WriteString(" WHERE ")
	if pk, ok := t.AutoIncrKey(); ok {
		buf.WriteString(g.QuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(len(t.Columns)))
		columns = append(columns, pk)
	} else if len(t.Keys) == 1 {
		pk := t.Keys[0]
		buf.WriteString(g.QuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(len(t.Columns)))
		columns = append(columns, pk)
	} else {
		keyNames := lo.Map(t.Keys, func(v *compiler.Column, _ int) string {
			return g.QuoteIdentifier(v.Name)
		})
		buf.WriteString("(" + strings.Join(keyNames, ",") + ")" + " = ")
		buf.WriteByte('(')
		for i, k := range t.Keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(g.sqlValuer(k, i+len(columns)))
		}
		buf.WriteByte(')')
		columns = append(columns, t.Keys...)
	}
	buf.WriteByte(';')
	g.L("func (v " + t.GoName + ") UpdateOneByPKStmt() (string, []any) {")
	g.WriteString("return " + query + g.Quote(buf.String()) + ", []any{")
	strpool.ReleaseString(buf)
	for i, f := range columns {
		if i > 0 {
			g.WriteByte(',')
		}
		g.WriteString(g.getOrValue(importPkgs, "v", f))
	}
	g.WriteString("}\n")
	g.L("}")
	strpool.ReleaseString(buf)
}

func (g *Generator) getOrValue(importPkgs *Package, obj string, f *compiler.Column) string {
	goPath := obj + f.GoPath
	if f.IsUnderlyingPtr() {
		return obj + "." + valueFunc(f)
	}
	return g.valuer(importPkgs, goPath, f.Type)
}

func (g *Generator) valuer(importPkgs *Package, goPath string, t types.Type) string {
	// This is to prevent auto cast if the value is driver.Value
	switch tv := t.(type) {
	case *types.Basic:
		switch tv.Kind() {
		case types.String, types.Bool, types.Int64, types.Float64:
			return goPath
		}
	case *types.Named:
		if tv.String() == typeOfTime {
			return goPath
		}
	}
	utype, isPtr := underlyingType(t)
	if columnType, ok := g.columnTypes[t.String()]; ok && columnType.Valuer != "" {
		return Expr(columnType.Valuer).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
	} else if _, wrong := types.MissingMethod(utype, goSqlValuer, true); wrong {
		if isPtr {
			return Expr("(database/sql/driver.Valuer)({{goPath}})").Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
		}
		return goPath
	} else if columnType, ok := g.columnDataType(t); ok && columnType.Valuer != "" {
		return Expr(columnType.Valuer).Format(importPkgs, ExprParams{GoPath: goPath, IsPtr: isPtr, Type: t, Len: arraySize(t)})
	} else if isImplemented(utype, textMarshaler) {
		return Expr(goutil.GenericFunc(encoding.TextValue[time.Time], "{{goPath}}")).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
	}
	return Expr(g.defaultColumnTypes["*"].Valuer).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
}

func (g *Generator) scanner(importPkgs *Package, goPath string, t types.Type) string {
	// This is to prevent auto cast if the value is driver.Value
	switch tv := t.(type) {
	case *types.Basic:
		switch tv.Kind() {
		case types.String, types.Bool, types.Int64, types.Float64:
			return goPath
		}
	case *types.Named:
		if tv.String() == typeOfTime {
			return goPath
		}
	}
	ptr, isPtr := pointerType(t)
	if columnType, ok := g.columnTypes[t.String()]; ok && columnType.Scanner != "" {
		return Expr(columnType.Scanner).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
	} else if isImplemented(ptr, goSqlScanner) {
		if isPtr {
			return Expr(goutil.GenericFunc(encoding.PtrScanner[sql.NullString], "{{addr}}")).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
		}
		return goPath
	} else if columnType, ok := g.columnDataType(t); ok && columnType.Scanner != "" {
		return Expr(columnType.Scanner).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr, Len: arraySize(t)})
	} else if isImplemented(ptr, textUnmarshaler) {
		return Expr(goutil.GenericFuncName(encoding.TextScanner[time.Time, *time.Time, *time.Time], "{{elemType}}", "{{addr}}")).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
	}
	return Expr(g.defaultColumnTypes["*"].Scanner).Format(importPkgs, ExprParams{GoPath: goPath, Type: t, IsPtr: isPtr})
}

func (g *Generator) sqlScanner(f *compiler.Column) string {
	// if sqlScanner, ok := f.sqlScanner(); ok {
	// 	matches := sqlFuncRegexp.FindStringSubmatch(sqlScanner("{}"))
	// 	if len(matches) > 4 {
	// 		return matches[1] + matches[2] + g.QuoteIdentifier(f.ColumnName()) + matches[4] + matches[5]
	// 	} else {
	// 		return g.QuoteIdentifier(f.ColumnName())
	// 	}
	// }
	return g.QuoteIdentifier(f.Name)
}

func (g *Generator) sqlValuer(col *compiler.Column, idx int) string {
	// if sqlValuer, ok := f.sqlValuer(); ok {
	// 	matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
	// 	// g.WriteString(fmt.Sprintf("%q + strconv.Itoa((row * noOfColumn) + %d) +%q", matches[1]+string(g.dialect.VarRune()), f.ColumnPos(), matches[5]))
	// 	if len(matches) > 4 {
	// 		return matches[1] + matches[2] + g.dialect.QuoteVar(idx+1) + matches[4] + matches[5]
	// 	}
	// 	return g.dialect.QuoteVar(idx + 1)
	// }
	return g.dialect.QuoteVar(idx + 1)
}

// // Generate migration files for models
// func (g *Generator) genMigrations(pkg *compiler.Package) error {
// 	unix := time.Now().Unix()
// 	for _, tableSchema := range pkg.Tables {
// 		// Each schema should have one migration files
// 		if err := g.genMigration(unix, tableSchema); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (g *Generator) genMigration(unix int64, t *compiler.Table) error {
// 	fileDest := fmt.Sprintf("%s/%d_%s.sql", g.config.Migration.Dir, unix, t.Name)
// 	f, err := os.OpenFile(fileDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	if err := g.dialect.Migrate(g.config.Migration.DSN, f, t); errors.Is(err, dialect.ErrNoNewMigration) {
// 		_ = syscall.Unlink(fileDest)
// 		return nil
// 	} else if err != nil {
// 		return err
// 	}
// 	return nil
// }

func valueFunc(f *compiler.Column) string {
	return f.GoName + "Value()"
}

func strwidth(n int) string {
	str := strconv.Itoa(n)
	return strconv.Itoa(len(str))
}
