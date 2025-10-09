package compiler

import (
	"errors"
	"go/ast"
	"go/types"
	"reflect"
	"regexp"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/go/packages"
)

var (
	ErrSkip = errors.New(`compiler: skip code generation`)

	nameRegex     = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9\_]*$`)
	typeOfTable   = reflect.TypeOf(sequel.TableName{})
	tableTypeName = typeOfTable.PkgPath() + "." + typeOfTable.Name()
)

type Matcher interface {
	Match(v string) bool
}

type Package struct {
	Pkg    *packages.Package
	Tables []*Table
}

type Table struct {
	GoName string
	// +sqlgen:database=name
	DbName string
	// +sqlgen:table=name
	Name    string
	Keys    []*Column
	Columns []*Column

	t types.Type
	// +sqlgen:readonly
	Readonly    bool
	autoIncrKey *Column
}

func (t *Table) ColumnsWithoutPK() []*Column {
	columns := make([]*Column, 0, len(t.Columns))
	nameMap := make(map[string]struct{})
	for _, k := range t.Keys {
		nameMap[k.Name] = struct{}{}
	}
	for _, c := range t.Columns {
		if _, ok := nameMap[c.Name]; ok || c.Readonly || c == t.autoIncrKey {
			continue
		}
		columns = append(columns, c)
	}
	return columns
}

func (t *Table) GoPtrPaths() []*FieldPath {
	paths := make([]*FieldPath, 0)
	pathMap := make(map[string]struct{})
	for _, f := range t.Columns {
		for _, p := range f.GoPtrPaths() {
			if _, ok := pathMap[p.GoPath]; ok {
				continue
			}
			paths = append(paths, p)
			pathMap[p.GoPath] = struct{}{}
		}
	}
	clear(pathMap)
	return paths
}

func (t *Table) HasPK() bool {
	return t.autoIncrKey != nil || len(t.Keys) > 0
}

func (t *Table) AutoIncrKey() (*Column, bool) {
	return t.autoIncrKey, t.autoIncrKey != nil
}

func (t *Table) InsertColumns() []*Column {
	return lo.Filter(t.Columns, func(v *Column, _ int) bool {
		return !v.Readonly && v != t.autoIncrKey
	})
}

func (t *Table) Implements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(t.t, T, true)
}

func (t *Table) PtrImplements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(types.NewPointer(t.t), T, true)
}

type Column struct {
	Name     string
	GoName   string
	GoPath   string
	Type     types.Type
	Pos      int
	Readonly bool
	// Store the position of structField
	goField *structField
	goSize  uint8
}

func (c *Column) GoPtrPaths() []*FieldPath {
	fieldPaths := c.goField.fieldPaths()
	paths := make([]*FieldPath, 0, len(fieldPaths))
	var goPath string
	for _, f := range fieldPaths {
		goPath += "." + f.name
		if f.IsPtr() {
			paths = append(paths, &FieldPath{
				Type:   f.t,
				GoName: f.name,
				GoPath: goPath,
			})
		}
	}
	return paths
}

func (c *Column) IsPtr() bool {
	return isPtr(c.Type)
}

func (c *Column) IsUnderlyingPtr() bool {
	for _, f := range c.goField.paths {
		if f.IsPtr() {
			return true
		}
	}
	return false
}

func (c *Column) IsNullable() bool {
	switch c.Type.(type) {
	case *types.Pointer,
		*types.Map,
		*types.Chan,
		*types.Interface,
		*types.Slice:
		return true
	default:
		return false
	}
}

func (c *Column) Size() uint8 {
	if c.goSize > 0 {
		return c.goSize
	}
	return 0
}

type structCache struct {
	name *ast.Ident
	doc  *ast.CommentGroup
	t    *ast.StructType
	pkg  *packages.Package
}

type structField struct {
	name string
	t    types.Type

	index    []int
	exported bool
	embedded bool
	parent   *structField
	tag      reflect.StructTag
	paths    []*structField
}

type FieldPath struct {
	GoName string
	GoPath string
	Type   types.Type
}

func (s *structField) FullGoPath() string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	paths := s.fieldPaths()
	for i := 0; i < len(paths); i++ {
		blr.WriteString("." + paths[i].name)
	}
	return blr.String()
}

func (s *structField) IsPtr() bool {
	return isPtr(s.t)
}

// // PtrPaths is helpful for initialise the addr of pointer
// func (s *structField) PtrPaths() []*FieldPath {
// 	paths := make([]*FieldPath, 0)
// 	fields := s.fieldPaths()
// 	path := ""
// 	for _, f := range fields {
// 		if isPtr(f.Type) {
// 			path += "." + f.GoName
// 			log.Println(path, f.Type.String())
// 			paths = append(paths, &FieldPath{
// 				GoName: f.GoName,
// 				Type:   f.Type,
// 			})
// 		}
// 	}
// 	return paths
// }

func (s *structField) fieldPaths() []*structField {
	if len(s.paths) > 0 {
		return s.paths
	}
	s.paths = append(s.paths, s)
	p := s.parent
	for p != nil {
		s.paths = append(s.paths, p)
		p = p.parent
	}
	reverse(s.paths)
	return s.paths
}

type typeQueue struct {
	idx  []int
	prev *structField
	doc  *ast.CommentGroup
	t    *ast.StructType
	pkg  *packages.Package
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func isPtr(v any) bool {
	_, ok := v.(*types.Pointer)
	return ok
}
