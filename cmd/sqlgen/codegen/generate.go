package codegen

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen"
	"github.com/si3nloong/sqlgen/cmd/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/internal/compiler"
	"github.com/si3nloong/sqlgen/cmd/internal/goutil"
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

func (g *Generator) Quote(str string) string {
	buf := make([]byte, 0, len(str))
	buf = append(buf, byte(g.quoteRune))
	for i := range str {
		buf = append(buf, byte(str[i]))
	}
	buf = append(buf, byte(g.quoteRune))
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

func (g *Generator) MustQuoteIdentifier(str string) string {
	return g.dialect.QuoteIdentifier(str)
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
	importPkgs := NewPackage(schema.Pkg.PkgPath, schema.Pkg.Name)
	importPkgs.Import(types.NewPackage("strings", ""))
	importPkgs.Import(types.NewPackage("strconv", ""))
	importPkgs.Import(types.NewPackage("database/sql/driver", ""))
	importPkgs.Import(types.NewPackage("github.com/si3nloong/sqlgen/sequel", ""))

	bw := bytes.NewBufferString(``)
	w := bufio.NewWriter(bw)

	for len(schema.Tables) > 0 {
		t := schema.Tables[0]

		fmt.Fprintln(w)

		// if method, isWrongType := t.Implements(sqlDatabaser); isWrongType {
		// 	g.LogError(fmt.Errorf(`sqlgen: struct %q has function "DatabaseName" but wrong footprint`, t.Name))
		// } else if method != nil && !isWrongType && t.dbName != "" {
		// 	g.L("func (" + t.GoName + ") DatabaseName() string {")
		// 	g.L(`return ` + g.Quote(g.QuoteIdentifier(t.dbName)))
		// 	g.L("}")
		// }
		var readonly bool
		if method, isWrongType := t.PtrImplements(locker); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.Name))
		} else if method == nil && !isWrongType {
			readonly = true
		}

		// Build the "TableName" function which return the table name
		if method, isWrongType := t.Implements(sqlTabler); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			fprintfln(w, "func (%s) %s() string {", t.GoName, methodName(sqlTabler))
			fprintfln(w, "return %s", g.Quote(g.QuoteIdentifier(t.Name)))
			fprintfln(w, "}")
		} else {
			// TODO: we need to do something when table name is declare by user
		}

		if t.HasPK() {
			fprintfln(w, "func (%s) HasPK() {}", t.GoName)
			pk, ok := t.AutoIncrKey()
			if ok {
				fprintfln(w, "func (%s) IsAutoIncr() {}", t.GoName)
				fprintfln(w, `func (v *%s) ScanAutoIncr(val int64) error {
	v.%s = %s(val)
	return nil
}`, t.GoName, pk.GoName, pk.Type)
			} else if len(t.Keys) == 1 {
				pk = t.Keys[0]
			}
			if pk != nil {
				fprintfln(w, `func (v %s) PK() (string, int, any) {
	return %s, %d, %s
}`, t.GoName, g.Quote(g.QuoteIdentifier(pk.Name)), pk.Pos, g.getOrValue(importPkgs, "v", pk))
			} else {
				g.buildCompositeKeys(w, importPkgs, t)
			}
		}

		// Build the "SQLColumns" function which return the column SQL query
		if method, isWrongType := t.Implements(sqlQueryColumner); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "SQLColumns" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			g.buildSqlColumns(w, t)

			// panic("")
			// if _, ok := lo.Find(t.Columns, func(v *compiler.Column) bool {
			// 	_, exists := v.sqlScanner()
			// 	return exists
			// }); ok {
			// fprintfln(w, "func (%s) %s() []string {", t.GoName, methodName(sqlQueryColumner))
			// fmt.Fprintf(w, "return []string{")
			// fmt.Fprintln(w, "}")
			// fmt.Fprintln(w, "}")

			// 	// 	g.WriteString("return []string{")
			// 	// 	for i, f := range t.columns {
			// 	// 		if i > 0 {
			// 	// 			g.WriteByte(',')
			// 	// 		}
			// 	// 		g.WriteString(g.Quote(g.sqlScanner(f)))
			// }
		}

		// Build the "Columns" function which return the column names
		if method, isWrongType := t.Implements(sqlColumner); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Columns" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			fprintfln(w, "func (%s) %s() []string {", t.GoName, methodName(sqlColumner))
			fmt.Fprint(w, "return []string{")
			if len(t.Columns) > 0 {
				fmt.Fprint(w, g.Quote(g.QuoteIdentifier(t.Columns[0].Name)))
				for i := 1; i < len(t.Columns); i++ {
					fmt.Fprint(w, ","+g.Quote(g.QuoteIdentifier(t.Columns[i].Name)))
				}
			}
			fprintfln(w, "} // %d", len(t.Columns))
			fprintfln(w, "}")
		}

		// Build the "Values" function which return the column values
		if method, isWrongType := t.Implements(sqlValuer); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Values" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			g.buildValuer(w, importPkgs, t)
		}

		// Build the "Addrs" function which return the column addressable values
		if method, isWrongType := t.PtrImplements(sqlScanner); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Addrs" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			g.buildScanner(w, importPkgs, t)
		}

		// If the struct has readonly columns
		if !readonly {
			insertColumns := t.InsertColumns()
			if len(insertColumns) > 0 {
				// If the insertable columns is not tally with the table columns means the struct has readonly columns
				if len(insertColumns) != len(t.Columns) {
					fprintfln(w, "func (%s) InsertColumns() []string {", t.GoName)
					fmt.Fprint(w, "return []string{")
					for i, col := range insertColumns {
						if i > 0 {
							fmt.Fprint(w, `,`)
						}
						fmt.Fprint(w, g.Quote(g.QuoteIdentifier(col.Name)))
					}
					fprintfln(w, "} // %d", len(insertColumns))
					fprintfln(w, "}")
				}

				if g.staticVar {
					fprintfln(w, `func (%s) InsertPlaceholders(row int) string {
	return "(%s)" // %d
}`, t.GoName, strings.Repeat(","+g.dialect.QuoteVar(0), len(insertColumns))[1:], len(insertColumns))
				} else {
					fprintfln(w, "func (%s) InsertPlaceholders(row int) string {", t.GoName)
					fprintfln(w, "const noOfColumn = %d", len(insertColumns))
					fmt.Fprint(w, `return "("+`)
					for i := range insertColumns {
						if i > 0 {
							fmt.Fprint(w, `+","+`)
						}
						fmt.Fprintf(w, `%q+ strconv.Itoa((row * noOfColumn) + %d)`, string(g.dialect.VarRune()), i+1)
					}
					fmt.Fprint(w, `+")"`)
					fprintfln(w, "}")
				}

				g.buildInsertOne(w, importPkgs, t)
			}
		}

		if t.HasPK() {
			g.buildFindByPK(w, importPkgs, t)
			if !readonly && len(t.ColumnsWithoutPK()) > 0 {
				g.buildUpdateByPK(w, importPkgs, t)
			}
		}

		// Build getter
		for _, f := range t.Columns {
			fprintfln(w, "func (v %s) %s any {", t.GoName, valueFunc(f))
			queue := []string{}
			// Find all the possible pointer paths
			ptrPaths := f.GoPtrPaths()
			for _, p := range ptrPaths {
				fprintfln(w, "if v%s != nil {", p.GoPath)
				queue = append(queue, "}")
			}

			if f.IsPtr() {
				// Deference the pointer value and return it
				fprintfln(w, "return %s", g.valuer(importPkgs, "*v"+f.GoPath, assertAsPtr[types.Pointer](f.Type).Elem()))
			} else {
				fprintfln(w, "return %s", g.valuer(importPkgs, "v"+f.GoPath, f.Type))
			}
			for len(queue) > 0 {
				fprintfln(w, queue[0])
				queue = queue[1:]
			}
			if len(ptrPaths) > 0 {
				fprintfln(w, "return nil")
			}
			fprintfln(w, "}")
			// 		if idx := strings.Index(typeStr, "."); idx > 0 {
			// 			typeStr = Expr(typeStr).Format(importPkgs)
			// 		}
			// 		var specificType string
			// 		if !typeInferred {
			// 			specificType = "[" + typeStr + "]"
			// 		}

			// 		if sqlValuer, ok := f.sqlValuer(); ok {
			// 			matches := sqlFuncRegexp.FindStringSubmatch(sqlValuer("{}"))
			// 			if len(matches) > 4 {
			// 				g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.SQLColumnValuer[", typeStr, "] {")
			// 				g.L(`return sequel.SQLColumn`+specificType+`(`, g.Quote(g.QuoteIdentifier(f.ColumnName())), `, v.`, f.GoPath()+",", fmt.Sprintf(`func(placeholder string) string { return %q+ placeholder + %q}`, matches[1]+matches[2], matches[4]+matches[5]), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.t), ` })`)
			// 				g.L("}")
			// 			} else {
			// 				g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
			// 				g.L(`return sequel.Column`, specificType, `(`, g.Quote(g.QuoteIdentifier(f.ColumnName())), `, v.`, f.GoPath(), `, func(val `, typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.t), ` })`)
			// 				g.L("}")
			// 			}
			// 		} else {
			// 			g.L("func (v "+t.goName+") ", g.config.Getter.Prefix+f.GoName(), "() sequel.ColumnValuer[", typeStr, "] {")
			// 			g.L("return sequel.Column", specificType, "(", g.Quote(g.QuoteIdentifier(f.ColumnName())), ", v.", f.GoPath(), ", func(val ", typeStr, `) driver.Value { return `, g.valuer(importPkgs, "val", f.t), ` })`)
			// 			g.L("}")
			// 		}
			// }
		}

		for _, f := range t.Columns {
			switch vt := f.Type.(type) {
			// First level struct data type
			case *types.Struct:
				buf := strpool.AcquireString()
				printStruct(buf, importPkgs, vt)
				aliasname := t.GoName + f.GoName + "Field"
				fprintfln(w, "type %s = %s", aliasname, buf)
				fmt.Println(buf.String())
				strpool.ReleaseString(buf)

				fprintfln(w, "func (v %s) %s() sequel.ColumnValuer[%s] {", t.GoName, g.config.Getter.Prefix+f.GoName, aliasname)
				fprintfln(w, "return sequel.Column(%s, v%s, func(val %s) driver.Value {", g.Quote(g.QuoteIdentifier(f.Name)), f.GoPath, aliasname)

			default:
				typeStr := f.Type.String()
				if idx := strings.Index(typeStr, "."); idx > 0 {
					typeStr = Expr(typeStr).Format(importPkgs)
				}
				fprintfln(w, "func (v %s) %s() sequel.ColumnValuer[%s] {", t.GoName, g.config.Getter.Prefix+f.GoName, typeStr)
				fprintfln(w, "return sequel.Column(%s, v%s, func(val %s) driver.Value {", g.Quote(g.QuoteIdentifier(f.Name)), f.GoPath, typeStr)
			}

			if f.IsPtr() {
				fprintfln(w, "if val != nil {")
				// Deference the pointer value and return it
				fprintfln(w, "return %s", g.valuer(importPkgs, "*val", assertAsPtr[types.Pointer](f.Type).Elem()))
				fprintfln(w, "}")
				fprintfln(w, "return nil")
			} else {
				fmt.Fprintf(w, "return %s", g.valuer(importPkgs, "val", f.Type))
			}
			fprintfln(w, "})")
			fprintfln(w, "}")
		}

		if err := w.Flush(); err != nil {
			return err
		}
		schema.Tables = schema.Tables[1:]
	}

	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return err
	}
	fileDest := filepath.Join(dstDir, g.config.Exec.Filename)
	f, err := os.OpenFile(fileDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	rw := bytes.NewBufferString(``)
	g.buildHeader(rw)
	fprintfln(rw, "package %s", schema.Pkg.Name)

	if len(importPkgs.imports) > 0 {
		fprintfln(rw, "import (")
		for _, pkg := range importPkgs.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				fprintfln(rw, strconv.Quote(pkg.Path()))
			} else {
				// If the import is alias import path
				fprintfln(rw, "%s %s", pkg.Name(), strconv.Quote(pkg.Path()))
			}
		}
		fprintfln(rw, ")")
	}

	mustNoError(rw.Write(bw.Bytes()))
	bw.Reset()

	formatted, err := imports.Process("", rw.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}

	mustNoError(f.Write(formatted))

	return f.Close()
}

func (g *Generator) buildHeader(w io.Writer) {
	if !g.config.SkipHeader {
		fmt.Fprintf(w, "// Code generated by sqlgen, version %s; DO NOT EDIT.\n\n", sqlgen.Version)
	}
}

func (g *Generator) buildCompositeKeys(w io.Writer, importPkgs *Package, table *compiler.Table) {
	fprintfln(w, "func (v %s) CompositeKey() ([]string, []int, []any) {", table.GoName)
	fmt.Fprint(w, `return []string{`)
	for i, f := range table.Keys {
		if i > 0 {
			fmt.Fprint(w, `,`)
		}
		fmt.Fprint(w, g.Quote(f.Name))
	}
	fmt.Fprint(w, `}, []int{`)
	for i, f := range table.Keys {
		if i > 0 {
			fmt.Fprint(w, `,`)
		}
		fmt.Fprintf(w, `%d`, f.Pos)
	}
	fmt.Fprintf(w, `}, []any{`)
	for i, k := range table.Keys {
		if i > 0 {
			fmt.Fprint(w, `,`)
		}
		fmt.Fprint(w, g.getOrValue(importPkgs, "v", k))
	}
	fprintfln(w, "}")
	fprintfln(w, "}")
}

func (g *Generator) buildSqlColumns(w io.Writer, t *compiler.Table) error {
	if len(g.config.DataTypes) == 0 {
		return nil
	}

	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	hasSQLScanner := false
	for i, column := range t.Columns {
		if i > 0 {
			fmt.Fprint(blr, ",")
		}
		t, _ := underlyingType(column.Type)
		if typeMapper, ok := g.config.DataTypes[t.String()]; ok && typeMapper.SQLScanner != nil {
			hasSQLScanner = true
			t := template.Must(template.New("scanner").Parse(*typeMapper.SQLScanner))
			buf := strpool.AcquireString()
			if err := t.Execute(buf, g.QuoteIdentifier(column.Name)); err != nil {
				strpool.ReleaseString(buf)
				return err
			}
			fmt.Fprint(blr, g.Quote(buf.String()))
			strpool.ReleaseString(buf)
		} else {
			fmt.Fprint(blr, g.Quote(g.QuoteIdentifier(column.Name)))
		}
	}
	if hasSQLScanner {
		fprintfln(w, "func (%s) %s() []string {", t.GoName, methodName(sqlQueryColumner))
		fprintfln(w, "return []string{%s} // %d", blr.String(), len(t.Columns))
		fprintfln(w, "}")
	}
	return nil
}

func (g *Generator) buildValuer(w io.Writer, importPkgs *Package, table *compiler.Table) {
	columns := table.InsertColumns()
	if len(columns) > 0 {
		fprintfln(w, "func (v %s) %s() []any {", table.GoName, methodName(sqlValuer))
		fprintfln(w, "return []any{")
		tmpl := "%s, // %" + strwidth(len(columns)) + "d - %s"
		for _, f := range columns {
			fprintfln(w, tmpl, g.getOrValue(importPkgs, "v", f), f.Pos, f.Name)
		}
		fprintfln(w, "}")
		fprintfln(w, "}")
	}
}

func (g *Generator) buildScanner(w io.Writer, importPkgs *Package, table *compiler.Table) {
	fprintfln(w, "func (v *%s) %s() []any {", table.GoName, methodName(sqlScanner))
	for _, f := range table.GoPtrPaths() {
		fprintfln(w, "if v%s == nil {", f.GoPath)
		fmt.Fprintf(w, `v%s = new(%s)`, f.GoPath, Expr(strings.TrimPrefix(f.Type.String(), "*")).Format(importPkgs, ExprParams{}))
		fprintfln(w, "}")
	}
	fprintfln(w, "return []any{")
	tmpl := "%s, // %" + strwidth(len(table.Columns)) + "d - %s"
	for _, f := range table.Columns {
		fprintfln(w, tmpl, g.scanner(importPkgs, "&v"+f.GoPath, f.Type), f.Pos, f.Name)
	}
	fprintfln(w, "}")
	fprintfln(w, "}")
}

func (g *Generator) buildFindByPK(w io.Writer, importPkgs *Package, t *compiler.Table) error {
	buf := strpool.AcquireString()
	buf.WriteString("SELECT ")
	for i, f := range t.Columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		scanner, err := g.sqlScanner(f)
		if err != nil {
			return err
		}
		buf.WriteString(scanner)
	}
	buf.WriteString(" FROM ")
	var query string
	if method, isWrongType := t.Implements(sqlTabler); isWrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		buf.WriteString(g.MustQuoteIdentifier(t.Name))
	} else {
		query = g.Quote(buf.String()) + "+ v.TableName() +"
		buf.Reset()
	}
	buf.WriteString(" WHERE ")
	if pk, ok := t.AutoIncrKey(); ok {
		buf.WriteString(g.MustQuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(1))
	} else if len(t.Keys) == 1 {
		pk := t.Keys[0]
		buf.WriteString(g.MustQuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(1))
	} else {
		keyNames := lo.Map(t.Keys, func(v *compiler.Column, _ int) string {
			return g.MustQuoteIdentifier(v.Name)
		})
		buf.WriteString("(" + strings.Join(keyNames, ",") + ")" + " = ")
		buf.WriteByte('(')
		for i, k := range t.Keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			valuer, err := g.sqlValuer(k, i)
			if err != nil {
				return err
			}
			buf.WriteString(valuer)
		}
		buf.WriteByte(')')
	}
	buf.WriteString(" LIMIT 1;")
	fprintfln(w, "func (v "+t.GoName+") FindOneByPKStmt() (string, []any) {")
	fmt.Fprintf(w, `return %s, []any{`, query+g.Quote(buf.String()))
	strpool.ReleaseString(buf)
	if pk, ok := t.AutoIncrKey(); ok {
		fmt.Fprint(w, g.getOrValue(importPkgs, "v", pk))
	} else {
		for i, f := range t.Keys {
			if i > 0 {
				fmt.Fprint(w, `,`)
			}
			fmt.Fprint(w, g.getOrValue(importPkgs, "v", f))
		}
	}
	fprintfln(w, "}")
	fprintfln(w, "}")
	return nil
}

func (g *Generator) buildInsertOne(w io.Writer, importPkgs *Package, t *compiler.Table) error {
	var query string
	buf := strpool.AcquireString()
	if method, isWrongType := t.Implements(sqlTabler); isWrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		buf.WriteString("INSERT INTO " + g.MustQuoteIdentifier(t.Name))
	} else {
		query = g.Quote("INSERT INTO ") + "+ v.TableName() +"
	}
	buf.WriteString(" (")
	columns := t.InsertColumns()
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(g.MustQuoteIdentifier(f.Name))
	}
	buf.WriteString(") VALUES (")
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		valuer, err := g.sqlValuer(f, i)
		if err != nil {
			return err
		}
		buf.WriteString(valuer)
	}
	buf.WriteByte(')')
	if g.config.Driver == Postgres {
		buf.WriteString(" RETURNING ")
		for i, f := range t.Columns {
			if i > 0 {
				buf.WriteByte(',')
			}
			scanner, err := g.sqlScanner(f)
			if err != nil {
				return err
			}
			buf.WriteString(scanner)
		}
	}
	buf.WriteByte(';')
	// If the columns and after filter columns is the same
	// mean it has no auto increment key
	fprintfln(w, `func (v %s) InsertOneStmt() (string, []any) {`, t.GoName)
	if len(columns) == len(t.Columns) {
		fprintfln(w, `return %s, v.Values()`, query+g.Quote(buf.String()))
		strpool.ReleaseString(buf)
	} else {
		fmt.Fprintf(w, `return %s, []any{`, query+g.Quote(buf.String()))
		strpool.ReleaseString(buf)
		for i, f := range columns {
			if i > 0 {
				fmt.Fprint(w, `,`)
			}
			fmt.Fprint(w, g.getOrValue(importPkgs, "v", f))
		}
		fprintfln(w, "}")
	}
	fprintfln(w, "}")
	return nil
}

func (g *Generator) buildUpdateByPK(w io.Writer, importPkgs *Package, t *compiler.Table) error {
	buf := strpool.AcquireString()
	var query string
	if method, isWrongType := t.Implements(sqlTabler); isWrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		buf.WriteString("UPDATE " + g.MustQuoteIdentifier(t.Name))
	} else {
		query = g.Quote("UPDATE ") + "+ v.TableName() +"
	}
	buf.WriteString(" SET ")
	columns := t.ColumnsWithoutPK()
	for i, f := range columns {
		if i > 0 {
			buf.WriteByte(',')
		}
		valuer, err := g.sqlValuer(f, i)
		if err != nil {
			return err
		}
		buf.WriteString(g.MustQuoteIdentifier(f.Name) + " = " + valuer)
	}
	buf.WriteString(" WHERE ")
	if pk, ok := t.AutoIncrKey(); ok {
		buf.WriteString(g.MustQuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(len(t.Columns)))
		columns = append(columns, pk)
	} else if len(t.Keys) == 1 {
		pk := t.Keys[0]
		buf.WriteString(g.MustQuoteIdentifier(pk.Name) + " = " + g.dialect.QuoteVar(len(t.Columns)))
		columns = append(columns, pk)
	} else {
		keyNames := lo.Map(t.Keys, func(v *compiler.Column, _ int) string {
			return g.MustQuoteIdentifier(v.Name)
		})
		buf.WriteString("(" + strings.Join(keyNames, ",") + ")" + " = ")
		buf.WriteByte('(')
		for i, k := range t.Keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			valuer, err := g.sqlValuer(k, i+len(columns))
			if err != nil {
				return err
			}
			buf.WriteString(valuer)
		}
		buf.WriteByte(')')
		columns = append(columns, t.Keys...)
	}
	buf.WriteByte(';')
	fprintfln(w, `func (v %s) UpdateOneByPKStmt() (string, []any) {`, t.GoName)
	fmt.Fprintf(w, "return %s, []any{", query+g.Quote(buf.String()))
	strpool.ReleaseString(buf)
	if len(columns) > 0 {
		fmt.Fprint(w, g.getOrValue(importPkgs, "v", columns[0]))
		for i := 1; i < len(columns); i++ {
			fmt.Fprint(w, ","+g.getOrValue(importPkgs, "v", columns[i]))
		}
	}
	fprintfln(w, "}")
	fprintfln(w, "}")
	strpool.ReleaseString(buf)
	return nil
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

func (g *Generator) sqlScanner(column *compiler.Column) (string, error) {
	t, _ := underlyingType(column.Type)
	if typeMapper, ok := g.config.DataTypes[t.String()]; ok && typeMapper.SQLScanner != nil {
		t := template.Must(template.New("scanner").Parse(*typeMapper.SQLScanner))
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		if err := t.Execute(buf, g.MustQuoteIdentifier(column.Name)); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return g.MustQuoteIdentifier(column.Name), nil
}

func (g *Generator) sqlValuer(column *compiler.Column, idx int) (string, error) {
	t, _ := underlyingType(column.Type)
	if typeMapper, ok := g.config.DataTypes[t.String()]; ok && typeMapper.SQLValuer != nil {
		t := template.Must(template.New("valuer").Parse(*typeMapper.SQLValuer))
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		if err := t.Execute(buf, g.dialect.QuoteVar(idx+1)); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return g.dialect.QuoteVar(idx + 1), nil
}

func printStruct(w io.Writer, importPkgs *Package, s *types.Struct) {
	fmt.Fprintln(w, "struct{")
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		switch ft := f.Type().(type) {
		case *types.Struct:
			if !f.Embedded() {
				fmt.Fprint(w, f.Name()+" ")
				printStruct(w, importPkgs, ft)
			} else {
				typeStr := f.Type().String()
				if idx := strings.Index(typeStr, "."); idx > 0 {
					typeStr = Expr(typeStr).Format(importPkgs)
				}
				fmt.Fprint(w, f.Name()+" "+typeStr)
				if tag := s.Tag(i); len(tag) > 0 {
					fmt.Fprintf(w, " `%s`", s.Tag(i))
				}
				fmt.Fprintln(w)
			}
		default:
			typeStr := f.Type().String()
			if idx := strings.Index(typeStr, "."); idx > 0 {
				typeStr = Expr(typeStr).Format(importPkgs)
			}
			if f.Embedded() {
				fmt.Fprint(w, typeStr)
			} else {
				fmt.Fprint(w, f.Name()+" "+typeStr)
			}
			if tag := s.Tag(i); len(tag) > 0 {
				fmt.Fprintf(w, " `%s`", s.Tag(i))
			}
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w, "}")
}

func methodName(i *types.Interface) string {
	return i.Method(0).Name()
}

func valueFunc(f *compiler.Column) string {
	return f.GoName + "Value()"
}

func strwidth(n int) string {
	str := strconv.Itoa(n)
	return strconv.Itoa(len(str))
}

func fprintfln(w io.Writer, format string, values ...any) {
	mustNoError(fmt.Fprintf(w, format+"\n", values...))
}
