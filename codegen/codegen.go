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
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/internal/strfmt"
	"github.com/si3nloong/sqlgen/sql/schema"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/imports"
)

var (
	schemaName = reflect.TypeOf(schema.Name{})
)

type RenameFunc func(string) string

type Generator struct {
	rename RenameFunc
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
	alias := p.newAliasIfExists(pkg)
	pkg.SetName(alias)
	p.cache[alias] = pkg
	p.importPkgs = append(p.importPkgs, pkg)
	return pkg, true
}

func (p *Package) newAliasIfExists(pkg *types.Package) string {
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
	}

	cfg.SrcDir = filepath.Dir(cfg.SrcDir)
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, cfg.SrcDir, func(fi fs.FileInfo) bool {
		filename := fi.Name()
		if strings.HasSuffix(filename, "_test.go") || strings.HasSuffix(filename, "_gen.go") || filename == "generated.go" {
			return false
		}
		return true
	}, parser.AllErrors)
	if err != nil {
		return err
	}

	if pkgs == nil {
		return nil
	}

	for k, pkg := range pkgs {
		if err := parsePackage(fset, pkg, cfg); err != nil {
			log.Println(err)
		}
		delete(pkgs, k)
	}

	return nil
}

func parsePackage(fset *token.FileSet, pkg *ast.Package, cfg *config.Config) error {
	fileSrc := cfg.SrcDir
	files := make([]*ast.File, 0)
	structTypes := make(map[string]*ast.StructType)

	for _, f := range pkg.Files {
		files = append(files, f)

		ast.Inspect(f, func(node ast.Node) bool {
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			// TODO: If it's an alias struct, we should skip right?
			structType, ok := typeSpec.Type.(*ast.StructType)
			if ok {
				structTypes[types.ExprString(typeSpec.Name)] = structType
			}
			return true
		})
	}

	conf := &types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types:  make(map[ast.Expr]types.TypeAndValue),
		Defs:   make(map[*ast.Ident]types.Object),
		Uses:   make(map[*ast.Ident]types.Object),
		Scopes: make(map[ast.Node]*types.Scope),
	}

	typePkg, err := conf.Check(fileSrc, fset, files, info)
	if err != nil {
		return err
	}

	impPkgs := new(Package)
	data := templates.ModelTmplParams{}
	data.GoPkg = typePkg.Name()

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
			model.Name = strfmt.ToSnakeCase(k)

			for _, f := range s.Fields.List {
				var tag reflect.StructTag
				if f.Tag != nil {
					// Trim backtick
					tag = reflect.StructTag(strings.TrimFunc(f.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				switch vi := ElemOf(f.Type).(type) {
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
						field.Name = strfmt.ToSnakeCase(field.GoName)
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
		"cast": func(n string, f *templates.Field) string {
			v := n + "." + f.GoName
			underType := getUnderlyingType(f.Type)
			if IsImplemented(f.Type, sqlValuer) {
				p, _ := impPkgs.Import(valuerPkg)
				return "(" + p.Name() + ".Valuer)(" + v + ")"
			}
			if typ, ok := typeMap[underType]; ok {
				if string(typ.Encoder) == f.Type.String() {
					return v
				}
				return typ.Encoder.CastOrInvoke(impPkgs, v)
			}
			return v
		},
		"addr": func(n string, f *templates.Field) string {
			v := "&" + n + "." + f.GoName
			underType := getUnderlyingType(f.Type)
			if types.Implements(types.NewPointer(f.Type), sqlScanner) {
				p, _ := impPkgs.Import(scannerPkg)
				return "(" + p.Name() + ".Scanner)(" + v + ")"
			}
			if typ, ok := typeMap[underType]; ok {
				if string(typ.Encoder) == f.Type.String() {
					return v
				}
				return typ.Decoder.CastOrInvoke(impPkgs, v)
			}
			return v
		},
	})

	b, err := os.ReadFile(filepath.Join(fileutil.CurDir(), "templates/model.go.tpl"))
	if err != nil {
		return err
	}

	t, err := tmpl.Parse(string(b))
	if err != nil {
		return err
	}

	blr := bytes.NewBufferString("")
	if err := t.Execute(blr, data); err != nil {
		return err
	}

	w := bytes.NewBufferString("")
	if cfg.IncludeHeader != nil && *cfg.IncludeHeader {
		w.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")
	}

	w.WriteString("package " + data.GoPkg + "\n\n")

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

	fileDest := filepath.Join(fileSrc, "generated.go")
	formatted, err := imports.Process(fileDest, w.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileDest, formatted, 0o644); err != nil {
		return err
	}
	return nil
}

func getUnderlyingType(t types.Type) string {
	typeStr := ""
	for t != nil {
		switch v := t.(type) {
		case *types.Pointer:
			typeStr += "*"
			t = v.Elem()
		case *types.Slice:
			typeStr += "[]"
			t = v.Elem()
		case *types.Map:
			typeStr += "map[" + getUnderlyingType(v.Key()) + "]"
			t = v.Elem()
		case *types.Named:
			switch vt := t.Underlying().(type) {
			case *types.Basic:
				return typeStr + vt.String()
			default:
				return typeStr + t.String()
			}
		default:
			return typeStr + t.Underlying().String()
		}
	}
	return typeStr
}

func IsImplemented(t types.Type, iv *types.Interface) bool {
	_, ok := types.MissingMethod(t, iv, true)
	return ok
}

func ElemOf(expr ast.Expr) ast.Expr {
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
