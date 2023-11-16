package codegen

import (
	"embed"
	"fmt"
	"go/ast"
	"go/types"
	"io/fs"
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
	TagOptionEncode        tagOption = "encode"
	TagOptionDecode        tagOption = "decode"
	TagOptionUnique        tagOption = "unique"
)

var (
	schemaName      = reflect.TypeOf(sequel.Table{})
	tableNameSchema = schemaName.PkgPath() + "." + schemaName.Name()
)

type typeQueue struct {
	path string
	idx  []int
	t    *ast.StructType
	pkg  *packages.Package
}

type structType struct {
	// pkg    *packages.Package
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
	embedded bool
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

var path2Regex = strings.NewReplacer(
	`.`, `\.`,
	`*`, `.*`,
	`\`, `[\\/]`,
	`/`, `[\\/]`,
)

var nameRegex = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9\_]*$`)

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

		if srcDir == "." {
			srcDir = fileutil.Getpwd()
			// If the prefix is ".", mean it's refer to current directory
		} else if srcDir[0] == '.' {
			srcDir = fileutil.Getpwd() + srcDir[1:]
		} else if srcDir[0] != '/' {
			srcDir = filepath.Join(fileutil.Getpwd(), srcDir)
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
				suffix = path2Regex.Replace(paths[1])
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
			matcher = &RegexMatcher{regexp.MustCompile(path2Regex.Replace(rootDir) + `([\\/][a-z0-9_-]+)*` + suffix)}
		} else if len(subMatches) > 0 {
			rootDir = strings.TrimSuffix(subMatches[1], "/")
			dirs = append(dirs, "")
			matcher = &RegexMatcher{regexp.MustCompile(path2Regex.Replace(srcDir))}
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

	if cfg.Database != nil {
		// Generate db code
		os.Remove(filepath.Join(cfg.Database.Dir, cfg.Database.Filename))
		if err := renderTemplate(
			"db.go.tpl",
			cfg.SkipHeader,
			cfg.Dialect(),
			"",
			cfg.Database.Package,
			cfg.Getter.Prefix,
			cfg.Database.Dir,
			cfg.Database.Filename,
			struct{}{},
		); err != nil {
			return err
		}
	}

	if cfg.Database.Operator != nil {
		os.Remove(filepath.Join(cfg.Database.Dir, cfg.Database.Operator.Filename))
		if err := renderTemplate(
			"operator.go.tpl",
			cfg.SkipHeader,
			cfg.Dialect(),
			"",
			cfg.Database.Operator.Package,
			cfg.Getter.Prefix,
			cfg.Database.Operator.Dir,
			cfg.Database.Operator.Filename,
			struct{}{},
		); err != nil {
			return err
		}
	}

	if cfg.SkipModTidy {
		return nil
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
			structCaches = make([]structCache, 0)
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
				//
				// The other is struct alias which we aren't cover, e.g :
				// ```go
				// type A = time.Time
				// ```
				switch t := typeSpec.Type.(type) {
				case *ast.StructType:
					structCaches = append(structCaches, structCache{name: typeSpec.Name, t: t, pkg: pkg})

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

					// Skip if unable to find the specific object
					if obj == nil {
						return true
					}

					decl := assertAsPtr[ast.TypeSpec](obj.Decl)
					if decl == nil {
						return true
					}

					if v := assertAsPtr[ast.StructType](decl.Type); v != nil {
						structCaches = append(structCaches, structCache{name: typeSpec.Name, t: v, pkg: importPkg})
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

		// Loop every struct and inspect the fields
		for len(structCaches) > 0 {
			var (
				s = structCaches[0]
				// Struct queue, this is useful when handling embedded struct
				q      = []typeQueue{{t: s.t, pkg: s.pkg}}
				f      typeQueue
				fields = make([]structField, 0)
			)

			for len(q) > 0 {
				f = q[0]

				// If the struct has empty field, just skip
				if len(f.t.Fields.List) == 0 {
					goto nextQueue
				}

				// Loop every struct field
				for i := range f.t.Fields.List {
					var (
						tag reflect.StructTag
						fi  = f.t.Fields.List[i]
					)
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
							q = append(q, typeQueue{path: path, idx: append(f.idx, i), t: t.Type.(*ast.StructType), pkg: f.pkg})

							fields = append(fields, structField{id: toID(append(f.idx, i)), exported: vi.IsExported(), embedded: true, name: types.ExprString(vi), tag: tag, path: path, t: f.pkg.TypesInfo.TypeOf(fi.Type)})
							continue

						// Embedded with imported struct
						case *ast.SelectorExpr:
							var (
								t         = f.pkg.TypesInfo.TypeOf(vi)
								pkgPath   = t.String()
								idx       = strings.LastIndex(pkgPath, ".")
								importPkg = f.pkg.Imports[pkgPath[:idx]]
								obj       *ast.Object
							)

							for i := range importPkg.Syntax {
								obj = importPkg.Syntax[i].Scope.Lookup(vi.Sel.Name)
								if obj != nil {
									break
								}
							}

							// Skip if unable to find the specific object
							if obj == nil {
								continue
							}

							decl := assertAsPtr[ast.TypeSpec](obj.Decl)
							if decl == nil {
								continue
							}

							path := types.ExprString(vi.Sel)
							if f.path != "" {
								path = f.path + "." + path
							}

							// If it's a embedded struct, we continue on next loop
							if st := assertAsPtr[ast.StructType](decl.Type); st != nil {
								q = append(q, typeQueue{path: path, idx: append(f.idx, i), t: st, pkg: importPkg})

								fields = append(fields, structField{id: toID(append(f.idx, i)), exported: vi.Sel.IsExported(), embedded: true, name: types.ExprString(vi.Sel), tag: tag, path: path, t: f.pkg.TypesInfo.TypeOf(fi.Type)})
							}
							continue
						}
					}

					for j, n := range fi.Names {
						path := types.ExprString(n)
						if f.path != "" {
							path = f.path + "." + path
						}

						fields = append(fields, structField{id: toID(append(f.idx, i+j)), exported: n.IsExported(), name: types.ExprString(n), tag: tag, path: path, t: f.pkg.TypesInfo.TypeOf(fi.Type)})
					}
				}

			nextQueue:
				q = q[1:]
			}

			if len(fields) > 0 {
				sort.Slice(fields, func(i, j int) bool {
					return fields[j].id > fields[i].id
				})
				structs = append(structs, structType{name: s.name, fields: fields, t: pkg.TypesInfo.TypeOf(s.name)})
			}

			structCaches = structCaches[1:]
		}

		// Generate interface code
		var (
			params = templates.ModelTmplParams{}
		)

		// Convert struct to models and generate code
		for _, s := range structs {
			var (
				index   int
				nameMap = make(map[string]struct{})
				model   = templates.Model{}
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
						if cfg.Strict && !nameRegex.MatchString(n) {
							return fmt.Errorf(`sqlgen: invalid column name %q in struct %q`, n, s.name)
						}
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
				if f.embedded {
					continue
				}

				tf.Type = f.t
				tf.GoName = f.name
				tf.GoPath = f.path
				tf.Index = index
				if val, ok := tag.Lookup(TagOptionEncode); ok {
					tf.CustomMarshaler = val
				}
				if val, ok := tag.Lookup(TagOptionDecode); ok {
					tf.CustomUnmarshaler = val
				}
				if _, ok := tag.Lookup(TagOptionBinary); ok {
					if isImplemented(tf.Type, binaryMarshaler) && isImplemented(newPointer(tf.Type), binaryUnmarshaler) {
						tf.IsBinary = true
					} else if cfg.Strict {
						return fmt.Errorf(`sqlgen: field %q of struct %q specific for "binary" must comply to encoding.BinaryMarshaler and encoding.BinaryUnmarshaler`, tf.GoName, model.GoName)
					}
				}
				if val, ok := tag.Lookup(TagOptionSize); ok {
					tf.Size, _ = strconv.Atoi(val)
				}
				tf.IsTextMarshaler = isImplemented(tf.Type, textMarshaler)
				tf.IsTextUnmarshaler = isImplemented(newPointer(tf.Type), textUnmarshaler)
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
			pkg.PkgPath,
			pkg.Name,
			cfg.Getter.Prefix,
			dir,
			cfg.Exec.Filename,
			params,
		); err != nil {
			return err
		}

	nextDir:
		dirs = dirs[1:]
	}

	return nil
}

func goModTidy() error {
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stdout
	return tidyCmd.Run()
}
