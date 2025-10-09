package codegen

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/Masterminds/semver/v3"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/compiler"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/fileutil"
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
	nameRegex = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9\_]*$`)
	go121     = lo.Must1(semver.NewConstraint(">= 1.2.1"))
	// goTagRegexp   = regexp.MustCompile(`(?i)^([a-z][a-z_]*[a-z])(\:(\w+))?$`)
	// sqlFuncRegexp = regexp.MustCompile(`(?i)\s*(\w+\()(\w+\s*\,\s*)?(\{\})(\s*\,\s*\w+)?(\))\s*`)
)

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
		return fmt.Errorf("sqlgen: missing dialect, please register dialect %q", cfg.Driver)
	}

	generator, err := newGenerator(cfg, dialect)
	if err != nil {
		return err
	}

	var (
		srcDir  string
		sources = make([]string, 0, len(cfg.Source))
	)
	sources = append(sources, cfg.Source...)

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

			if fi.IsDir() {
				// If it's just a folder
				matcher = FolderMatcher(srcDir)
			} else {
				// If it's just a file
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
	g *Generator,
	rootDir string,
	dirs []string,
	matcher Matcher,
) error {
	var dir string
	var filename string
	rename := g.config.RenameFunc()

	for len(dirs) > 0 {
		dir = path.Join(rootDir, dirs[0])

		// Sometimes user might place db destination in the source as well
		// In this situation, we're not process the folder, we will skip it
		// if the file is exists in db folder
		pwd := fileutil.Getpwd()
		if idx := lo.IndexOf([]string{
			path.Join(pwd, g.config.Database.Dir),
			path.Join(pwd, g.config.Database.Operator.Dir),
		}, dir); idx >= 0 {
			dirs = dirs[1:]
			continue
		}

		slog.Info("Process", "dir", dir)
		if fileutil.IsDirEmptyFiles(dir, g.config.Exec.Filename) {
			slog.Info("Folder is empty, so not processing")
			dirs = dirs[1:]
			continue
		}

		filename = path.Join(dir, g.config.Exec.Filename)
		// Unlink the generated file, ignore the error
		_ = syscall.Unlink(filename)

		slog.Info("Parse go package", "dir", dir)
		// Since we're loading one directory at a time,
		// the return results will only return one package back
		schema, err := compiler.Parse(dir, &compiler.Config{
			Tag:        g.config.Tag,
			RenameFunc: rename,
			Matcher:    matcher,
		})
		if errors.Is(err, compiler.ErrSkip) {
			goto nextDir
		} else if err != nil {
			return err
		}

		if err := g.generateModels(dir, schema); err != nil {
			return err
		}

		// If the `skip_empty` is true,
		// we do not generate the go file
		if g.config.Exec.SkipEmpty {
			goto nextDir
		}

		// if g.config.Migration != nil {
		// 	if err := os.MkdirAll(g.config.Migration.Dir, os.ModePerm); err != nil {
		// 		return err
		// 	}

		// 	if err := g.genMigrations(schema); err != nil {
		// 		return err
		// 	}
		// }

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

// func isGoEnum(pkg *packages.Package, ident *ast.Ident) *ast.Object {
// 	for _, f := range pkg.Syntax {
// 		if obj := f.Scope.Lookup(types.ExprString(ident)); obj != nil {
// 			return obj
// 		}
// 	}
// 	return nil
// }

// func mapGoEnums(enumCache map[string]*enum, pkg *packages.Package, f *ast.Ident) *enum {
// 	// enumMap := make(map[string][]*parser.EnumValue)
// 	key := types.ExprString(f)
// 	if enum, ok := enumCache[key]; ok {
// 		return enum
// 	}

// 	// Loop thru every files
// 	for _, f := range pkg.Syntax {
// 		for _, d := range f.Decls {
// 			decl := assertAsPtr[ast.GenDecl](d)
// 			if decl == nil {
// 				continue
// 			}

// 			var typeName *ast.Ident
// 			for _, s := range decl.Specs {
// 				spec := assertAsPtr[ast.ValueSpec](s)
// 				if spec == nil {
// 					continue
// 				}

// 				n := assertAsPtr[ast.Ident](spec.Type)
// 				if n == nil {
// 					n = typeName
// 				}

// 				// If it's still empty type, we skip
// 				if n == nil {
// 					continue
// 				}

// 				if _, ok := enumCache[types.ExprString(n)]; !ok {
// 					enumCache[types.ExprString(n)] = new(enum)
// 				}

// 				values := make([]*enumValue, 0)
// 				// 	v := s.(*ast.ValueSpec) // safe because decl.Tok == token.VAR
// 				for _, name := range spec.Names {
// 					obj := pkg.TypesInfo.ObjectOf(name)
// 					switch v := obj.(type) {
// 					case *types.Const:
// 						values = append(values, &enumValue{name: types.ExprString(name), value: v.Val().ExactString()})
// 					case *types.Var:
// 						values = append(values, &enumValue{name: types.ExprString(name), value: v.String()})
// 					}
// 				}

// 				enumCache[types.ExprString(n)].values = append(enumCache[types.ExprString(n)].values, values...)
// 				typeName = n
// 			}
// 		}
// 	}

// 	return enumCache[key]
// }

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
