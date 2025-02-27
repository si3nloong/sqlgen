package codegen

import (
	"embed"
	"fmt"
	"go/ast"
	"go/types"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/Masterminds/semver/v3"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/sequel"
	"golang.org/x/tools/go/packages"
)

var (
	//go:embed templates/*.go.tpl
	codegenTemplates embed.FS

	// https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source
	path2Regex = strings.NewReplacer(
		`.`, `\.`,
		`*`, `.*`,
		`\`, `[\\/]`,
		`/`, `[\\/]`,
	)
	nameRegex     = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9\_]*$`)
	go121         = lo.Must1(semver.NewConstraint(">= 1.2.1"))
	goTagRegexp   = regexp.MustCompile(`(?i)^([a-z][a-z_]*[a-z])(\:(\w+))?$`)
	sqlFuncRegexp = regexp.MustCompile(`(?i)\s*(\w+\()(\w+\s*\,\s*)?(\{\})(\s*\,\s*\w+)?(\))\s*`)
	typeOfTable   = reflect.TypeOf(sequel.TableName{})
	tableNameType = typeOfTable.PkgPath() + "." + typeOfTable.Name()
)

const (
	TagOptionAutoIncrement = "auto_increment"
	TagOptionBinary        = "binary"
	TagOptionPKAlias       = "pk"
	TagOptionPK            = "primary_key"
	TagOptionFKAlias       = "fk"
	TagOptionFK            = "foreign_key"
	TagOptionUnsigned      = "unsigned"
	TagOptionSize          = "size"
	TagOptionDataType      = "data_type"
	TagOptionEncode        = "encode"
	TagOptionDecode        = "decode"
	TagOptionUnique        = "unique"
	TagOptionIndex         = "index"
)

type typeQueue struct {
	paths []string
	idx   []int
	prev  *structFieldType
	t     *ast.StructType
	pkg   *packages.Package
}

func Generate(c *Config) error {
	vldr := validator.New()
	if err := vldr.Struct(c); err != nil {
		return err
	}

	cfg := DefaultConfig()
	if c != nil {
		cfg = cfg.Merge(c)
	}

	dialect, ok := dialect.GetDialect((string)(cfg.Driver))
	if !ok {
		return fmt.Errorf("sqlgen: missing dialect, please register your dialect first")
	}

	generator, err := newGenerator(cfg, dialect)
	if err != nil {
		return err
	}

	var (
		srcDir  string
		sources = make([]string, len(cfg.Source))
	)
	copy(sources, cfg.Source)

	// Resolve every source provided
	for len(sources) > 0 {
		srcDir = strings.TrimSpace(sources[0])
		if srcDir == "" {
			return fmt.Errorf("sqlgen: source directory %q is empty path", srcDir)
		}

		if srcDir == "." {
			srcDir = fileutil.Getpwd()
			// If the prefix is ".", mean it's refer to current directory
		} else if srcDir[0] == '.' {
			srcDir = fileutil.Getpwd() + srcDir[1:]
		} else if srcDir[0] != '/' {
			srcDir = filepath.Join(fileutil.Getpwd(), srcDir)
		}

		// If suffix is *, we will add go extension to it
		if srcDir[len(srcDir)-1] == '*' {
			srcDir = srcDir + ".go"
		}

		slog.Info("Processing", "dir", srcDir)

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
			slog.Info("Submatch", "rootDir", rootDir, "dir", path2Regex.Replace(srcDir))
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

		if err := parseGoPackage(generator, rootDir, dirs, matcher); err != nil {
			return err
		}

	nextSrc:
		sources = sources[1:]
	}

	if cfg.Database != nil {
		// Generate db code
		_ = syscall.Unlink(filepath.Join(cfg.Database.Dir, cfg.Database.Filename))
		if err := renderTemplate(
			generator,
			"db.go.tpl",
			"",
			cfg.Database.Package,
			cfg.Database.Dir,
			cfg.Database.Filename,
		); err != nil {
			return err
		}
	}

	if cfg.Database.Operator != nil {
		_ = syscall.Unlink(filepath.Join(cfg.Database.Dir, cfg.Database.Operator.Filename))
		if err := renderTemplate(
			generator,
			"operator.go.tpl",
			"",
			cfg.Database.Operator.Package,
			cfg.Database.Operator.Dir,
			cfg.Database.Operator.Filename,
		); err != nil {
			return err
		}
	}

	if cfg.SkipModTidy {
		return nil
	}
	return goModTidy()
}

func parseGoPackage(
	gen *Generator,
	rootDir string,
	dirs []string,
	matcher Matcher,
) error {
	type structCache struct {
		name *ast.Ident
		t    *ast.StructType
		pkg  *packages.Package
	}

	var (
		dir      string
		filename string
		rename   = gen.config.RenameFunc()
	)

	for len(dirs) > 0 {
		dir = path.Join(rootDir, dirs[0])

		// Sometimes user might place db destination in the source as well
		// In this situation, we're not process the folder, we will skip it
		// if the file is exists in db folder
		if idx := lo.IndexOf([]string{
			path.Join(fileutil.Getpwd(), gen.config.Database.Dir),
			path.Join(fileutil.Getpwd(), gen.config.Database.Operator.Dir),
		}, dir); idx >= 0 {
			dirs = dirs[1:]
			continue
		}

		slog.Info("Process", "dir", dir)
		if fileutil.IsDirEmptyFiles(dir, gen.config.Exec.Filename) {
			slog.Info("Folder is empty, so not processing")
			dirs = dirs[1:]
			continue
		}

		filename = path.Join(dir, gen.config.Exec.Filename)
		// Unlink the generated file, ignore the error
		_ = syscall.Unlink(filename)

		// slog.Info("Load package", "dir", dir)
		// Since we're loading one directory at a time,
		// the return results will only return one package back
		pkgs, err := packages.Load(&packages.Config{
			Dir:  dir,
			Mode: pkgMode,
		})
		if err != nil {
			return err
		} else if len(pkgs) == 0 {
			return nil
		}

		var (
			pkg = pkgs[0]
			// enum cache in the packages
			enumCache    = make(map[string]*enum)
			typeInferred = false
			structCaches = make([]structCache, 0)
		)

		for _, file := range pkg.Syntax {
			// If it's the generated code, we will skip it
			if ast.IsGenerated(file) {
				continue
			}

			if pkg.Module != nil {
				// If go version is 1.21, then it don't have infer type
				if go121.Check(lo.Must1(semver.NewVersion(pkg.Module.GoVersion))) {
					typeInferred = true
				}
			}

			// ast.Print(pkg.Fset, f)
			ast.Inspect(file, func(node ast.Node) bool {
				if node == nil {
					return true
				}

				typeSpec := assertAsPtr[ast.TypeSpec](node)
				if typeSpec == nil {
					return true
				}

				// We only interested on Type Definition, or else we will skip
				// e.g: `type Entity sql.NullString`
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
			structs = make([]*structType, 0)
		)

		// Loop every struct and inspect the fields
		for len(structCaches) > 0 {
			var (
				s = structCaches[0]
				// Struct queue, this is useful when handling embedded struct
				q            = []typeQueue{{t: s.t, pkg: s.pkg}}
				f            typeQueue
				structFields = make([]*structFieldType, 0)
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
						var (
							t    = fi.Type
							path string
						)

						// If it's an embedded struct with pointer
						// we need to get the underlying type
						if ut := assertAsPtr[ast.StarExpr](fi.Type); ut != nil {
							path = "*"
							t = ut.X
						}

						switch vi := t.(type) {
						// Local struct
						case *ast.Ident:
							// Object is nil when it's not found in current scope (different file)
							obj := vi.Obj
							if vi.Obj == nil {
								// Since it's a local struct, we will find it in the local module files
								for i := range f.pkg.Syntax {
									obj = f.pkg.Syntax[i].Scope.Lookup(vi.Name)
									// exit when found the struct
									if obj != nil {
										break
									}
								}
							}
							// After lookup still cannot find the type, then skip
							if obj == nil {
								continue
							}

							path += types.ExprString(vi)
							// if f.path != "" {
							// 	path = f.path + "." + path
							// }
							t := obj.Decl.(*ast.TypeSpec)

							ft := &structFieldType{
								name:     types.ExprString(vi),
								index:    append(f.idx, i),
								paths:    append(f.paths, path),
								t:        f.pkg.TypesInfo.TypeOf(fi.Type),
								exported: vi.IsExported(),
								embedded: true,
								tag:      tag,
							}
							structFields = append(structFields, ft)

							q = append(q, typeQueue{
								paths: append(f.paths, path),
								idx:   append(f.idx, i),
								prev:  ft,
								t:     t.Type.(*ast.StructType),
								pkg:   f.pkg,
							})
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

							path += types.ExprString(vi.Sel)

							// If it's a embedded struct, we continue on next loop
							if st := assertAsPtr[ast.StructType](decl.Type); st != nil {
								paths := append(f.paths, path)
								ft := &structFieldType{
									name:     types.ExprString(vi.Sel),
									index:    append(f.idx, i),
									paths:    paths,
									t:        f.pkg.TypesInfo.TypeOf(fi.Type),
									exported: vi.Sel.IsExported(),
									embedded: true,
									parent:   f.prev,
									tag:      tag,
								}

								q = append(q, typeQueue{
									paths: paths,
									idx:   append(f.idx, i),
									t:     st,
									prev:  ft,
									pkg:   importPkg,
								})
								structFields = append(structFields, ft)
							}
							continue
						}
					}

					var goEnum *enum
					// Every struct field
					switch fv := fi.Type.(type) {
					// Imported types
					case *ast.SelectorExpr:
						// If the field type is a Go imported enum,
						// we will inspect it
						// importPkg, ok := f.pkg.Imports[types.ExprString(fv.X)]
						// if ok {
						// 	for _, file := range importPkg.Syntax {
						// 		ast.Inspect(file, func(n ast.Node) bool {
						// 			mapEnumIfExists(importPkg, n, enumMap)
						// 			return true
						// 		})
						// 	}
						// }

					// Local types
					case *ast.Ident:
						if fv.Obj != nil {
							mapGoEnums(enumCache, pkg, fv)
							// goEnum = mapGoEnums(enumCache, pkg, fv)
						}
					}

					for j, n := range fi.Names {
						path := types.ExprString(n)

						structFields = append(structFields, &structFieldType{
							name:     types.ExprString(n),
							index:    append(f.idx, i+j),
							paths:    append(f.paths, path),
							t:        f.pkg.TypesInfo.TypeOf(fi.Type),
							enums:    goEnum,
							exported: n.IsExported(),
							parent:   f.prev,
							tag:      tag,
						})
					}
				}

			nextQueue:
				q = q[1:]
			}

			if len(structFields) > 0 {
				sort.Slice(structFields, func(i, j int) bool {
					for k, f := range structFields[i].index {
						if k >= len(structFields[j].index) {
							return false
						}
						if f != structFields[j].index[k] {
							return f < structFields[j].index[k]
						}
					}
					return len(structFields[i].index) < len(structFields[j].index)
				})

				structs = append(structs, &structType{
					name:   types.ExprString(s.name),
					t:      pkg.TypesInfo.TypeOf(s.name),
					fields: structFields,
				})
			}

			structCaches = structCaches[1:]
		}

		// Sort the struct in package in ascending order
		sort.Slice(structs, func(i, j int) bool {
			return structs[i].name < structs[j].name
		})

		schemas := make([]*tableInfo, 0)
		for len(structs) > 0 {
			s := structs[0]

			// To store struct fields, to prevent field name collision
			nameDict := make(map[string]struct{})

			table := new(tableInfo)
			table.goName = s.name
			table.tableName = rename(table.goName)
			table.t = s.t

			var pos int
			for _, f := range s.fields {
				column := new(columnInfo)
				column.goName = f.name
				column.goPaths = f.paths
				column.t = f.t
				column.columnName = rename(f.name)
				column.columnPos = pos
				column.enums = f.enums

				name := s.name
				tagVal := strings.TrimSpace(f.tag.Get(gen.config.Tag))
				if tagVal != "" {
					tagPaths := strings.Split(tagVal, ",")
					name = strings.TrimSpace(tagPaths[0])
					// Skip field if user mentioned skip
					if name == "-" {
						continue
					} else if name != "" { // Column name must follow convention
						if !nameRegex.MatchString(name) {
							return fmt.Errorf(`sqlgen: invalid column name %q in struct %q`, name, s.name)
						}
						column.columnName = name
					}

					for _, v := range tagPaths[1:] {
						submatches := goTagRegexp.FindStringSubmatch(v)
						if len(submatches) > 2 {
							column.tags = append(column.tags, goTag{key: submatches[1], value: submatches[3]})
						}
					}
				}

				switch f.t.String() {
				// If the type is table name, then we replace table name
				// and continue on next property
				//
				// Check type is sequel.Name, then override name
				case tableNameType:
					if name != "" {
						table.tableName = name
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

				// Check uniqueness of the column
				if _, ok := nameDict[column.columnName]; ok {
					return fmt.Errorf("sqlgen: struct %q has duplicate column name %q in directory %q", s.name, column.columnName, dir)
				}

				if v, ok := column.getOptionValue(TagOptionSize); ok {
					column.size, err = strconv.Atoi(v)
					if err != nil {
						return fmt.Errorf(`sqlgen: invalid size value %q %w`, v, err)
					}
				}

				if columnType, ok := gen.columnTypes[f.t.String()]; ok {
					column.mapper = columnType
				} else if columnType, ok := gen.columnDataType(f.t); ok {
					column.mapper = columnType
				} else if columnType, ok := gen.defaultColumnTypes["*"]; ok {
					column.mapper = columnType
				} else {
					return fmt.Errorf(`sqlgen: missing data type mapping for data type %T`, f.t)
				}

				if column.hasOption(TagOptionAutoIncrement) {
					if table.autoIncrKey != nil {
						return fmt.Errorf(`sqlgen: you cannot have a composite key if you define auto increment key`)
					}
					table.autoIncrKey = column
					table.keys = append(table.keys, column)
				} else if column.hasOption(TagOptionPK) || column.hasOption(TagOptionPKAlias) {
					table.keys = append(table.keys, column)
				}

				table.columns = append(table.columns, column)
				nameDict[column.columnName] = struct{}{}
				pos++
			}

			// If the model doesn't consist any field,
			// we don't really want to generate the boilerplate code
			if len(table.columns) > 0 {
				schemas = append(schemas, table)
			}

			clear(nameDict)
			structs = structs[1:]
		}

		// If the `skip_empty` is true,
		// we do not generate the go file
		if gen.config.Exec.SkipEmpty && len(schemas) == 0 {
			goto nextDir
		}

		if err := gen.genModels(pkg, dir, typeInferred, schemas); err != nil {
			return err
		}

		if gen.config.Migration != nil {
			if err := os.MkdirAll(gen.config.Migration.Dir, os.ModePerm); err != nil {
				return err
			}

			if err := gen.genMigrations(schemas); err != nil {
				return err
			}
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

func isGoEnum(pkg *packages.Package, ident *ast.Ident) *ast.Object {
	for _, f := range pkg.Syntax {
		if obj := f.Scope.Lookup(types.ExprString(ident)); obj != nil {
			return obj
		}
	}
	return nil
}

func mapGoEnums(enumCache map[string]*enum, pkg *packages.Package, f *ast.Ident) *enum {
	// enumMap := make(map[string][]*parser.EnumValue)
	key := types.ExprString(f)
	if enum, ok := enumCache[key]; ok {
		return enum
	}

	// Loop thru every files
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			decl := assertAsPtr[ast.GenDecl](d)
			if decl == nil {
				continue
			}

			var typeName *ast.Ident
			for _, s := range decl.Specs {
				spec := assertAsPtr[ast.ValueSpec](s)
				if spec == nil {
					continue
				}

				n := assertAsPtr[ast.Ident](spec.Type)
				if n == nil {
					n = typeName
				}

				// If it's still empty type, we skip
				if n == nil {
					continue
				}

				if _, ok := enumCache[types.ExprString(n)]; !ok {
					enumCache[types.ExprString(n)] = new(enum)
				}

				values := make([]*enumValue, 0)
				// 	v := s.(*ast.ValueSpec) // safe because decl.Tok == token.VAR
				for _, name := range spec.Names {
					obj := pkg.TypesInfo.ObjectOf(name)
					switch v := obj.(type) {
					case *types.Const:
						values = append(values, &enumValue{name: types.ExprString(name), value: v.Val().ExactString()})
					case *types.Var:
						values = append(values, &enumValue{name: types.ExprString(name), value: v.String()})
					}
				}

				enumCache[types.ExprString(n)].values = append(enumCache[types.ExprString(n)].values, values...)
				typeName = n
			}
		}
	}

	return enumCache[key]
}

// func mapEnumIfExists(pkg *packages.Package, node ast.Node, enumMap map[string][]*goEnum) {
// 	switch v := node.(type) {
// 	case *ast.GenDecl:
// 		switch v.Tok {
// 		case token.CONST:
// 			var prevGoType string
// 			for _, spec := range v.Specs {
// 				valueSpec := assertAsPtr[ast.ValueSpec](spec)
// 				if valueSpec == nil {
// 					continue
// 				}

// 				typeName := assertAsPtr[ast.Ident](valueSpec.Type)
// 				if typeName == nil {
// 					if prevGoType == "" {
// 						return
// 					}

// 					mapGoEnums := enumMap[prevGoType]
// 					switch v := mapGoEnums[len(mapGoEnums)-1].value.(type) {
// 					case goIotaEnum:
// 						for _, n := range valueSpec.Names {
// 							enumMap[prevGoType] = append(enumMap[prevGoType], &goEnum{name: n, value: v + 1})
// 						}
// 					case goStringEnum:
// 						for _, n := range valueSpec.Names {
// 							enumMap[prevGoType] = append(enumMap[prevGoType], &goEnum{name: n, value: v})
// 						}
// 					}
// 				} else if typeName.IsExported() {
// 					obj := pkg.TypesInfo.ObjectOf(typeName)
// 					if !obj.Exported() {
// 						return
// 					}

// 					goType := obj.Type().String()
// 					val := types.ExprString(valueSpec.Values[0])
// 					if strings.Contains(val, "iota") {
// 						for _, n := range valueSpec.Names {
// 							enum := &goEnum{name: n, value: goStringEnum(val)}
// 							enumMap[goType] = append(enumMap[goType], enum)
// 						}
// 					} else {
// 						for _, n := range valueSpec.Names {
// 							enum := &goEnum{name: n, value: goIotaEnum(0)}
// 							enumMap[goType] = append(enumMap[goType], enum)
// 						}
// 					}
// 					prevGoType = goType
// 				}
// 			}
// 		}
// 	}
// }
