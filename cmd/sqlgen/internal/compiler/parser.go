package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

const pkgMode = packages.NeedName |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule |
	packages.NeedDeps

type Config struct {
	Tag        string
	RenameFunc func(string) string
	Matcher    Matcher
}

var annotationRegexp = regexp.MustCompile(`(?i)\s*\+sqlgen\:(.*)`)

func Parse(dir string, cfg *Config) (*Package, error) {
	// Load single go package and inspect the structs
	pkgs, err := packages.Load(&packages.Config{
		Dir:  dir,
		Mode: pkgMode,
	})
	if err != nil {
		return nil, err
	} else if len(pkgs) == 0 {
		return nil, ErrSkip
	}

	pkg := pkgs[0]
	structCaches := make([]structCache, 0)

	// Loop thru every go files
	for _, file := range pkg.Syntax {
		// If it's the generated code, we will skip it
		if ast.IsGenerated(file) {
			continue
		}

		// if pkg.Module != nil {
		// 	// If go version is 1.21, then it don't have infer type
		// 	if go121.Check(lo.Must1(semver.NewVersion(pkg.Module.GoVersion))) {
		// 		typeInferred = true
		// 	}
		// }

		// ast.Print(pkg.Fset, f)
		ast.Inspect(file, func(node ast.Node) bool {
			if node == nil {
				return true
			}

			genDecl := assertAsPtr[ast.GenDecl](node)
			if genDecl == nil || genDecl.Tok != token.TYPE {
				return true
			}

			for _, spec := range genDecl.Specs {
				typeSpec := assertAsPtr[ast.TypeSpec](spec)
				if typeSpec == nil {
					continue
				}

				// We only interested on Type Definition, or else we will skip
				// e.g: `type Entity sql.NullString`
				if typeSpec.Assign > 0 {
					continue
				}

				filename := pkg.Fset.Position(typeSpec.Name.NamePos).Filename
				if !cfg.Matcher.Match(filename) {
					continue
				}

				objType := pkg.TypesInfo.ObjectOf(typeSpec.Name)
				// We're not interested in the unexported type
				if !objType.Exported() {
					continue
				}
				// There are 2 types we're interested in
				// 1. struct (*ast.StructType)
				// 2. Type Definition from external package (*ast.SelectorExpr)
				//
				// The other situation like struct alias which we aren't cover it, e.g :
				// ```go
				// type A = time.Time
				// ```
				switch t := typeSpec.Type.(type) {
				case *ast.StructType:
					structCaches = append(structCaches, structCache{name: typeSpec.Name, doc: genDecl.Doc, t: t, pkg: pkg})

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
						continue
					}

					decl := assertAsPtr[ast.TypeSpec](obj.Decl)
					if decl == nil {
						continue
					}

					if v := assertAsPtr[ast.StructType](decl.Type); v != nil {
						structCaches = append(structCaches, structCache{name: typeSpec.Name, doc: genDecl.Doc, t: v, pkg: importPkg})
					}

				default:
					return true
				}
			}
			return true
		})
	}

	// If no files matched the pattern, skip it
	if len(structCaches) == 0 {
		return nil, ErrSkip
	}

	// Sort the struct inside the package by ascending order
	// Why do we need to sort the order because the output will look identical
	// everytimes it generated the codes
	sort.Slice(structCaches, func(i, j int) bool {
		return types.ExprString(structCaches[i].name) < types.ExprString(structCaches[j].name)
	})

	goPkg := new(Package)
	goPkg.Pkg = pkg
	goPkg.Tables = make([]*Table, 0, len(structCaches))

	// Loop every struct and inspect the fields
	for len(structCaches) > 0 {
		// Pop the first item for further inspection
		s := structCaches[0]
		// Struct queue, this is useful when handling embedded struct
		q := []typeQueue{{t: s.t, doc: s.doc, pkg: s.pkg}}
		structFields := make([]*structField, 0)

		var f typeQueue
		// Resolve every field on the struct, including nested struct
		for len(q) > 0 {
			f = q[0]

			// If the struct has empty field, just skip
			if len(f.t.Fields.List) == 0 {
				goto nextQueue
			}

			// Loop every struct field
			for i := range f.t.Fields.List {
				var tag reflect.StructTag
				fi := f.t.Fields.List[i]
				if fi.Tag != nil {
					// Trim backtick
					tag = reflect.StructTag(strings.TrimFunc(fi.Tag.Value, func(r rune) bool {
						return r == '`'
					}))
				}

				// If the field is embedded struct
				// `Type` can be either *ast.Ident or *ast.SelectorExpr
				if fi.Names == nil {
					t := fi.Type

					// If it's an embedded struct with pointer
					// we need to get the underlying type
					if ut := assertAsPtr[ast.StarExpr](fi.Type); ut != nil {
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

						t := obj.Decl.(*ast.TypeSpec)
						switch t.Type.(type) {
						case *ast.StructType:
							ft := &structField{
								name:  types.ExprString(vi),
								t:     f.pkg.TypesInfo.TypeOf(fi.Type),
								index: append(f.idx, i),
								// paths:    append(f.paths, path),
								exported: vi.IsExported(),
								embedded: true,
								parent:   f.prev,
								tag:      tag,
							}
							structFields = append(structFields, ft)

							q = append(q, typeQueue{
								idx:  append(f.idx, i),
								prev: ft,
								t:    t.Type.(*ast.StructType),
								pkg:  f.pkg,
							})
							continue
						case *ast.ArrayType:
							continue
						}

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

						// If it's a embedded struct, we continue on next loop
						if st := assertAsPtr[ast.StructType](decl.Type); st != nil {
							ft := &structField{
								name:     types.ExprString(vi.Sel),
								t:        f.pkg.TypesInfo.TypeOf(fi.Type),
								index:    append(f.idx, i),
								exported: vi.Sel.IsExported(),
								embedded: true,
								parent:   f.prev,
								tag:      tag,
							}

							q = append(q, typeQueue{
								idx:  append(f.idx, i),
								t:    st,
								prev: ft,
								pkg:  importPkg,
							})
							structFields = append(structFields, ft)
						}
						continue
					}
				}

				// Every struct field
				switch fv := fi.Type.(type) {
				// Imported types
				case *ast.SelectorExpr:
					// If the field type is a Go imported enum,
					// we will inspect it
					importPkg, ok := f.pkg.Imports[types.ExprString(fv.X)]
					if ok {
						log.Println(importPkg)
						// 	for _, file := range importPkg.Syntax {
						// 		ast.Inspect(file, func(n ast.Node) bool {
						// 			mapEnumIfExists(importPkg, n, enumMap)
						// 			return true
						// 		})
						// 	}
					}

				// Local types
				case *ast.Ident:
					if fv.Obj != nil {

					}
				}

				for j, n := range fi.Names {
					structFields = append(structFields, &structField{
						name:     types.ExprString(n),
						t:        f.pkg.TypesInfo.TypeOf(fi.Type),
						index:    append(f.idx, i+j),
						exported: n.IsExported(),
						parent:   f.prev,
						tag:      tag,
					})
				}
			}

		nextQueue:
			q = q[1:]
		}

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

		table := &Table{
			GoName: types.ExprString(s.name),
			Name:   cfg.RenameFunc(types.ExprString(s.name)),
			t:      pkg.TypesInfo.TypeOf(s.name),
		}

		if s.doc != nil {
			comment := f.doc.Text()
			if annotationRegexp.MatchString(comment) {
				annotations := annotationRegexp.FindStringSubmatch(comment)
				flags := strings.Split(annotations[1], ",")
				for len(flags) > 0 {
					switch flags[0] {
					case "readonly":
						table.Readonly = true
					case "ignore":
						// Skip this struct because of ignore
						structCaches = structCaches[1:]
						continue
					}
					flags = flags[1:]
				}
			}
		}

		// To store struct fields, to prevent field name collision
		nameMap := make(map[string]struct{})
		pos := 0
		for _, f := range structFields {
			if ptrCount(f.t) > 1 {
				return nil, fmt.Errorf(`sqlgen: pointer of pointer is not supported`)
			}

			var tagPaths []string
			tagVal := f.tag.Get(cfg.Tag)
			name := cfg.RenameFunc(f.name)
			if tagVal != "" {
				tagPaths = strings.Split(tagVal, ",")
				if len(tagPaths) > 0 {
					n := strings.TrimSpace(tagPaths[0])
					tagPaths = tagPaths[1:]
					if n == "-" {
						continue
					} else if n != "" {
						if !nameRegex.MatchString(n) {
							return nil, fmt.Errorf(`sqlgen: invalid column name %q in struct %q`, n, s.name)
						}
						name = n
					}
				}
			}

			switch f.t.String() {
			// If the type is table name, then we replace table name
			// and continue on next property
			//
			// Check type is sequel.Name, then override name
			case tableTypeName:
				if name != "" {
					table.Name = name
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

			if _, ok := nameMap[name]; ok {
				return nil, fmt.Errorf("sqlgen: struct %q has duplicate column name %q in directory %q", s.name, name, dir)
			}

			tag := parseTag(tagPaths)
			column := &Column{
				GoName:   f.name,
				GoPath:   f.FullGoPath(),
				Name:     name,
				Pos:      pos,
				Type:     f.t,
				Readonly: tag.hasOpts(TagOptionReadonly),
				goField:  f,
			}

			if tag.hasOpts(TagOptionAutoIncrement, TagOptionPK, TagOptionPKAlias) {
				if tag.hasOpts(TagOptionAutoIncrement) {
					if table.autoIncrKey != nil {
						return nil, fmt.Errorf(`sqlgen: you cannot have multiple auto increment key`)
					}
					table.autoIncrKey = column
				} else {
					table.Keys = append(table.Keys, column)
				}
				if column.Readonly {
					return nil, fmt.Errorf(`sqlgen: primary key cannot be readonly`)
				}
			}

			table.Columns = append(table.Columns, column)
			pos++
			nameMap[name] = struct{}{}
		}
		clear(nameMap)

		if len(table.Columns) > 0 {
			goPkg.Tables = append(goPkg.Tables, table)
		}
		structCaches = structCaches[1:]
	}
	return goPkg, nil
}

// This is to count the underlying pointer of a type
func ptrCount(t types.Type) int {
	var total int
loop:
	for t != nil {
		if v, ok := t.(*types.Pointer); ok {
			total++
			t = v.Elem()
		} else {
			break loop
		}
	}
	return total
}

func assertAsPtr[T any](v any) *T {
	t, ok := v.(*T)
	if ok {
		return t
	}
	return nil
}
