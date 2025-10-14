package codegen

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"go/types"
	"io"
	"iter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/goutil"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/go/packages"
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
	pkg *packages.Package,
	tables iter.Seq2[*compiler.Table, error],
) error {
	importPkgs := NewPackage(pkg.PkgPath, pkg.Name)
	importPkgs.Import(types.NewPackage("strings", ""))
	importPkgs.Import(types.NewPackage("strconv", ""))
	importPkgs.Import(types.NewPackage("database/sql/driver", ""))
	importPkgs.Import(types.NewPackage("github.com/si3nloong/sqlgen/sequel", ""))

	bw := bytes.NewBufferString("")
	w := bufio.NewWriter(bw)

	next, stop := iter.Pull2(tables)
	defer stop()

loop:
	for {
		t, err, ok := next()
		if err != nil {
			return err
		} else if !ok {
			break loop
		}

		fmt.Fprintln(w)

		// if method, isWrongType := t.Implements(sqlDatabaser); isWrongType {
		// 	g.LogError(fmt.Errorf(`sqlgen: struct %q has function "DatabaseName" but wrong footprint`, t.Name))
		// } else if method != nil && !isWrongType && t.dbName != "" {
		// 	g.L("func (" + t.GoName + ") DatabaseName() string {")
		// 	g.L(`return ` + g.Quote(g.QuoteIdentifier(t.dbName)))
		// 	g.L("}")
		// }

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

		pk, hasPK := t.PK()
		if hasPK {
			fprintfln(w, "func (%s) HasPK() {}", t.GoName)
			switch v := pk.(type) {
			case *compiler.AutoIncrPrimaryKey:
				fprintfln(w, "func (%s) IsAutoIncr() {}", t.GoName)
				fprintfln(w, "func (v *%s) ScanAutoIncr(val int64) error {", t.GoName)
				fprintfln(w, "v.%s = %s(val)", v.GoName(), v.GoType())
				fprintfln(w, "return nil")
				fprintfln(w, "}")
				fprintfln(w, "func (v %s) PK() (string, int, any) {", t.GoName)
				fprintfln(w, "return %s, %d, %s", g.Quote(g.QuoteIdentifier(v.Name())), v.Pos(), g.getOrValue(importPkgs, "v", v))
				fprintfln(w, "}")
			case *compiler.PrimaryKey:
				fprintfln(w, "func (v %s) PK() (string, int, any) {", t.GoName)
				fprintfln(w, "return %s, %d, %s", g.Quote(g.QuoteIdentifier(v.Name())), v.Pos(), g.getOrValue(importPkgs, "v", v))
				fprintfln(w, "}")
			case *compiler.CompositePrimaryKey:
				g.buildCompositeKeys(w, importPkgs, t.GoName, v)
			default:
				panic("unreachable")

			}
		}

		// Build the "Columns" function which return the column names
		if method, isWrongType := t.Implements(sqlColumner); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Columns" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			fprintfln(w, "func (%s) %s() []string {", t.GoName, methodName(sqlColumner))
			fmt.Fprint(w, "return []string{")
			if len(t.Columns) > 0 {
				fmt.Fprint(w, g.Quote(t.Columns[0].Name()))
				for i := 1; i < len(t.Columns); i++ {
					fmt.Fprint(w, ","+g.Quote(t.Columns[i].Name()))
				}
			}
			fprintfln(w, "} // %d", len(t.Columns))
			fprintfln(w, "}")
		}

		if !t.Readonly {
			// Build the "Values" function which return the column values
			if method, isWrongType := t.Implements(sqlValuer); isWrongType {
				g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Values" but wrong footprint`, t.Name))
			} else if method != nil && !isWrongType {
				g.buildValuer(w, importPkgs, t)
			}
		}

		// Build the "Addrs" function which return the column addressable values
		if method, isWrongType := t.PtrImplements(sqlScanner); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "Addrs" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			g.buildScanner(w, importPkgs, t)
		}

		// If the struct has readonly columns
		if !t.Readonly {
			insertColumns := t.InsertColumns()
			if n := len(insertColumns); n > 0 {
				// If the insertable columns is not tally with the table columns means the struct has readonly columns
				if n != len(t.Columns) {
					fprintfln(w, "func (%s) InsertColumns() []string {", t.GoName)
					fmt.Fprint(w, "return []string{")
					fmt.Fprint(w, g.Quote(g.QuoteIdentifier(insertColumns[0].Name())))
					for i := 1; i < n; i++ {
						fmt.Fprint(w, ","+g.Quote(g.QuoteIdentifier(insertColumns[i].Name())))
					}
					fprintfln(w, "} // %d", n)
					fprintfln(w, "}")
				}

				fprintfln(w, "func (%s) InsertPlaceholders(row int) string {", t.GoName)
				if g.staticVar {
					fprintfln(w, `return "(%s)" // %d`, strings.Repeat(","+g.dialect.QuoteVar(0), len(insertColumns))[1:], len(insertColumns))
				} else {
					fprintfln(w, "const noOfColumn = %d", len(insertColumns))
					fmt.Fprint(w, `return "("+`)
					for i := range insertColumns {
						if i > 0 {
							fmt.Fprint(w, `+","+`)
						}
						fmt.Fprintf(w, `%c+ strconv.Itoa((row * noOfColumn) + %d)`, g.dialect.VarRune(), i+1)
					}
					fmt.Fprint(w, `+")"`)
				}
				fprintfln(w, "}")
			}
			g.buildInsertOne(w, importPkgs, t)
		}

		if hasPK {
			g.buildFindByPK(w, importPkgs, t)
			if !t.Readonly {
				g.buildUpdateByPK(w, importPkgs, t)
			}
		}

		// Build the "SQLColumns" function which return the column SQL query
		if method, isWrongType := t.Implements(sqlQueryColumner); isWrongType {
			g.LogError(fmt.Errorf(`sqlgen: struct %q has function "SQLColumns" but wrong footprint`, t.Name))
		} else if method != nil && !isWrongType {
			g.buildSqlColumns(w, t)
		}

		// Build getter function for each column
		if !t.Readonly {
			for _, f := range t.Columns {
				fprintfln(w, "func (v %s) %s any {", t.GoName, valueFunc(f))
				queue := []string{}
				// Find all the possible pointer paths
				for _, paths := range f.GoPtrPaths() {
					for _, p := range paths {
						fprintfln(w, "if v.%s != nil {", p.GoPath())
						queue = append(queue, "}")
					}
				}

				if f.IsGoPtr() {
					// Deference the pointer value and return it
					fprintfln(w, "return %s", g.valuer(importPkgs, "*v."+f.GoPath(), assertAsPtr[types.Pointer](f.GoType()).Elem()))
				} else {
					fprintfln(w, "return %s", g.valuer(importPkgs, "v."+f.GoPath(), f.GoType()))
				}
				if len(queue) > 0 {
					for len(queue) > 0 {
						fprintfln(w, queue[0])
						queue = queue[1:]
					}
					fprintfln(w, "return nil")
				}
				fprintfln(w, "}")
			}
		}

		// Build the valuer function for each column
		for _, f := range t.Columns {
			var typeStr string
			isBasic := g.isBasicType(f.GoType())
			// underlyingType, _ := underlyingType(f.GoType())
			switch vt := f.GoType().(type) {
			// First level struct data type
			case *types.Struct:
				buf := strpool.AcquireString()
				printStruct(buf, importPkgs, vt)
				aliasname := t.GoName + f.GoName() + "InlineStruct"
				fprintfln(w, "type %s = %s", aliasname, buf)
				strpool.ReleaseString(buf)

				fprintfln(w, "func (v %s) %s() sequel.ColumnConvertClause[%s] {", t.GoName, g.config.Getter.Prefix+f.GoName(), aliasname)
				// fprintfln(w, "return sequel.Column(%s, v%s, func(val %s) any {", g.Quote(g.QuoteIdentifier(f.Name)), f.GoPath, aliasname)
				typeStr = aliasname

			default:
				typeStr = f.GoType().String()
				if idx := strings.Index(typeStr, "."); idx > 0 {
					typeStr = Expr(typeStr).Format(importPkgs)
				}
				if isBasic {
					fprintfln(w, "func (v %s) %s() sequel.ColumnClause[%s] {", t.GoName, g.config.Getter.Prefix+f.GoName(), typeStr)
				} else {
					fprintfln(w, "func (v %s) %s() sequel.ColumnConvertClause[%s] {", t.GoName, g.config.Getter.Prefix+f.GoName(), typeStr)
				}
			}

			if isBasic {
				fprintfln(w, "return sequel.BasicColumn(%s, v.%s)", g.Quote(g.QuoteIdentifier(f.Name())), f.GoPath())
			} else {
				fprintfln(w, "return sequel.Column(%s, v.%s, func(val %s) any {", g.Quote(g.QuoteIdentifier(f.Name())), f.GoPath(), typeStr)
				// fprintfln(w, "if val != nil {")
				// fprintfln(w, "return %s", g.valuer(importPkgs, "*val", assertAsPtr[types.Pointer](f.Type).Elem()))
				if f.IsGoPtr() {
					fprintfln(w, "if val != nil {")
					// Deference the pointer value and return it
					fprintfln(w, "return %s", g.valuer(importPkgs, "*val", assertAsPtr[types.Pointer](f.GoType()).Elem()))
					fprintfln(w, "}")
					fprintfln(w, "return nil")
				} else {
					fmt.Fprintf(w, "return %s", g.valuer(importPkgs, "val", f.GoType()))
				}
				fprintfln(w, "})")
			}
			fprintfln(w, "}")
		}

		if err := w.Flush(); err != nil {
			return err
		}
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

	fw := bytes.NewBufferString("")
	g.buildHeader(fw)
	fprintfln(fw, "package %s", pkg.Name)

	if len(importPkgs.imports) > 0 {
		fprintfln(fw, "import (")
		for _, pkg := range importPkgs.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				fprintfln(fw, strconv.Quote(pkg.Path()))
			} else {
				// If the import is alias import path
				fprintfln(fw, "%s %s", pkg.Name(), strconv.Quote(pkg.Path()))
			}
		}
		fprintfln(fw, ")")
	}

	if _, err := fw.Write(bw.Bytes()); err != nil {
		return err
	}
	bw.Reset()

	// println(fw.String())
	formatted, err := imports.Process("", fw.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}

	if _, err := f.Write(formatted); err != nil {
		return err
	}
	return f.Close()
}

func (g *Generator) buildHeader(w io.Writer) {
	fmt.Fprintf(w, "// Code generated by sqlgen. DO NOT EDIT.\n\n")
}

func (g *Generator) buildCompositeKeys(w io.Writer, importPkgs *Package, goName string, k *compiler.CompositePrimaryKey) {
	// column names, indexes, values
	fprintfln(w, "func (v %s) CompositeKey() ([]string, []int, []any) {", goName)
	w1 := strpool.AcquireString()
	w2 := strpool.AcquireString()
	w3 := strpool.AcquireString()
	if n := len(k.Columns); n > 0 {
		column := k.Columns[0]
		fmt.Fprint(w1, g.Quote(column.Name()))
		fmt.Fprint(w2, strconv.Itoa(column.Pos()))
		fmt.Fprint(w3, g.valuer(importPkgs, "v."+column.GoPath(), column.GoType()))
		for i := 1; i < n; i++ {
			column := k.Columns[i]
			fmt.Fprint(w1, ","+g.Quote(column.Name()))
			fmt.Fprint(w2, ","+strconv.Itoa(column.Pos()))
			fmt.Fprint(w3, ","+g.valuer(importPkgs, "v."+column.GoPath(), column.GoType()))
		}
	}
	fprintfln(w, "return []string{%s}, []int{%s}, []any{%s}", w1, w2, w3)
	strpool.ReleaseString(w3)
	strpool.ReleaseString(w2)
	strpool.ReleaseString(w1)
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
		t, _ := underlyingType(column.GoType())
		if typeMapper, ok := g.config.DataTypes[t.String()]; ok && typeMapper.SQLScanner != nil {
			hasSQLScanner = true
			t := template.Must(template.New("scanner").Parse(*typeMapper.SQLScanner))
			buf := strpool.AcquireString()
			if err := t.Execute(buf, g.MustQuoteIdentifier(column.Name())); err != nil {
				strpool.ReleaseString(buf)
				return err
			}
			fmt.Fprint(blr, g.Quote(buf.String()))
			strpool.ReleaseString(buf)
		} else {
			fmt.Fprint(blr, g.Quote(g.MustQuoteIdentifier(column.Name())))
		}
	}
	if hasSQLScanner {
		fprintfln(w, "func (%s) %s() []string {", t.GoName, methodName(sqlQueryColumner))
		fprintfln(w, "return []string{%s} // %d", blr.String(), len(t.Columns))
		fprintfln(w, "}")
	}
	return nil
}

func (g *Generator) buildValuer(w io.Writer, importPkgs *Package, t *compiler.Table) {
	if n := len(t.Columns); n > 0 {
		fprintfln(w, "func (v %s) %s() []any {", t.GoName, methodName(sqlValuer))
		fprintfln(w, "return []any{")
		tmpl := "%s, // %" + stfwidth(n) + "d - %s"
		for _, f := range t.Columns {
			fprintfln(w, tmpl, g.getOrValue(importPkgs, "v", f), f.Pos(), f.Name())
		}
		fprintfln(w, "}")
		fprintfln(w, "}")
	}
}

func (g *Generator) buildScanner(w io.Writer, importPkgs *Package, table *compiler.Table) {
	fprintfln(w, "func (v *%s) %s() []any {", table.GoName, methodName(sqlScanner))
	// Initialize all pointer fields before we passed those property addr
	for p := range table.ColumnGoPtrPaths() {
		fprintfln(w, "if v.%s == nil {", p.GoPath())
		fmt.Fprintf(w, "v.%s = new(%s)", p.GoPath(), Expr(strings.TrimPrefix(p.GoType().String(), "*")).Format(importPkgs, ExprParams{}))
		fprintfln(w, "}")
	}
	fprintfln(w, "return []any{")
	tmpl := "%s, // %" + stfwidth(len(table.Columns)) + "d - %s"
	for _, f := range table.Columns {
		fprintfln(w, tmpl, g.scanner(importPkgs, "&v."+f.GoPath(), f.GoType()), f.Pos(), f.Name())
	}
	fprintfln(w, "}")
	fprintfln(w, "}")
}

func (g *Generator) buildFindByPK(w io.Writer, importPkgs *Package, t *compiler.Table) error {
	w1 := bytes.NewBufferString("")
	defer w1.Reset()
	fmt.Fprintf(w1, "%cSELECT ", g.quoteRune)
	if n := len(t.Columns); n > 0 {
		column := t.Columns[0]
		scanner, err := g.sqlScanner(column)
		if err != nil {
			return err
		}
		fmt.Fprint(w1, scanner)
		for i := 1; i < n; i++ {
			column = t.Columns[i]
			scanner, err := g.sqlScanner(column)
			if err != nil {
				return err
			}
			fmt.Fprint(w1, ","+scanner)
		}
		fmt.Fprint(w1, " FROM ")
	}
	if method, isWrongType := t.Implements(sqlTabler); isWrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		fmt.Fprint(w1, g.MustQuoteIdentifier(t.Name))
	} else {
		fmt.Fprintf(w1, "%c+ v.TableName() +%c", g.quoteRune, g.quoteRune)
	}
	fmt.Fprint(w1, " WHERE ")
	pk, ok := t.PK()
	if !ok {
		return fmt.Errorf(`sqlgen:`)
	}
	w2 := bytes.NewBufferString("")
	defer w2.Reset()
	switch v := pk.(type) {
	case *compiler.AutoIncrPrimaryKey:
		fmt.Fprint(w1, g.MustQuoteIdentifier(v.Name())+" = "+g.dialect.QuoteVar(1))
		fmt.Fprint(w2, g.valuer(importPkgs, "v."+v.GoPath(), v.GoType()))
	case *compiler.PrimaryKey:
		fmt.Fprint(w1, g.MustQuoteIdentifier(v.Name())+" = "+g.dialect.QuoteVar(1))
		fmt.Fprint(w2, g.valuer(importPkgs, "v."+v.GoPath(), v.GoType()))
	case *compiler.CompositePrimaryKey:
		if n := len(v.Columns); n > 0 {
			w3 := strpool.AcquireString()
			column := v.Columns[0]
			fmt.Fprint(w1, "("+g.MustQuoteIdentifier(column.Name()))
			fmt.Fprint(w2, g.valuer(importPkgs, "v."+column.GoPath(), column.GoType()))
			fmt.Fprint(w3, g.dialect.QuoteVar(1))
			for i := 1; i < n; i++ {
				column = v.Columns[i]
				fmt.Fprint(w1, ","+g.MustQuoteIdentifier(column.Name()))
				fmt.Fprint(w2, ","+g.valuer(importPkgs, "v."+column.GoPath(), column.GoType()))
				fmt.Fprint(w3, ","+g.dialect.QuoteVar(i+1))
			}
			fmt.Fprintf(w1, ") = (%s)", w3)
			strpool.ReleaseString(w3)
		}
	}
	fmt.Fprintf(w1, " LIMIT 1;%c", g.quoteRune)
	fprintfln(w, "func (v "+t.GoName+") FindOneByPKStmt() (string, []any) {")
	fprintfln(w, "return %s, []any{%s}", w1, w2)
	fprintfln(w, "}")
	return nil
}

func (g *Generator) buildInsertOne(w io.Writer, importPkgs *Package, t *compiler.Table) error {
	columns := t.InsertColumns()
	noOfColumns := len(columns)
	if noOfColumns == 0 {
		return nil
	}

	// Build the insert statement
	w1 := bytes.NewBufferString("")
	defer w1.Reset()
	if method, isWrongType := t.Implements(sqlTabler); isWrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
		fmt.Fprintf(w1, "%cINSERT INTO %s (", g.quoteRune, g.MustQuoteIdentifier(t.Name))
	} else if method != nil {
		fmt.Fprintf(w1, "%cINSERT INTO %s (", g.quoteRune, g.MustQuoteIdentifier(t.Name))
	} else {
		fmt.Fprintf(w1, "%cINSERT INTO %c+ v.TableName() +%c (", g.quoteRune, g.quoteRune, g.quoteRune)
	}
	w2 := bytes.NewBufferString("")
	if g.config.Driver == Postgres {
		w3 := bytes.NewBufferString("")
		w4 := bytes.NewBufferString("")
		defer w3.Reset()
		defer w4.Reset()
		column := columns[0]
		valuer, err := g.sqlValuer(column, 0)
		if err != nil {
			return err
		}
		scanner, err := g.sqlScanner(column)
		if err != nil {
			return err
		}
		fmt.Fprint(w1, g.MustQuoteIdentifier(column.Name()))
		fmt.Fprint(w2, g.getOrValue(importPkgs, "v", column))
		fmt.Fprint(w3, valuer)
		fmt.Fprint(w4, scanner)
		for i := 1; i < noOfColumns; i++ {
			column = columns[i]
			valuer, err := g.sqlValuer(column, i)
			if err != nil {
				return err
			}
			scanner, err := g.sqlScanner(column)
			if err != nil {
				return err
			}
			fmt.Fprint(w1, ","+g.MustQuoteIdentifier(column.Name()))
			fmt.Fprint(w2, ","+g.getOrValue(importPkgs, "v", column))
			fmt.Fprint(w3, ","+valuer)
			fmt.Fprint(w4, ","+scanner)
		}
		fmt.Fprintf(w1, "(%s) VALUES (%s) RETURNING (%s)", w2, w3, w4)
	} else {
		column := columns[0]
		valuer, err := g.sqlValuer(column, 0)
		if err != nil {
			return err
		}
		w3 := bytes.NewBufferString("")
		defer w3.Reset()
		fmt.Fprint(w1, g.MustQuoteIdentifier(column.Name()))
		fmt.Fprint(w2, g.getOrValue(importPkgs, "v", column))
		fmt.Fprint(w3, valuer)
		for i := 1; i < noOfColumns; i++ {
			column = columns[i]
			valuer, err := g.sqlValuer(column, i)
			if err != nil {
				return err
			}
			fmt.Fprint(w1, ","+g.MustQuoteIdentifier(column.Name()))
			fmt.Fprint(w2, ","+g.getOrValue(importPkgs, "v", column))
			fmt.Fprint(w3, ","+valuer)
		}
		fmt.Fprintf(w1, ") VALUES (%s)", w3)
	}
	fmt.Fprintf(w1, ";%c", g.quoteRune)
	// If the columns and after filter columns is the same
	// mean it has no auto increment key
	fprintfln(w, "func (v %s) InsertOneStmt() (string, []any) {", t.GoName)
	if len(columns) == len(t.Columns) {
		fprintfln(w, "return %s, v.Values()", w1)
	} else {
		fprintfln(w, "return %s, []any{%s}", w1, w2)
	}
	fprintfln(w, "}")
	return nil
}

func (g *Generator) buildUpdateByPK(w io.Writer, importPkgs *Package, t *compiler.Table) error {
	columns := t.ColumnsExceptPK()
	noOfColumns := len(columns)
	if noOfColumns == 0 {
		return nil
	}
	// Build the update statement
	w1 := bytes.NewBufferString("")
	defer w1.Reset()
	if method, isWrongType := t.Implements(sqlTabler); isWrongType {
		g.LogError(fmt.Errorf(`sqlgen: struct %q has function "TableName" but wrong footprint`, t.GoName))
	} else if method != nil {
		fmt.Fprintf(w1, "%cUPDATE %s", g.quoteRune, g.MustQuoteIdentifier(t.Name))
	} else {
		fmt.Fprintf(w1, "%cUPDATE %c+ v.TableName() +%c", g.quoteRune, g.quoteRune, g.quoteRune)
	}
	fmt.Fprint(w1, " SET ")
	w2 := bytes.NewBufferString("")
	defer w2.Reset()
	column := columns[0]
	valuer, err := g.sqlValuer(column, 0)
	if err != nil {
		return err
	}
	fmt.Fprintf(w1, "%s = %s", g.MustQuoteIdentifier(column.Name()), valuer)
	fmt.Fprint(w2, g.getOrValue(importPkgs, "v", column)+",")
	for i := 1; i < noOfColumns; i++ {
		column = columns[i]
		valuer, err := g.sqlValuer(column, i)
		if err != nil {
			return err
		}
		fmt.Fprintf(w1, ",%s = %s", g.MustQuoteIdentifier(column.Name()), valuer)
		fmt.Fprint(w2, g.getOrValue(importPkgs, "v", column)+",")
	}
	fmt.Fprint(w1, " WHERE ")
	switch pk := t.MustPK().(type) {
	case *compiler.AutoIncrPrimaryKey:
		fmt.Fprintf(w1, "%s = %s", g.MustQuoteIdentifier(pk.Name()), g.dialect.QuoteVar(noOfColumns))
		fmt.Fprint(w2, g.getOrValue(importPkgs, "v", pk))
	case *compiler.PrimaryKey:
		fmt.Fprintf(w1, "%s = %s", g.MustQuoteIdentifier(pk.Name()), g.dialect.QuoteVar(noOfColumns))
		fmt.Fprint(w2, g.getOrValue(importPkgs, "v", pk))
	case *compiler.CompositePrimaryKey:
		if n := len(pk.Columns); n > 0 {
			column := pk.Columns[0]
			w3 := strpool.AcquireString()
			fmt.Fprintf(w1, "(%s", g.MustQuoteIdentifier(column.Name()))
			fmt.Fprint(w2, g.getOrValue(importPkgs, "v", column))
			fmt.Fprint(w3, g.dialect.QuoteVar(noOfColumns+1))
			for i := 1; i < n; i++ {
				column = pk.Columns[i]
				fmt.Fprintf(w1, ",%s", g.MustQuoteIdentifier(column.Name()))
				fmt.Fprintf(w2, ",%s", g.getOrValue(importPkgs, "v", column))
				fmt.Fprintf(w3, ",%s", g.dialect.QuoteVar(noOfColumns+i+1))
			}
			fmt.Fprintf(w1, ") = (%s)", w3)
			strpool.ReleaseString(w3)
		}
	}
	fmt.Fprintf(w1, ";%c", g.quoteRune)
	fprintfln(w, "func (v %s) UpdateOneByPKStmt() (string, []any) {", t.GoName)
	fprintfln(w, "return %s, []any{%s}", w1, w2)
	fprintfln(w, "}")
	return nil
}

func (g *Generator) getOrValue(importPkgs *Package, obj string, f compiler.Column) string {
	goPath := obj + "." + f.GoPath()
	if f.IsUnderlyingPtr() {
		return obj + "." + valueFunc(f)
	}
	return g.valuer(importPkgs, goPath, f.GoType())
}

func (g *Generator) isBasicType(t types.Type) bool {
	// This is to prevent auto cast if the value is driver.Value
	switch tv := t.(type) {
	case *types.Basic:
		switch tv.Kind() {
		case types.String, types.Bool, types.Int64, types.Float64:
			return true
		}
	case *types.Named:
		if tv.String() == typeOfTime {
			return true
		}
	}
	return false
}

func (g *Generator) valuer(importPkgs *Package, goPath string, t types.Type) string {
	// This is to prevent auto cast if the value is driver.Value
	if g.isBasicType(t) {
		return goPath
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

func (g *Generator) sqlScanner(column compiler.Column) (string, error) {
	t, _ := underlyingType(column.GoType())
	if typeMapper, ok := g.config.DataTypes[t.String()]; ok && typeMapper.SQLScanner != nil {
		t := template.Must(template.New("scanner").Parse(*typeMapper.SQLScanner))
		buf := strpool.AcquireString()
		defer strpool.ReleaseString(buf)
		if err := t.Execute(buf, g.MustQuoteIdentifier(column.Name())); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return g.MustQuoteIdentifier(column.Name()), nil
}

func (g *Generator) sqlValuer(column compiler.Column, idx int) (string, error) {
	t, _ := underlyingType(column.GoType())
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

func valueFunc(f compiler.Column) string {
	return f.GoName() + "Value()"
}

func stfwidth(n int) string {
	str := strconv.Itoa(n)
	return strconv.Itoa(len(str))
}

func fprintfln(w io.Writer, format string, values ...any) {
	mustNoError(fmt.Fprintf(w, format+"\n", values...))
}
