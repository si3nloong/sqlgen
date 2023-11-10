package codegen

import (
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

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/go/packages"
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
	TagOptionUnsigned      tagOption = "unsigned"
	TagOptionSize          tagOption = "size"
	TagOptionDataType      tagOption = "data_type"
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
	pkg    *packages.Package
	t      types.Type
	name   *ast.Ident
	fields []structField
}

type structField struct {
	id       string
	name     string
	path     string
	t        types.Type
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

	copy(sources, cfg.Source)

	// Resolve every source provided
	for len(sources) > 0 {
		srcDir = strings.TrimSpace(sources[0])
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
			rootDir = strings.TrimSuffix(subMatches[1], "/")
			dirs = append(dirs, "")
			matcher = &RegexMatcher{regexp.MustCompile(path2regex.Replace(srcDir))}
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
	type structCache struct {
		name *ast.Ident
		t    *ast.StructType
		pkg  *packages.Package
	}

	var (
		rename   = cfg.RenameFunc()
		dialect  = cfg.Dialect()
		dir      string
		filename string
	)

	// If `skip_escape` is false, we escape the table and column value
	if !cfg.SkipEscape {
		dialect = cfg.Dialect()
	}

	for len(dirs) > 0 {
		dir = path.Join(rootDir, dirs[0])
		if fileutil.IsDirEmptyFiles(dir, cfg.Exec.Filename) {
			dirs = dirs[1:]
			continue
		}

		filename = path.Join(dir, cfg.Exec.Filename)
		// Remove the generated file, ignore the error
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
			structTypes = make([]structCache, 0)
		)

		for _, f := range pkg.Syntax {
			// ast.Print(pkg.Fset, f)
			ast.Inspect(f, func(node ast.Node) bool {
				typeSpec := assertAsPtr[ast.TypeSpec](node)
				if typeSpec == nil {
					return true
				}

				// We only interested on Type Definition, or else we will skip
				// e.g: `type Model sql.NullString`
				if typeSpec.Assign > 0 {
					return true
				}

				filename = pkg.Fset.Position(typeSpec.Name.NamePos).Filename
				if !matcher.Match(filename) {
					return true
				}

				objType := pkg.TypesInfo.ObjectOf(typeSpec.Name)
				// We're not interested in the unexported type
				if !objType.Exported() {
					return true
				}

				// There are 2 types we're interested in
				// 1. struct (*ast.StructType)
				// 2. Type Definition from external package (*ast.SelectorExpr)
				switch t := typeSpec.Type.(type) {
				case *ast.StructType:
					structTypes = append(structTypes, structCache{name: typeSpec.Name, t: t, pkg: pkg})

				case *ast.SelectorExpr:
					var (
						pkgPath   = pkg.TypesInfo.ObjectOf(t.Sel).Pkg()
						importPkg = pkg.Imports[pkgPath.Path()]
						obj       *ast.Object
					)

					for i := range importPkg.Syntax {
						obj = importPkg.Syntax[i].Scope.Lookup(t.Sel.Name)
						if obj != nil {
							break
						}
					}

					decl := assertAsPtr[ast.TypeSpec](obj.Decl)
					if decl == nil {
						return true
					}

					if v := assertAsPtr[ast.StructType](decl.Type); v != nil {
						structTypes = append(structTypes, structCache{name: typeSpec.Name, t: v, pkg: importPkg})
					}
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
				// Struct queue, this is useful when handling embedded struct
				q      = []typeQueue{{t: s.t}}
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
							// TODO: support imported struct
							log.Println("debug =>", vi)
						}
					}

					for j, n := range fi.Names {
						path := types.ExprString(n)
						if f.path != "" {
							path = f.path + "." + path
						}

						fields = append(fields, structField{id: toID(append(f.idx, i+j)), exported: n.IsExported(), name: types.ExprString(n), tag: tag, path: path, t: s.pkg.TypesInfo.TypeOf(fi.Type)})
					}
				}

			next:
				q = q[1:]
			}
			if len(fields) > 0 {
				sort.Slice(fields, func(i, j int) bool {
					return fields[j].id > fields[i].id
				})

				structs = append(structs, structType{name: s.name, fields: fields, pkg: s.pkg, t: pkg.TypesInfo.TypeOf(s.name)})
			}
		}

		// Generate interface code
		var (
			nameMap map[string]struct{}
			params  = templates.ModelTmplParams{}
		)

		// Convert struct to models and generate code
		for _, s := range structs {
			nameMap = make(map[string]struct{})

			var (
				index int
				model = templates.Model{}
			)
			model.GoName = types.ExprString(s.name)
			model.TableName = rename(model.GoName)

			// Check struct implements `sequel.Tabler`
			if m, w := types.MissingMethod(s.t, sqlTabler, true); cfg.Strict && w {
				return fmt.Errorf(`sqlgen: struct %q has implements "sequel.Tabler" but wrong footprint`, s.name)
			} else if m == nil && !w {
				model.HasTableName = true
			}

			// Check struct implements `sequel.Columner`
			if m, w := types.MissingMethod(s.t, sqlColumner, true); cfg.Strict && w {
				return fmt.Errorf(`sqlgen: struct %q has implements "sequel.Columner" but wrong footprint`, s.name)
			} else if m == nil && !w {
				model.HasColumn = true
			}

			for _, f := range s.fields {
				var (
					tv  = f.tag.Get(cfg.Tag)
					tf  = &templates.Field{}
					tag = make(tagOpts)
					n   string
				)

				tf.ColumnName = rename(f.name)

				if tv != "" {
					tags := strings.Split(tv, ",")
					n = strings.TrimSpace(tags[0])
					if n == "-" {
						continue
					} else if n != "" {
						tf.ColumnName = n
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

				switch f.t.String() {
				// If the type is table name, then we replace table name
				// and continue on next property
				case tableNameSchema:
					if n != "" {
						model.TableName = n
					}
					continue
				}

				// If it's a unexported field, skip!
				if !f.exported {
					continue
				}

				tf.Type = f.t
				tf.GoName = f.name
				tf.GoPath = f.path
				tf.Index = index
				if _, ok := tag.Lookup(TagOptionBinary); ok {
					if IsImplemented(tf.Type, binaryMarshaler) && IsImplemented(newPointer(tf.Type), binaryUnmarshaler) {
						tf.IsBinary = true
					} else if cfg.Strict {
						return fmt.Errorf(`sqlgen: field %q of struct %q specific for "binary" must comply to encoding.BinaryMarshaler and encoding.BinaryUnmarshaler`, tf.GoName, model.GoName)
					}
				}
				if v, ok := tag.Lookup(TagOptionSize); ok {
					tf.Size, _ = strconv.Atoi(v)
				}
				tf.IsTextMarshaler = IsImplemented(tf.Type, textMarshaler)
				tf.IsTextUnmarshaler = IsImplemented(newPointer(tf.Type), textUnmarshaler)
				index++

				if _, ok := tag.Lookup(TagOptionPK, TagOptionPKAlias, TagOptionAutoIncrement); ok {
					if cfg.Strict && model.PK != nil {
						return fmt.Errorf(`sqlgen: a model can only allow one primary key, else it will get overriden`)
					}

					// Check auto increment
					pk := templates.PK{Field: tf}
					_, pk.IsAutoIncr = tag.Lookup(TagOptionAutoIncrement)
					model.PK = &pk
				}

				if cfg.Strict {
					// Check uniqueness of the column
					if _, ok := nameMap[tf.ColumnName]; ok {
						return fmt.Errorf("sqlgen: struct %q has duplicate key %q in %s", s.name, tf.ColumnName, dir)
					}
				}
				nameMap[tf.ColumnName] = struct{}{}

				// Check type is sequel.Name, then override name
				model.Fields = append(model.Fields, tf)
			}

			clear(nameMap)

			// If the model doesn't consist any field,
			// we don't really want to generate the boilerplate code
			if len(model.Fields) > 0 {
				params.Models = append(params.Models, &model)
			}
		}

		if cfg.Exec.SkipEmpty && len(params.Models) == 0 {
			goto nextDir
		}

		if err := renderTemplate(
			"model.go.tpl",
			cfg.SkipHeader,
			dialect,
			pkg.Name,
			dir,
			cfg.Exec.Filename,
			params,
		); err != nil {
			return err
		}

	nextDir:
		dirs = dirs[1:]
	}

	if cfg.Database != nil {
		// Generate db code
		if err := renderTemplate(
			"db.go.tpl",
			cfg.SkipHeader,
			dialect,
			cfg.Database.Package,
			cfg.Database.Dir,
			cfg.Database.Filename,
			templates.DBTmplParams{},
		); err != nil {
			return err
		}
	}

	if cfg.SkipModTidy {
		return nil
	}
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
	buf := strpool.AcquireString()
	defer strpool.ReleaseString(buf)
	for i, v := range val {
		if i > 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(strconv.Itoa(v))
	}
	return buf.String()
}

func assertAsPtr[T any](v any) *T {
	t, ok := v.(*T)
	if ok {
		return t
	}
	return nil
}

func UnderlyingType(t types.Type) (*Mapping, bool) {
	var (
		typeStr string
		prev    = t
	)

loop:
	for t != nil {
		switch v := t.(type) {
		case *types.Basic:
			typeStr += v.String()
			break loop
		case *types.Named:
			if _, ok := v.Underlying().(*types.Struct); ok {
				typeStr += v.String()
				break loop
			}
			typeStr += v.Underlying().String()
			prev = t.Underlying()
		case *types.Pointer:
			typeStr += "*"
			prev = v.Elem()
		case *types.Slice:
			typeStr += "[]"
			prev = v.Elem()
		default:
			break loop
		}
		if v, ok := typeMap[typeStr]; ok {
			return v, ok
		}
		if prev == t {
			break loop
		}
		t = prev
	}
	if v, ok := typeMap[typeStr]; ok {
		return v, ok
	}
	return nil, false
}
