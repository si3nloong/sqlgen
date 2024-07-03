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
	"strings"
	"syscall"

	"github.com/Masterminds/semver/v3"
	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/config"
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
	codegenRegex  = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)
	go121         = lo.Must1(semver.NewConstraint(">= 1.2.1"))
	sqlFuncRegexp = regexp.MustCompile(`(?i)\s*(\w+\()(\w+\s*\,\s*)?(\{\})(\s*\,\s*\w+)?(\))\s*`)
)

const fileMode = 0o755

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
	t      types.Type
	name   *ast.Ident
	fields []structField
}

type structField struct {
	index    []int
	name     string
	path     string
	t        types.Type
	exported bool
	embedded bool
	tag      reflect.StructTag
}

type tagOpts map[string]string

func (t tagOpts) Lookup(key string, keys ...string) (string, bool) {
	keys = append(keys, key)
	for k, val := range t {
		if lo.IndexOf(keys, k) >= 0 {
			return val, true
		}
	}
	return "", false
}

type tableInfo struct {
	goName      string
	dbName      string
	tableName   string
	t           types.Type
	autoIncrKey *columnInfo
	keys        []*columnInfo
	columns     []*columnInfo
	indexes     []*indexInfo
}

func (b *tableInfo) GoName() string {
	return b.goName
}

func (b *tableInfo) DatabaseName() string {
	return b.dbName
}

func (b *tableInfo) TableName() string {
	return b.tableName
}

func (b *tableInfo) AutoIncrKey() (sequel.GoColumnSchema, bool) {
	return b.autoIncrKey, b.autoIncrKey != nil
}

func (b *tableInfo) Keys() []string {
	return lo.Map(b.keys, func(c *columnInfo, _ int) string {
		return c.colName
	})
}

func (b *tableInfo) Columns() []string {
	return lo.Map(b.columns, func(c *columnInfo, _ int) string {
		return c.colName
	})
}

func (b *tableInfo) Implements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(b.t, T, true)
}

func (b *tableInfo) PtrImplements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(types.NewPointer(b.t), T, true)
}

func (b *tableInfo) Column(i int) sequel.GoColumnSchema {
	return b.columns[i]
}

func (b *tableInfo) Index(i int) sequel.GoIndexSchema {
	return b.indexes[i]
}

// Mean table has only pk
func (b tableInfo) hasNoColsExceptPK() bool {
	return len(b.keys) == len(b.columns)
}

type columnInfo struct {
	goName  string
	goPath  string
	colName string
	colPos  int
	t       types.Type
	tag     tagOpts
	model   *config.Model
	size    int64
}

var (
	_ (sequel.GoColumnSchema) = (*columnInfo)(nil)
)

func (i columnInfo) SQLValuer() sequel.QueryFunc {
	if i.model == nil {
		return nil
	}
	return func(placeholder string) string {
		return strings.Replace(i.model.SQLValuer, "{placeholder}", placeholder, 1)
	}
}

func (i columnInfo) SQLScanner() sequel.QueryFunc {
	if i.model == nil {
		return nil
	}
	return func(column string) string {
		return strings.Replace(i.model.SQLScanner, "{column}", column, 1)
	}
}

func (c *columnInfo) GoName() string {
	return c.goName
}

func (c *columnInfo) GoPath() string {
	return c.goPath
}

func (c *columnInfo) Type() types.Type {
	return c.t
}

func (c *columnInfo) Size() int64 {
	return c.size
}

func (c *columnInfo) ColumnName() string {
	return c.colName
}

func (c *columnInfo) ColumnPos() int {
	return c.colPos
}

func (c *columnInfo) Implements(T *types.Interface) (wrongType bool) {
	_, wrongType = types.MissingMethod(c.t, T, true)
	return
}

type indexInfo struct {
	columns   []string
	indexType string
}

func (i indexInfo) Columns() []string {
	return i.columns
}

func (i indexInfo) Type() string {
	return i.indexType
}

func Generate(c *config.Config) error {
	cfg := config.DefaultConfig()
	if c != nil {
		cfg = c
	}

	cfg.Init()

	dialect, ok := sequel.GetDialect(string(cfg.Driver))
	if !ok {
		panic("sqlgen: missing dialect, please register your dialect first")
	}

	var (
		srcDir  string
		sources = make([]string, len(cfg.Source))
		gen     = newGenerator(cfg, dialect)
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
				matcher = &FileMatcher{filepath.Join(srcDir, fi.Name()): struct{}{}}
			}

			rootDir = srcDir
			dirs = append(dirs, "")
		}

		if err := parseGoPackage(gen, rootDir, dirs, matcher); err != nil {
			return err
		}

	nextSrc:
		sources = sources[1:]
	}

	if cfg.Database != nil {
		// Generate db code
		_ = syscall.Unlink(filepath.Join(cfg.Database.Dir, cfg.Database.Filename))
		if err := renderTemplate(
			dialect,
			"db.go.tpl",
			// true,
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
			dialect,
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

		// Skip if the file is exists in db folder
		if idx := sort.SearchStrings([]string{
			path.Join(fileutil.Getpwd(), gen.config.Database.Dir),
			path.Join(fileutil.Getpwd(), gen.config.Database.Operator.Dir),
		}, dir); idx != 2 {
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
			pkg          = pkgs[0]
			typeInferred = false
			structCaches = make([]structCache, 0)
		)

		for _, file := range pkg.Syntax {
			if len(file.Comments) > 0 && len(file.Comments[0].List) > 0 {
				// If the first comment is "^// Code generated .* DO NOT EDIT\.$"
				// we will skip it
				if codegenRegex.MatchString(file.Comments[0].List[0].Text) {
					continue
				}
			}

			if pkg.Module != nil {
				// If go version is 1.21, then it don't have infer type
				if go121.Check(lo.Must1(semver.NewVersion(pkg.Module.GoVersion))) {
					typeInferred = true
				}
			}

			// ast.Print(pkg.Fset, f)
			ast.Inspect(file, func(node ast.Node) bool {
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

							if obj == nil {
								continue
							}

							path := types.ExprString(vi)
							if f.path != "" {
								path = f.path + "." + path
							}
							t := obj.Decl.(*ast.TypeSpec)
							q = append(q, typeQueue{path: path, idx: append(f.idx, i), t: t.Type.(*ast.StructType), pkg: f.pkg})

							fields = append(fields, structField{index: append(f.idx, i), exported: vi.IsExported(), embedded: true, name: types.ExprString(vi), tag: tag, path: path, t: f.pkg.TypesInfo.TypeOf(fi.Type)})
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

								fields = append(fields, structField{index: append(f.idx, i), exported: vi.Sel.IsExported(), embedded: true, name: types.ExprString(vi.Sel), tag: tag, path: path, t: f.pkg.TypesInfo.TypeOf(fi.Type)})
							}
							continue
						}
					}

					for j, n := range fi.Names {
						path := types.ExprString(n)
						if f.path != "" {
							path = f.path + "." + path
						}

						fields = append(fields, structField{index: append(f.idx, i+j), exported: n.IsExported(), name: types.ExprString(n), tag: tag, path: path, t: f.pkg.TypesInfo.TypeOf(fi.Type)})
					}
				}

			nextQueue:
				q = q[1:]
			}

			if len(fields) > 0 {
				sort.Slice(fields, func(i, j int) bool {
					for k, xik := range fields[i].index {
						if k >= len(fields[j].index) {
							return false
						}
						if xik != fields[j].index[k] {
							return xik < fields[j].index[k]
						}
					}
					return len(fields[i].index) < len(fields[j].index)
				})
				structs = append(structs, structType{name: s.name, fields: fields, t: pkg.TypesInfo.TypeOf(s.name)})
			}

			structCaches = structCaches[1:]
		}

		schemaList := make([]*tableInfo, 0)

		for _, s := range structs {
			var (
				pos      int
				table    = new(tableInfo)
				nameDict = make(map[string]struct{})
			)

			table.t = s.t
			table.goName = types.ExprString(s.name)
			table.tableName = rename(table.goName)

			for _, f := range s.fields {
				var (
					column = new(columnInfo)
					name   string
				)

				column.goName = f.name
				column.goPath = f.path
				column.t = f.t
				column.colName = rename(f.name)
				column.colPos = pos
				column.tag = make(tagOpts)

				if model, ok := gen.config.Models[f.t.String()]; ok {
					column.model = model
				}

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
						column.colName = name
					}

					for _, v := range tagPaths[1:] {
						kv := strings.SplitN(v, ":", 2)
						k := strings.TrimSpace(strings.ToLower(kv[0]))
						if len(kv) > 1 {
							column.tag[k] = kv[1]
						} else {
							column.tag[k] = ""
						}
					}
				}

				switch f.t.String() {
				// If the type is table name, then we replace table name
				// and continue on next property
				case tableNameSchema:
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
				if _, ok := nameDict[column.colName]; ok {
					return fmt.Errorf("sqlgen: struct %q has duplicate key %q in %s", s.name, column.colName, dir)
				}

				if _, ok := column.tag.Lookup(TagOptionPK, TagOptionPKAlias, TagOptionAutoIncrement); ok {
					// Check auto increment
					_, autoIncr := column.tag.Lookup(TagOptionAutoIncrement)
					if autoIncr {
						if table.autoIncrKey != nil {
							return fmt.Errorf(`sqlgen: you cannot have a composite key if you define auto increment key`)
						}
						table.autoIncrKey = column
					}
					table.keys = append(table.keys, column)
				}

				nameDict[column.colName] = struct{}{}

				column.colName = gen.QuoteIdentifier(column.colName)
				// TODO: Check type is sequel.Name, then override name
				table.columns = append(table.columns, column)
				pos++
			}

			clear(nameDict)
			// FIXME: we should allow database name to declare
			// table.dbName = gen.QuoteIdentifier(table.dbName)
			table.tableName = gen.QuoteIdentifier(table.tableName)

			// If the model doesn't consist any field,
			// we don't really want to generate the boilerplate code
			if len(table.columns) > 0 {
				schemaList = append(schemaList, table)
			}
		}

		if gen.config.Exec.SkipEmpty && len(schemaList) == 0 {
			goto nextDir
		}

		if err := gen.generate(pkg, dir, typeInferred, schemaList); err != nil {
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
