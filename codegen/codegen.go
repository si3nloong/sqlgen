package codegen

import (
	"bytes"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
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
	"github.com/si3nloong/sqlgen/sql/schema"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/imports"
)

var (
	schemaName = reflect.TypeOf(schema.Name{})

	typeMap = map[string]Mapping{
		"string":     {"string", "github.com/si3nloong/sqlgen/sql/types.String"},
		"[]byte":     {"string", "github.com/si3nloong/sqlgen/sql/types.String"},
		"bool":       {"bool", "github.com/si3nloong/sqlgen/sql/types.Bool"},
		"uint":       {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint8":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint16":     {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint32":     {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint64":     {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int":        {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int8":       {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int16":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int32":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int64":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"float32":    {"float64", "github.com/si3nloong/sqlgen/sql/types.Float"},
		"float64":    {"float64", "github.com/si3nloong/sqlgen/sql/types.Float"},
		"*string":    {"github.com/si3nloong/sqlgen/sql/types.String", "github.com/si3nloong/sqlgen/sql/types.PtrOfString"},
		"*[]byte":    {"github.com/si3nloong/sqlgen/sql/types.String", "github.com/si3nloong/sqlgen/sql/types.PtrOfString"},
		"*bool":      {"github.com/si3nloong/sqlgen/sql/types.Bool", "github.com/si3nloong/sqlgen/sql/types.PtrOfBool"},
		"*uint":      {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint8":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint16":    {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint32":    {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint64":    {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int":       {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int8":      {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int16":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int32":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int64":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*float32":   {"github.com/si3nloong/sqlgen/sql/types.Float", "github.com/si3nloong/sqlgen/sql/types.PtrOfFloat"},
		"*float64":   {"github.com/si3nloong/sqlgen/sql/types.Float", "github.com/si3nloong/sqlgen/sql/types.PtrOfFloat"},
		"*time.Time": {"github.com/si3nloong/sqlgen/sql/types.Time", "github.com/si3nloong/sqlgen/sql/types.PtrOfTime"},
	}
)

type RenameFunc func(string) string

type Generator struct {
	rename RenameFunc
}

type Codec string

func (c Codec) IsPkgFunc() (*types.Package, string, bool) {
	pkg := string(c)
	idx := strings.LastIndexByte(pkg, '.')
	if idx > 0 {
		path := pkg[:idx]
		cb := pkg[idx+1:]
		return types.NewPackage(path, filepath.Base(path)), cb, true
	}
	return nil, "", false
}

func (c Codec) CastOrInvoke(pkg *Package, v string) string {
	if p, invoke, ok := c.IsPkgFunc(); ok {
		p, _ = pkg.Import(p)
		return p.Name() + "." + invoke + "(" + v + ")"
	}
	return string(c) + "(" + v + ")"
}

type Mapping struct {
	Encoder Codec
	Decoder Codec
}

type Package struct {
	cache      map[string]*types.Package
	importPkgs []*types.Package
}

func (p *Package) Import(pkg *types.Package) (*types.Package, bool) {
	if i := slices.IndexFunc(p.importPkgs, func(item *types.Package) bool {
		return pkg.Path() == item.Path()
	}); i > -1 {
		return p.importPkgs[i], false
	}
	if p.cache == nil {
		p.cache = make(map[string]*types.Package)
	}
	alias := p.newAliasIfRequired(pkg)
	pkg.SetName(alias)
	p.cache[alias] = pkg
	p.importPkgs = append(p.importPkgs, pkg)
	return pkg, true
}

func (p *Package) newAliasIfRequired(pkg *types.Package) string {
	pkgName, newPkgName := pkg.Name(), pkg.Name()
	for i := 1; ; i++ {
		if _, ok := p.cache[newPkgName]; ok {
			newPkgName = pkgName + strconv.Itoa(i)
			continue
		}
		break
	}
	return newPkgName
}

func Generate(cfg *config.Config) error {
	gen := new(Generator)
	gen.rename = strfmt.ToSnakeCase

	switch strings.ToLower(cfg.NamingConvention) {
	case "snakecase":
		gen.rename = strfmt.ToSnakeCase
	case "camelcase":
		gen.rename = strfmt.ToCamelCase
	case "-":
		gen.rename = func(s string) string {
			return s
		}
	}

	fileSrc := filepath.Dir(cfg.SrcDir)
	fset := token.NewFileSet()

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

	impPkgs := new(Package)
	data := templates.ModelTmplParams{}
	data.GoPkg = pkg.Name()

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

				// fi := valueOf(f.Type).(*ast.Ident)
				switch vi := valueOf(f.Type).(type) {
				case *ast.Ident:
					// Check and process embedded struct
					if f.Names == nil && vi.Obj != nil {
						typeSpec, ok := vi.Obj.Decl.(*ast.TypeSpec)
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
					continue
				} else if name != "" {
					log.Println(typ, schemaName.PkgPath()+"."+schemaName.Name())
					if typ == schemaName.PkgPath()+"."+schemaName.Name() {
						model.Name = name
						continue
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

	log.Println(data)

	sort.Slice(data.Models, func(i, j int) bool {
		return data.Models[i].GoName < data.Models[j].GoName
	})

	tmpl := template.New("template.go").Funcs(template.FuncMap{
		"quote": strconv.Quote,
		"reserveImport": func(pkgPath string, aliases ...string) string {
			name := filepath.Base(pkgPath)
			if len(aliases) > 0 {
				name = aliases[0]
			}
			impPkgs.Import(types.NewPackage(pkgPath, name))
			return ""
		},
		"isValuer": func(f *templates.Field) bool {
			return IsImplemented(f.Type, sqlValuer)
		},
		"cast": func(n string, f *templates.Field) (string, error) {
			v := n + "." + f.GoName
			underType := getUnderlyingType(f.Type)
			if IsImplemented(f.Type, sqlValuer) {
				p, _ := impPkgs.Import(valuerPkg)
				return "(" + p.Name() + ".Valuer)(" + v + ")", nil
			}
			if typ, ok := typeMap[underType]; ok {
				return typ.Encoder.CastOrInvoke(impPkgs, v), nil
			}
			return v, nil
		},
		"addr": func(n string, f *templates.Field) string {
			v := "&" + n + "." + f.GoName
			underType := getUnderlyingType(f.Type)
			if types.Implements(types.NewPointer(f.Type), sqlScanner) {
				p, _ := impPkgs.Import(scannerPkg)
				return "(" + p.Name() + ".Scanner)(" + v + ")"
			}
			if typ, ok := typeMap[underType]; ok {
				return typ.Decoder.CastOrInvoke(impPkgs, v)
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

	if len(impPkgs.importPkgs) > 0 {
		w.WriteString("import (\n")
		for _, pkg := range impPkgs.importPkgs {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				w.WriteString(strconv.Quote(pkg.Path()) + "\n")
			} else {
				w.WriteString(pkg.Name() + " " + strconv.Quote(pkg.Path()) + "\n")
			}
		}
		w.WriteString(")\n")
	}
	blr.WriteTo(w)

	log.Println(w.String())

	formatted, err := imports.Process(fileDest, w.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		log.Println(err)
		return err
	}

	if err := os.WriteFile(fileDest, formatted, 0644); err != nil {
		return err
	}

	return nil
}

func getUnderlyingType(t types.Type) string {
	switch v := t.(type) {
	case *types.Slice:
		return "[]" + getUnderlyingType(v.Elem())
	case *types.Named:
		return t.Underlying().String()
	case *types.Pointer:
		return "*" + getUnderlyingType(v.Elem())
	default:
		return t.Underlying().String()
	}
}

func IsImplemented(t types.Type, iv *types.Interface) bool {
	_, ok := types.MissingMethod(t, iv, true)
	return ok
}

func valueOf(expr ast.Expr) ast.Expr {
	for expr != nil {
		switch v := expr.(type) {
		case *ast.StarExpr:
			expr = v.X
		default:
			return v
		}
	}
	return expr
}
