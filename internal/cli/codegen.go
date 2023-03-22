package cli

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
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
}

type RenameFunc func(string) string

func getType(expr ast.Expr) string {
	var (
		actualType = "any"
	)

loop:
	for expr != nil {
		log.Println(reflect.TypeOf(expr))
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
			actualType = el + getType(t.Elt)
			break loop
		// Map type
		case *ast.MapType:
			actualType = `map[` + getType(t.Key) + `]` + getType(t.Value)
			break loop
		// case *ast.SelectorExpr:
		// case *ast.InterfaceType:
		case *ast.IndexExpr:
		default:
			break loop
		}
	}
	return actualType
}

type codegen struct {
	// Go template engine
	tmpl *template.Template
}

func newCodegen() *codegen {
	return &codegen{
		tmpl: template.New("codegen").Funcs(template.FuncMap{
			"quote": strconv.Quote,
			"cast": func(n string, f Field) string {
				if v, ok := sqlTypes[f.Type]; ok && v[0] != f.Type {
					return v[0] + "(" + n + ")"
				} else if f.BaseType != f.Type {
					return f.Type + "(" + n + ")"
				}
				return n
			},
			"addr": func(n string, f Field) string {
				if v, ok := sqlTypes[f.Type]; ok && v[0] != f.Type {
					return v[1] + "(" + n + ")"
				}
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

func (cg *codegen) run(filesrc string) error {
	// cli := "xxx"
	rename := RenameFunc(func(s string) string {
		return s
	})
	b, _ := os.ReadFile("./internal/template/template.go.tmpl")
	t, err := cg.tmpl.Parse(string(b))
	if err != nil {
		return err
	}

	fset := token.NewFileSet() // positions are relative to fset
	pwd, _ := os.Getwd()
	// src, _ := os.ReadFile(filesrc)
	// Parse src but stop after processing the imports.
	gofile, err := parser.ParseFile(fset, filesrc, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	if gofile.Name == nil {
		return errors.New(`missing go package name`)
	}

	// cache := make(map[string]*Entity)
	ent := Entity{}
	ent.Pkg = types.ExprString(gofile.Name)

	log.Println(len(gofile.Decls), gofile.Decls)
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
				log.Println(getType(decl.Recv.List[0].Type))
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

		log.Println("token ->", decl.Tok.String())
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

			log.Println("t ->", reflect.TypeOf(typeSpec.Type))
			// If it's not a `ast.StructType`, we are not interested
			structType := assertAs[ast.StructType](typeSpec.Type)
			if structType == nil {
				continue
			}

			// Set the struct name
			ent.Name = typeSpec.Name.String()

			for _, f := range structType.Fields.List {
				var (
					tag reflect.StructTag
				)

				if f.Tag != nil {
					tag = reflect.StructTag(strings.TrimFunc(f.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				// f.Type.Pos()
				log.Println("Type ->", f.Type, types.ExprString(f.Type))
				fieldType := getType(f.Type)

				for _, n := range f.Names {
					if !n.IsExported() {
						continue
					}

					paths := strings.Split(tag.Get("sql"), ",")
					column := strings.TrimSpace(paths[0])
					// Skip if the property set skip
					if column == "-" {
						continue
					} else if column == "" {
						column = rename(n.Name)
					}

					p := new(Field)
					p.Name = n.Name
					p.Column = column
					p.BaseType = types.ExprString(f.Type)
					p.Type = fieldType

					ent.FieldList = append(ent.FieldList, p)

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

	// for _, f := range ent.FieldList {
	// 	log.Println(f)
	// }

	ext := filepath.Ext(filesrc)
	filename := strings.Replace("{name}_gen.go", "{name}", strings.TrimSuffix(filepath.Base(filesrc), ext), -1)
	filedst := filepath.Join(pwd, filepath.Dir(filesrc), filename)
	outfile, err := os.OpenFile(filedst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// outfile.WriteString("// Code generated by " + cli + ", version 1.0.0. DO NOT EDIT.\n\n")

	if err := t.Execute(outfile, ent); err != nil {
		return err
	}
	return nil
}
