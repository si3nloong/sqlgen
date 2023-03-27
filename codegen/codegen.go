package codegen

import (
	"bytes"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/internal/strfmt"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/imports"
)

type RenameFunc func(string) string

type Generator struct {
	rename RenameFunc
}

var (
	sqlScanner, sqlValuer *types.Interface
)

func init() {
	var lconf loader.Config
	lconf.Import("database/sql")
	lprog, err := lconf.Load()
	if err != nil {
		panic(err)
	}

	sqlScanner = lprog.Package("database/sql").Pkg.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
	sqlValuer = lprog.Package("database/sql/driver").Pkg.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
}

func getUnderlyingType(t types.Type) string {
	switch v := t.(type) {
	case *types.Slice:
		return "[]" + v.Elem().Underlying().String()
	case *types.Named:
		return t.String()
	default:
		return t.Underlying().String()
	}
}

type Mapping struct {
	Encoder string
	Decoder string
}

type ImportPkgs []*types.Package

func (p *ImportPkgs) Load(imp *types.Package) bool {
	if i := slices.IndexFunc(*p, func(pkg *types.Package) bool {
		return pkg.Path() == imp.Path()
	}); i > -1 {
		return false
	}
	*p = append(*p, imp)
	return true
}

func Generate(cfg *config.Config) error {
	gen := new(Generator)
	gen.rename = strfmt.ToSnakeCase

	switch strings.ToLower(cfg.NamingConvention) {
	case "snakecase":
		gen.rename = strfmt.ToSnakeCase
	case "camelcase":
		gen.rename = strfmt.ToCamelCase
	}

	fset := token.NewFileSet() // positions are relative to fset
	// pwd, _ := os.Getwd()

	fileSrc := filepath.Dir(cfg.SrcDir)
	pkgs, err := parser.ParseDir(fset, fileSrc, func(fi fs.FileInfo) bool {
		filename := fi.Name()
		if strings.HasSuffix(filename, "_test.go") || strings.HasSuffix(filename, "_gen.go") || filename == "generated.go" {
			return false
		}
		return true
	}, parser.AllErrors)
	if err != nil {
		return err
	}

	structTypes := make(map[string]*ast.StructType)
	files := make([]*ast.File, 0)
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			files = append(files, f)
			ast.Inspect(f, func(node ast.Node) bool {
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if ok {
					structTypes[types.ExprString(typeSpec.Name)] = structType
				}
				return true
			})
		}

	}
	conf := &types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types:  make(map[ast.Expr]types.TypeAndValue),
		Defs:   make(map[*ast.Ident]types.Object),
		Uses:   make(map[*ast.Ident]types.Object),
		Scopes: make(map[ast.Node]*types.Scope),
	}

	pkg, err := conf.Check(fileSrc, fset, files, info)
	if err != nil {
		return err
	}

	impPkgs := make(ImportPkgs, 0)
	// checker := types.NewChecker(conf, fset, pkg, info)

	data := templates.ModelTmplParams{}
	data.GoPkg = pkg.Name()
	// log.Println(structs)

	// for k, _ := range info.Types {
	// 	log.Println(k, reflect.TypeOf(k))
	// }

	for k, st := range structTypes {
		var (
			model = new(templates.Model)
			queue = []*ast.StructType{st}
		)

		for len(queue) > 0 {
			s := queue[0]
			if len(s.Fields.List) == 0 {
				goto next
			}

			model.GoName = k
			model.Name = gen.rename(k)

			for _, f := range s.Fields.List {
				var tag reflect.StructTag
				if f.Tag != nil {
					// Trim backtick
					tag = reflect.StructTag(strings.TrimFunc(f.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				fi := f.Type.(*ast.Ident)
				// Check and process embedded struct
				if f.Names == nil && fi.Obj != nil {
					typeSpec, ok := fi.Obj.Decl.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					queue = append(queue, structType)
					continue
				}

				t := info.TypeOf(f.Type)
				typ := t.String()
				paths := strings.Split(tag.Get(cfg.Tag), ",")
				tagOpts := make(map[string]string)
				name := strings.TrimSpace(paths[0])
				for _, v := range paths[1:] {
					v = strings.ToLower(v)
					tagOpts[v] = v
				}
				if name == "-" {
					goto next
				} else if name != "" {
					if typ == "github.com/si3nloong/sqlgen/sql/schema.Name" {
						model.Name = name
						goto next
					}
				}

				for _, n := range f.Names {
					if !n.IsExported() {
						continue
					}

					field := new(templates.Field)
					field.GoName = types.ExprString(n)
					field.Type = t
					if name == "" {
						field.Name = gen.rename(field.GoName)
					} else {
						field.Name = name
					}

					if _, ok := tagOpts["pk"]; ok {
						model.PK = field
					}

					model.Fields = append(model.Fields, field)
				}
			}

		next:
			queue = queue[1:]
		}

		if len(model.Fields) == 0 {
			continue
		}

		data.Models = append(data.Models, model)
	}

	sort.Slice(data.Models, func(i, j int) bool {
		return data.Models[i].GoName < data.Models[j].GoName
	})

	typeMap := map[string]Mapping{
		"string":  {"string", "github.com/si3nloong/gqlgen/sql/types.String"},
		"[]byte":  {"string", "github.com/si3nloong/gqlgen/sql/types.String"},
		"bool":    {"bool", "github.com/si3nloong/gqlgen/sql/types.Bool"},
		"uint":    {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"uint8":   {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"uint16":  {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"uint32":  {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"uint64":  {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"int":     {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"int8":    {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"int16":   {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"int32":   {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"int64":   {"int64", "github.com/si3nloong/gqlgen/sql/types.Integer"},
		"float32": {"float64", "github.com/si3nloong/gqlgen/sql/types.Float"},
		"float64": {"float64", "github.com/si3nloong/gqlgen/sql/types.Float"},
	}

	tmpl := template.New("template.go").Funcs(template.FuncMap{
		"reserveImport": func(imp string, aliases ...string) string {
			name := filepath.Base(imp)
			if len(aliases) > 0 {
				name = aliases[0]
			}
			impPkgs.Load(types.NewPackage(imp, name))
			return ""
		},
		"isValuer": func(f *templates.Field) bool {
			_, ok := types.MissingMethod(f.Type, sqlValuer, true)
			return ok
		},
		"quote": strconv.Quote,
		"cast": func(n string, f *templates.Field) string {
			v := n + "." + f.GoName
			underType := f.Type.Underlying().String()
			if _, ok := types.MissingMethod(f.Type, sqlValuer, true); ok {
				impPkgs.Load(types.NewPackage("database/sql/driver", "driver"))
				return "(driver.Valuer)(" + v + ")"
			} else if typ, ok := typeMap[underType]; ok {
				return typ.Encoder + "(" + v + ")"
			}
			return v
		},
		"addr": func(n string, f *templates.Field) string {
			v := "&" + n + "." + f.GoName
			underType := f.Type.Underlying().String()
			if types.Implements(types.NewPointer(f.Type), sqlScanner) {
				impPkgs.Load(types.NewPackage("database/sql", "sql"))
				return "(sql.Scanner)(" + v + ")"
			} else if typ, ok := typeMap[underType]; ok {
				return filepath.Base(typ.Decoder) + "(" + v + ")"
			}
			return v
		},
	})

	fileDest := fileSrc + "/generated.go"
	w := bytes.NewBufferString("")
	// outFile, err := os.OpenFile(fileDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	return err
	// }
	w.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")
	w.WriteString("package " + data.GoPkg + "\n\n")
	b, _ := os.ReadFile("./codegen/templates/model.go.tpl")

	// log.Println(string(b))
	t, _ := tmpl.Parse(string(b))
	blr := bytes.NewBufferString(``)
	log.Println(data)
	if err := t.Execute(blr, data); err != nil {
		return err
	}

	if len(impPkgs) > 0 {
		w.WriteString("import (\n")
		for _, pkg := range impPkgs {
			w.WriteString(strconv.Quote(pkg.Path()) + "\n")
		}
		w.WriteString(")\n")
	}
	blr.WriteTo(w)

	log.Println(w.String())

	formatted, err := imports.Process(fileDest, w.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileDest, formatted, 0644); err != nil {
		return err
	}

	return nil
}
