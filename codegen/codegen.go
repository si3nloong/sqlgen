package codegen

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/types"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/valyala/bytebufferpool"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

//go:embed templates/*.go.tpl
var codegenTemplates embed.FS

const fileMode = 0o755

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

type typeQueue struct {
	path string
	idx  []int
	t    *ast.StructType
}

type structType struct {
	name   *ast.Ident
	fields []structField
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

var path2regex = strings.NewReplacer(
	`.`, `\.`,
	`*`, `.*`,
	`\`, `[\\/]`,
	`/`, `[\\/]`,
)

func Generate(c *config.Config) error {
	var (
		srcDir  string
		cfg     = c.Clone()
		sources = make([]string, len(cfg.Source))
	)

	log.Println(cfg)
	copy(sources, cfg.Source)

	// Resolve every source provided
	for len(sources) > 0 {
		srcDir = strings.TrimSpace(sources[0])
		log.Println(srcDir)
		if srcDir == "" {
			return fmt.Errorf(`sqlgen: src is empty path`)
		}

		// If the prefix is ".", mean it's refer to current directory
		if srcDir[0] == '.' {
			srcDir = fileutil.Getpwd() + srcDir[1:]
		}

		// File: examples/testdata/test.go
		// Folder: examples/testdata
		// Wildcard: [examples/**, examples/testdata/**/*.go,  examples/testdata/**/*]
		// File wildcard: [examples/testdata/*model.go, examples/testdata/*_model.go]
		var (
			rootDir    string
			r                  = regexp.MustCompile(`(?i)((?:\/)([a-z][a-z0-9-_.]+\/)*)\w*\*\w*(?:\.go)`)
			subMatches         = r.FindStringSubmatch(srcDir)
			matcher    Matcher = new(EmptyMatcher)
			dirs               = make([]string, 0)
		)

		if strings.Contains(srcDir, "**") {
			paths := strings.SplitN(srcDir, "**", 2)
			rootDir = strings.TrimSuffix(strings.TrimSpace(paths[0]), "/")
			suffix := `(?:[\\/]\w+\.\w+)`
			if paths[1] != "" {
				suffix = path2regex.Replace(paths[1])
			}
			if err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
				// If the directory is not exists, the "d" will be nil
				if d == nil || !d.IsDir() {
					// If it's not a folder, we skip!
					return nil
				}
				dirs = append(dirs, strings.TrimPrefix(path, rootDir))
				return nil
			}); err != nil {
				return fmt.Errorf(`sqlgen: failed to walk schema %s: %w`, paths[0], err)
			}
			matcher = &RegexMatcher{regexp.MustCompile(path2regex.Replace(rootDir) + `([\\/][a-z0-9_-]+)*` + suffix)}
		} else if len(subMatches) > 0 {
			matcher = &RegexMatcher{regexp.MustCompile(path2regex.Replace(srcDir))}
			rootDir = strings.TrimSuffix(subMatches[1], "/")
			dirs = append(dirs, "")
		} else {
			fi, err := os.Stat(srcDir)
			// If the file or folder not exists, we skip!
			if os.IsNotExist(err) {
				goto nextSrc
			} else if err != nil {
				return err
			}

			// If it's just a file
			if !fi.IsDir() {
				srcDir = filepath.Dir(srcDir)
				matcher = FileMatcher{filepath.Join(srcDir, fi.Name()): struct{}{}}
			}

			rootDir = srcDir
			dirs = append(dirs, "")
		}

		if err := parseGoPackage(cfg, rootDir, dirs, matcher); err != nil {
			return err
		}

	nextSrc:
		sources = sources[1:]
	}

	return goModTidy()
}

func parseGoPackage(cfg *config.Config, rootDir string, dirs []string, matcher Matcher) error {
	var (
		rename   = cfg.RenameFunc()
		dialect  = cfg.Dialect()
		dir      string
		filename string
		w, blr   = bytes.NewBufferString(""), bytes.NewBufferString("")
	)

	for len(dirs) > 0 {
		dir = path.Join(rootDir, dirs[0])
		if fileutil.IsDirEmptyFiles(dir, cfg.Exec.Filename) {
			dirs = dirs[1:]
			continue
		}

		filename = path.Join(dir, cfg.Exec.Filename)
		// Remove generated file, ignore the error
		os.Remove(filename)

		pkgs, err := packages.Load(&packages.Config{
			Dir:  dir,
			Mode: mode,
		})
		if err != nil {
			return err
		} else if len(pkgs) == 0 {
			return nil
		}

		var (
			pkg = pkgs[0]
		)

		if len(pkg.Errors) > 0 {
			return pkg.Errors[0]
		}

		var (
			impPkg      = new(Package)
			structTypes = make([]*ast.TypeSpec, 0)
		)

		for _, f := range pkg.Syntax {
			// ast.Print(pkg.Fset, f)
			ast.Inspect(f, func(node ast.Node) bool {
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				filename = pkg.Fset.Position(typeSpec.Name.NamePos).Filename
				if !matcher.Match(filename) {
					return true
				}

				obj := pkg.TypesInfo.ObjectOf(typeSpec.Name)
				// We're not interested in the unexported type
				if !obj.Exported() {
					return true
				}

				// TODO: If it's an alias struct, we should skip right?
				if _, ok := typeSpec.Type.(*ast.StructType); ok {
					structTypes = append(structTypes, typeSpec)
				}
				return true
			})
		}

		// If we want to preserve the ordering,
		// we must use array instead of map
		var (
			structs = make([]structType, 0)
		)

		// Loop every struct and map the fields
		for _, s := range structTypes {
			var (
				// queue to store struct, this is useful
				// when handling embedded struct
				q      = []typeQueue{{t: s.Type.(*ast.StructType)}}
				f      typeQueue
				fields = make([]structField, 0)
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
							// Object can be nil
							if vi.Obj == nil {
								continue
							}
							path := types.ExprString(vi)
							if f.path != "" {
								path = f.path + "." + path
							}
							t := vi.Obj.Decl.(*ast.TypeSpec)
							q = append(q, typeQueue{path: path, idx: append(f.idx, i), t: t.Type.(*ast.StructType)})
							continue

						// Imported struct
						case *ast.SelectorExpr:
							// log.Println(vi)
						}
					}

					for j, n := range fi.Names {
						path := types.ExprString(n)
						if f.path != "" {
							path = f.path + "." + path
						}

						fields = append(fields, structField{id: toID(append(f.idx, i+j)), exported: n.IsExported(), name: types.ExprString(n), tag: tag, path: path, t: fi.Type})
					}
				}

			next:
				q = q[1:]
			}
			if len(fields) > 0 {
				sort.Slice(fields, func(i, j int) bool {
					return fields[j].id > fields[i].id
				})

				structs = append(structs, structType{name: s.Name, fields: fields})
			}
		}

		// Generate interface code
		var (
			t       types.Type
			nameMap map[string]struct{}
			params  = &templates.ModelTmplParams{}
		)

		// Convert struct to models and generate code
		for i := range structs {
			t = pkg.TypesInfo.TypeOf(structs[i].name)
			nameMap = make(map[string]struct{})

			var (
				index int
				model = templates.Model{}
			)
			model.GoName = types.ExprString(structs[i].name)
			model.TableName = rename(model.GoName)
			model.HasTableName = !IsImplemented(t, sqlTabler)
			model.HasColumn = !IsImplemented(t, sqlColumner)
			// model.HasRow = !IsImplemented(t, sqlRower)

			for _, f := range structs[i].fields {
				tv := f.tag.Get(cfg.Tag)

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
					return fmt.Errorf("sqlgen: struct %q has duplicate key %q in %s", structs[i].name, tf.ColumnName, dir)
				}
				nameMap[tf.ColumnName] = struct{}{}

				// Check type is sequel.Name, then override name
				model.Fields = append(model.Fields, tf)
			}

			clear(nameMap)
			params.Models = append(params.Models, &model)
		}

		tmplName := "model.go.tpl"
		tmpl, err := template.New(tmplName).Funcs(template.FuncMap{
			"quote":         strconv.Quote,
			"createTable":   createTableStmt(dialect),
			"alterTable":    alterTableStmt(dialect),
			"reserveImport": reserveImport(impPkg),
			"castAs":        castAs(impPkg),
			"addrOf":        addrOf(impPkg),
			"wrap":          dialect.Wrap,
			"var":           dialect.Var,
		}).ParseFS(codegenTemplates, "templates/"+tmplName)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(blr, params); err != nil {
			return err
		}

		if !cfg.SkipHeader {
			w.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")
		}

		w.WriteString("package " + pkg.Name + "\n\n")

		if len(impPkg.imports) > 0 {
			w.WriteString("import (\n")
			for _, pkg := range impPkg.imports {
				if filepath.Base(pkg.Path()) == pkg.Name() {
					w.WriteString("\t" + strconv.Quote(pkg.Path()) + "\n")
				} else {
					w.WriteString("\t" + pkg.Name() + " " + strconv.Quote(pkg.Path()) + "\n")
				}
			}
			w.WriteString(")\n")
		}

		w.WriteString(blr.String())

		fileDest := filepath.Join(dir, cfg.Exec.Filename)
		formatted, err := imports.Process(fileDest, w.Bytes(), &imports.Options{Comments: true})
		if err != nil {
			return err
		}
		blr.Reset()
		w.Reset()

		log.Println("Creating " + fileDest)
		if err := os.WriteFile(fileDest, formatted, fileMode); err != nil {
			return err
		}

		dirs = dirs[1:]
	}

	tmplName := "db.go.tpl"
	impPkg := new(Package)
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{
		"createTable":   createTableStmt(dialect),
		"alterTable":    alterTableStmt(dialect),
		"reserveImport": reserveImport(impPkg),
		"quote":         strconv.Quote,
		"wrap":          dialect.Wrap,
		"var":           dialect.Var,
	}).ParseFS(codegenTemplates, "templates/"+tmplName)
	if err != nil {
		return err
	}

	os.MkdirAll(cfg.Database.Dir, fileMode)
	filename = filepath.Join(cfg.Database.Dir, cfg.Database.Filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	if !cfg.SkipHeader {
		f.WriteString("// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.\n\n")
	}

	f.WriteString("package " + cfg.Database.Package + "\n\n")

	if err := tmpl.Execute(blr, &templates.DatabaseParams{
		PkgName: cfg.Database.Package,
	}); err != nil {
		return err
	}

	if len(impPkg.imports) > 0 {
		f.WriteString("import (\n")
		for _, pkg := range impPkg.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				f.WriteString("\t" + strconv.Quote(pkg.Path()) + "\n")
			} else {
				f.WriteString("\t" + pkg.Name() + " " + strconv.Quote(pkg.Path()) + "\n")
			}
		}
		f.WriteString(")\n")
	}
	f.WriteString(blr.String())

	return goModTidy()
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
