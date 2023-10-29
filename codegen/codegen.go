package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/internal/gosyntax"
	"github.com/si3nloong/sqlgen/internal/strfmt"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/valyala/bytebufferpool"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type tagOption string

const (
	TagOptionAutoIncrement tagOption = "auto_increment"
	TagOptionBinary        tagOption = "binary"
	TagOptionPKAlias       tagOption = "pk"
	TagOptionPK            tagOption = "primary_key"
	TagOptionSize          tagOption = "size"
	TagOptionDataType      tagOption = "datatype"
	TagOptionUnique        tagOption = "unique"
)

var (
	schemaName      = reflect.TypeOf(sequel.Name{})
	tableNameSchema = schemaName.PkgPath() + "." + schemaName.Name()
)

type RenameFunc func(string) string

type typeQueue struct {
	path string
	idx  []int
	t    *ast.StructType
}

type Generator struct {
	rename RenameFunc
}

type structField struct {
	id       string
	name     string
	path     string
	t        ast.Expr
	exported bool
	tag      reflect.StructTag
}

type tagOpts map[string]string

func (t tagOpts) Lookup(key tagOption, keys ...tagOption) (v string, ok bool) {
	keys = append(keys, key)
	for k, v := range t {
		if lo.IndexOf(keys, tagOption(k)) >= 0 {
			return v, true
		}
	}
	return
}

func Generate(cfg *config.Config) error {
	// gen := new(Generator)
	rename := strfmt.ToSnakeCase

	switch strings.ToLower(cfg.NamingConvention) {
	case "snakecase":
		rename = strfmt.ToSnakeCase
	case "camelcase":
		rename = strfmt.ToCamelCase
	case "no":
		rename = func(s string) string { return s }
	}

	log.Println(cfg)
	fi, err := os.Stat(cfg.SrcDir)
	if err != nil {
		return err
	}

	dir := cfg.SrcDir
	if !fi.IsDir() {
		dir = filepath.Join(fileutil.Getpwd(), filepath.Dir(cfg.SrcDir))
	}
	if cfg.SrcDir == "." {
		dir = fileutil.Getpwd()
	}

	// For example:
	// docs/*.md
	// docs/**/*
	//
	// Get all the folders, filter files
	//
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		return nil
	})

	filename := path.Join(dir, "generated.go")
	os.Remove(filename) // remove "generated.go"

	pkgs, err := packages.Load(&packages.Config{
		Dir:  dir,
		Mode: mode,
	})
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return nil
	}

	var (
		pkg         = pkgs[0]
		impPkgs     = new(Package)
		structTypes = make(map[*ast.TypeSpec]*ast.StructType)
	)

	for _, f := range pkg.Syntax {
		// ast.Print(pkg.Fset, f)
		ast.Inspect(f, func(node ast.Node) bool {
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}

			obj := pkg.TypesInfo.ObjectOf(typeSpec.Name)
			// We're not interested in the unexported type
			if !obj.Exported() {
				return true
			}

			// TODO: If it's an alias struct, we should skip right?
			structType, ok := typeSpec.Type.(*ast.StructType)
			if ok {
				structTypes[typeSpec] = structType
			}
			return true
		})
	}

	var (
		structs = make(map[*ast.Ident][]*structField, 0)
	)

	// Loop every struct and map the fields
	for k, s := range structTypes {
		var (
			// queue to store struct, this is useful
			// when handling embedded struct
			q      = []typeQueue{{t: s}}
			f      typeQueue
			fields = make([]*structField, 0)
		)

		for len(q) > 0 {
			f = q[0]

			// If the struct has empty field, just skip
			if len(f.t.Fields.List) == 0 {
				goto next
			}

			// Loop every struct field
			for i, fi := range f.t.Fields.List {
				var tag reflect.StructTag
				if fi.Tag != nil {
					// Trim backtick
					tag = reflect.StructTag(strings.TrimFunc(fi.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				// If the field is embedded struct
				// `Type` can be either *ast.Ident or *ast.SelectorExpr
				if fi.Names == nil {
					switch vi := fi.Type.(type) {
					// Local struct
					case *ast.Ident:
						path := types.ExprString(vi)
						if f.path != "" {
							path = f.path + "." + path
						}
						t := vi.Obj.Decl.(*ast.TypeSpec)
						q = append(q, typeQueue{path: path, idx: append(f.idx, i), t: t.Type.(*ast.StructType)})
						continue

					// Imported struct
					case *ast.SelectorExpr:
						log.Println(vi)
					}
				}

				for j, n := range fi.Names {
					path := types.ExprString(n)
					if f.path != "" {
						path = f.path + "." + path
					}

					fields = append(fields, &structField{id: toID(append(f.idx, i+j)), exported: n.IsExported(), name: types.ExprString(n), tag: tag, path: path, t: fi.Type})
				}
			}

		next:
			q = q[1:]
		}

		if len(fields) > 0 {
			sort.Slice(fields, func(i, j int) bool {
				return fields[j].id > fields[i].id
			})

			structs[k.Name] = fields
		}
	}

	// Generate interface code
	var (
		nameMap map[string]struct{}
		t       types.Type
		params  = &templates.ModelTmplParams{}
	)

	// Convert struct to models and generate code
	for n, s := range structs {
		t = pkg.TypesInfo.TypeOf(n)
		nameMap = make(map[string]struct{})

		var (
			index int
			model = templates.Model{}
		)
		model.GoName = types.ExprString(n)
		model.TableName = rename(model.GoName)
		model.HasTableName = !IsImplemented(t, sqlTabler)
		model.HasColumn = !IsImplemented(t, sqlColumner)
		// model.HasRow = !IsImplemented(t, sqlRower)

		for _, f := range s {
			tv := f.tag.Get("sql")

			switch pkg.TypesInfo.TypeOf(f.t).String() {
			// If the type is table name
			case tableNameSchema:
				model.TableName = ""
				continue
			}

			// If it's a unexported field, skip!
			if !f.exported {
				continue
			}

			tf := &templates.Field{}
			tf.ColumnName = rename(f.name)
			tf.Type = pkg.TypesInfo.TypeOf(f.t)
			tag := make(tagOpts)

			if tv != "" {
				tags := strings.Split(tv, ",")
				name := strings.TrimSpace(tags[0])
				if name == "-" {
					continue
				} else if name != "" {
					tf.ColumnName = name
				}
				for _, v := range tags[1:] {
					kv := strings.SplitN(v, ":", 2)
					k := strings.TrimSpace(strings.ToLower(kv[0]))
					if len(kv) > 1 {
						tag[k] = kv[1]
					} else {
						tag[k] = ""
					}
				}
			}

			tf.GoName = f.name
			tf.GoPath = f.path
			tf.Index = index
			_, tf.IsBinary = tag.Lookup(TagOptionBinary)
			if v, ok := tag.Lookup(TagOptionSize); ok {
				tf.Size, _ = strconv.Atoi(v)
			}
			index++

			if _, ok := tag.Lookup(TagOptionPK, TagOptionPKAlias, TagOptionAutoIncrement); ok {
				if model.PK != nil {
					return fmt.Errorf(`sqlgen: a model can only allow one primary key, else it will get overriden`)
				}

				// Check auto increment
				pk := templates.PK{Field: tf}
				_, pk.IsAutoIncr = tag.Lookup(TagOptionAutoIncrement)
				model.PK = &pk
			}

			// Check uniqueness of the column
			if _, ok := nameMap[tf.ColumnName]; ok {
				return fmt.Errorf("sqlgen: duplicate key %s in struct", tf.ColumnName)
			}
			nameMap[tf.ColumnName] = struct{}{}

			// Check type is sequel.Name, then override name
			model.Fields = append(model.Fields, tf)
		}

		clear(nameMap)
		params.Models = append(params.Models, &model)
	}

	tmpl := template.New(dir).Funcs(template.FuncMap{
		"quote":         strconv.Quote,
		"createTable":   createTableStmt,
		"alterTable":    alterTableStmt,
		"reserveImport": reserveImport(impPkgs),
		"castAs":        castAs(impPkgs),
		"addrOf":        addrOf(impPkgs),
	})

	tmpFile := filepath.Join(fileutil.CurDir(), "templates/model.gtpl")
	b, err := os.ReadFile(tmpFile)
	if err != nil {
		return err
	}

	tt, err := tmpl.Parse(string(b))
	if err != nil {
		return err
	}

	blr := bytes.NewBufferString("")
	if err := tt.Execute(blr, params); err != nil {
		return err
	}

	w := bytes.NewBufferString("")
	if cfg.IncludeHeader {
		w.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")
	}

	w.WriteString("package " + pkg.Name + "\n\n")

	if len(impPkgs.importPkgs) > 0 {
		w.WriteString("import (\n")
		for _, pkg := range impPkgs.importPkgs {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				w.WriteString("\t" + strconv.Quote(pkg.Path()) + "\n")
			} else {
				w.WriteString("\t" + pkg.Name() + " " + strconv.Quote(pkg.Path()) + "\n")
			}
		}
		w.WriteString(")\n")
	}

	w.WriteString(blr.String())
	// log.Println(w.String())

	fileDest := filepath.Join(dir, "generated.go")
	formatted, err := imports.Process(fileDest, w.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileDest, formatted, 0o644); err != nil {
		return err
	}

	return goModTidy()
}

func (g *Generator) parsePackage(fset *token.FileSet, files []*ast.File, cfg *config.Config) error {
	fileSrc := filepath.Dir(cfg.SrcDir)
	// files := make([]*ast.File, 0)
	structTypes := make(map[string]*ast.StructType)

	for _, f := range files {
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

	// conf := &types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	info := &types.Info{
		Scopes: make(map[ast.Node]*types.Scope),
		Types:  make(map[ast.Expr]types.TypeAndValue),
		Defs:   make(map[*ast.Ident]types.Object),
		Uses:   make(map[*ast.Ident]types.Object),
	}

	impPkgs := new(Package)
	data := templates.ModelTmplParams{}

	for k, st := range structTypes {
		var (
			model = new(templates.Model)
			queue = []*ast.StructType{st}
		)

		for len(queue) > 0 {
			s := queue[0]
			index := uint(0)
			if len(s.Fields.List) == 0 {
				goto next
			}

			model.GoName = k
			model.TableName = g.rename(k)

			for _, f := range s.Fields.List {
				var tag reflect.StructTag
				if f.Tag != nil {
					// Trim backtick
					tag = reflect.StructTag(strings.TrimFunc(f.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				switch vi := gosyntax.ElemOf(f.Type).(type) {
				// Check and process the embedded struct
				case *ast.Ident:
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
				tagOpts := make([]string, 0, len(paths))
				name := strings.TrimSpace(paths[0])
				for _, v := range paths[1:] {
					tagOpts = append(tagOpts, strings.ToLower(v))
				}

				if name == "-" {
					continue
				} else if name != "" {
					if typ == schemaName.PkgPath()+"."+schemaName.Name() {
						model.GoName = name
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
					// field.Tag = tagOpts
					if name == "" {
						field.GoName = g.rename(field.GoName)
					} else {
						field.GoName = name
					}

					index++
					if slices.Index(tagOpts, "pk") < 0 {
						model.Fields = append(model.Fields, field)
						continue
					}

					// model.PK = field
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
		"wrap": func(str string) string {
			return "`" + str + "`"
		},
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
			actualType := gosyntax.UnderlyingType(f.Type)
			if IsImplemented(f.Type, sqlValuer) {
				p, _ := impPkgs.Import(types.NewPackage("", ""))
				return "(" + p.Name() + ".Valuer)(" + v + ")"
			}
			if typ, ok := typeMap[actualType]; ok {
				if string(typ.Encoder) == f.Type.String() {
					return v
				}
				return typ.Encoder.Format(impPkgs, v)
			}
			return v
		},
		"addr": func(n string, f *templates.Field) string {
			v := "&" + n + "." + f.GoName
			actualType := gosyntax.UnderlyingType(f.Type)
			if types.Implements(types.NewPointer(f.Type), sqlScanner) {
				p, _ := impPkgs.Import(types.NewPackage("", ""))
				return "(" + p.Name() + ".Scanner)(" + v + ")"
			}
			if typ, ok := typeMap[actualType]; ok {
				if string(typ.Encoder) == f.Type.String() {
					return v
				}
				return typ.Decoder.Format(impPkgs, v)
			}
			return v
		},
	})

	tmpFile := filepath.Join(fileutil.CurDir(), "templates/model.gtpl")
	b, err := os.ReadFile(tmpFile)
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
	if cfg.IncludeHeader {
		w.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")
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

	fileDest := filepath.Join(fileSrc, "generated.go")
	formatted, err := imports.Process(fileDest, w.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileDest, formatted, 0o644); err != nil {
		return err
	}

	return g.goModTidy()
}

func (g *Generator) goModTidy() error {
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stdout
	return tidyCmd.Run()
}

func IsImplemented(t types.Type, iv *types.Interface) bool {
	method, wrongType := types.MissingMethod(t, iv, true)
	return method == nil && !wrongType
}

func goModTidy() error {
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stdout
	return tidyCmd.Run()
}

func toID(val []int) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	for i, v := range val {
		if i > 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(strconv.Itoa(v))
	}
	return buf.String()
}

func UnderlyingType(t types.Type) (*Mapping, bool) {
	var (
		typeStr string
		prev    = t
	)
	for t != nil {
		switch v := t.(type) {
		case *types.Basic:
			typeStr = v.String()
			prev = t.Underlying()
		case *types.Named:
			typeStr = v.String()
			prev = t.Underlying()
		case *types.Pointer:
			typeStr = v.Underlying().String()
			prev = v.Elem()
		default:
			break
		}
		if v, ok := typeMap[typeStr]; ok {
			return v, ok
		}
		if prev == t {
			break
		}
		t = prev
	}
	return nil, false
}
