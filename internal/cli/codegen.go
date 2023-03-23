package cli

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type formatCase string

const (
	snakeCase  formatCase = "snakecase"
	kebabCase  formatCase = "kebabcase"
	paskalCase formatCase = "paskalcase"
	camelCase  formatCase = "camelcase"
)

var sqlTypes = map[string][2]string{
	"string":    {"string", "types.String"},
	"int":       {"int64", "types.Integer"},
	"int8":      {"int64", "types.Integer"},
	"int16":     {"int64", "types.Integer"},
	"int32":     {"int64", "types.Integer"},
	"int64":     {"int64", "types.Integer"},
	"uint":      {"uint64", "types.Integer"},
	"uint8":     {"uint64", "types.Integer"},
	"uint16":    {"uint64", "types.Integer"},
	"uint32":    {"uint64", "types.Integer"},
	"uint64":    {"uint64", "types.Integer"},
	"float32":   {"float64", "types.Float"},
	"time.Time": {"time.Time", "types.String"},
	"[]string":  {"github.com/si3nloong/sqlgen/sql/encoding.MarshalStringList", "github.com/si3nloong/sqlgen/sql/encoding.MarshalStringList"},
	"[][]byte":  {"[][]byte", "github.com/si3nloong/sqlgen/sql/encoding.MarshalStringList"},
	"[]bool":    {"[]bool", "github.com/si3nloong/sqlgen/sql/encoding.MarshalBoolList"},
	"[]uint64":  {"[]uint64", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]uint32":  {"[]uint32", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]uint16":  {"[]uint16", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]uint8":   {"[]uint8", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]uint":    {"[]uint", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]int64":   {"[]int64", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]int32":   {"[]int32", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]int16":   {"[]int16", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]int8":    {"[]int8", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]int":     {"[]int", "github.com/si3nloong/sqlgen/sql/encoding.MarshalIntList"},
	"[]float32": {"[]float32", "github.com/si3nloong/sqlgen/sql/encoding.MarshalFloatList"},
	"[]float64": {"[]float64", "github.com/si3nloong/sqlgen/sql/encoding.MarshalFloatList"},
}

type RenameFunc func(string) string

type codegen struct {
	// Go template engine
	tmpl *template.Template

	// Go ast
	goast *ast.File

	importCache map[string]string
}

func newCodegen() *codegen {
	return &codegen{
		tmpl: template.New("template.go").Funcs(template.FuncMap{
			"quote": strconv.Quote,
			"import": func(pkg string, aliases ...string) string {
				// add to cache
				log.Println("Imported pkg ->", pkg)
				return pkg
			},
			"cast": func(n string, f Field) string {
				log.Println("actualType ->", f.ActualType, ", type ->", f.Type)
				if v, ok := sqlTypes[f.ActualType]; ok && v[0] != f.Type {
					pos := strings.LastIndex(v[0], ".")
					if pos >= 0 {
						pkg := string(v[0][:pos])
						// log.Println(pkg, path.Base(pkg), string(v[0][pos:]))
						// log.Println("pos is package")
						return path.Base(pkg) + string(v[0][pos:]) + "(" + n + ")"
					}
					return v[0] + "(" + n + ")"
				}
				return n
			},
			"addr": func(n string, f Field) string {
				log.Println(f.Type, f.ActualType)
				// if v, ok := sqlTypes[f.ActualType]; ok && v[0] != f.Type {
				// 	return v[1] + "(" + n + ")"
				// }
				return n
			},
		}),
	}
}

func valueOf(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return valueOf(t.X)
	default:
		return expr
	}
}

func (c *codegen) Generate(filesrc string) error {
	format := RenameFunc(func(s string) string {
		return s
	})
	log.Println("File ->", filesrc)
	b, _ := os.ReadFile("./internal/templates/template.go.tpl")
	t, err := c.tmpl.Parse(string(b))
	if err != nil {
		return err
	}

	fset := token.NewFileSet() // positions are relative to fset
	pwd, _ := os.Getwd()

	// Parse src but stop after error.
	gofile, err := parser.ParseFile(fset, filesrc, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	// cache := make(map[string]*Entity)
	ent := Entity{}
	ent.GoPkg = types.ExprString(gofile.Name)

	c.importCache = make(map[string]string)
	for _, imp := range gofile.Imports {
		path, _ := strconv.Unquote(types.ExprString(imp.Path))
		// c.importCache[types.ExprString(imp.Name)] = path
		log.Println("Import ->", imp.Name, path)
	}

	for _, d := range gofile.Decls {
		// log.Println("decls ->", reflect.TypeOf(d))
		switch decl := d.(type) {
		// Determine which type implements driver.Scanner / driver.Valuer
		case *ast.FuncDecl:
			if len(decl.Recv.List) < 1 {
				continue
			}
			funcName := types.ExprString(decl.Name)
			noOfReturns := len(decl.Type.Results.List)
			if funcName == "Scan" && noOfReturns == 1 {

			}
			if funcName == "Value" && noOfReturns == 2 {

				// log.Println("Func ->", types.ExprString(decl))
				log.Println("FuncDeclRecv ->", decl.Recv.List[0])
				// log.Println(getType(decl.Recv.List[0].Type))
				log.Println("FuncDeclName ->", types.ExprString(valueOf(decl.Recv.List[0].Type)))
			}
			log.Println("FuncDeclType ->", len(decl.Type.Params.List))
			log.Println("FuncDeclResult ->")
		case *ast.GenDecl:

		}

		// If it's not IMPORT, CONST, TYPE, or VAR, we are not interested
		decl := assertAs[ast.GenDecl](d)
		if decl == nil {
			continue
		}

		// If the token isn't type, we are not interested
		if decl.Tok != token.TYPE {
			continue
		}

		for _, spec := range decl.Specs {
			// If it's not a `ast.TypeSpec`, we are not interested
			typeSpec := assertAs[ast.TypeSpec](spec)
			if typeSpec == nil {
				continue
			}

			// If it's not a `ast.StructType`, we are not interested
			structType := assertAs[ast.StructType](typeSpec.Type)
			if structType == nil {
				continue
			}

			// Set the struct name
			ent.GoName = typeSpec.Name.String()
			ent.Name = format(typeSpec.Name.String())

			for _, f := range structType.Fields.List {
				var (
					tag reflect.StructTag
				)

				if f.Tag != nil {
					// Trim backtick
					tag = reflect.StructTag(strings.TrimFunc(f.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				fieldType := c.getType(f.Type)
				log.Println("Type ->", fieldType, types.ExprString(f.Type))

				for _, n := range f.Names {
					// Skip private property
					if !n.IsExported() {
						continue
					}

					paths := strings.Split(tag.Get("sql"), ",")
					column := strings.TrimSpace(paths[0])
					// Skip if the property set skip
					if column == "-" {
						continue
					} else if column == "" {
						column = format(n.Name)
					}

					p := new(Field)
					p.GoName = n.Name
					p.Name = column
					p.Type = types.ExprString(f.Type)
					p.ActualType = fieldType

					log.Println("ActualType ->", fieldType, f.Type)

					ent.Fields = append(ent.Fields, p)

					packages.Load(&packages.Config{
						Mode: packages.NeedName |
							packages.NeedFiles |
							packages.NeedImports |
							packages.NeedTypes |
							packages.NeedSyntax |
							packages.NeedTypesInfo |
							packages.NeedModule |
							packages.NeedDeps,
					})
				}
			}
		}
		// ast.Inspect(f, func(n ast.Node) bool {
	}

	for _, f := range ent.Fields {
		log.Println(f)
	}

	filename := strings.Replace("{name}_gen.go", "{name}", strings.TrimSuffix(filepath.Base(filesrc), filepath.Ext(filesrc)), -1)
	filedst := filepath.Join(pwd, filepath.Dir(filesrc), filename)
	outfile, err := os.OpenFile(filedst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// outfile.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")

	if err := t.Execute(outfile, ent); err != nil {
		return err
	}
	return nil
}

func (c *codegen) getType(expr ast.Expr) string {
	var (
		actualType = "any"
	)

loop:
	for expr != nil {
		switch t := expr.(type) {
		// Base primitive type
		case *ast.Ident:
			if t.Obj == nil {
				actualType = t.String()
				break loop
			}
			typeSpec := assertAs[ast.TypeSpec](t.Obj.Decl)
			if typeSpec == nil {
				break loop
			}
			expr = typeSpec.Type
		// Array or Slice type
		case *ast.ArrayType:
			el := `[`
			if v := assertAs[ast.BasicLit](t.Len); v != nil {
				el += v.Value
			}
			el += `]`
			actualType = el + c.getType(t.Elt)
			break loop
		// Map type
		case *ast.MapType:
			actualType = `map[` + c.getType(t.Key) + `]` + c.getType(t.Value)
			break loop
		// Imported package
		case *ast.SelectorExpr:
			impPath := c.importCache[types.ExprString(t.X)]
			actualType = impPath + "." + c.getType(t.Sel)
			break loop
		// case *ast.InterfaceType:
		// case *ast.IndexExpr:
		default:
			break loop
		}
	}
	return actualType
}
