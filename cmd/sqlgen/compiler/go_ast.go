package compiler

import (
	"go/types"
	"strings"
)

type GoDecl interface {
	GoName() string
	GoType() types.Type
	GoIndexes() []int
	GoPath() string
	IsGoPtr() bool
	// This will return all the pointer paths from the root to this field
	// eg. A.nested.Field []*A -> *nested -> Field
	// will return [[*A], [*nested]]
	GoPtrPaths() [][]GoDecl
}

type Column interface {
	GoDecl
	IsUnderlyingPtr() bool
	Pos() int
	Name() string
	columnType()
}

type goInfo struct {
	parent    *goInfo
	goName    string
	goIndexes []int
	goType    types.Type
}

func (c goInfo) GoName() string { return c.goName }
func (c goInfo) GoPath() string {
	paths := []string{c.goName}
	t := c.parent
	for t != nil {
		paths = append([]string{t.goName}, paths...)
		t = t.parent
	}
	return strings.Join(paths, ".")
}
func (c goInfo) GoType() types.Type { return c.goType }
func (c goInfo) GoIndexes() []int   { return c.goIndexes }
func (c goInfo) IsGoPtr() bool {
	return isPtr(c.goType)
}
func (c goInfo) GoPtrPaths() [][]GoDecl {
	result := [][]GoDecl{}
	paths := []GoDecl{}
	if c.IsGoPtr() {
		paths = append(paths, c)
		result = append(result, paths)
		paths = []GoDecl{}
	}
	t := c.parent
	for t != nil {
		// prepend
		paths = append([]GoDecl{t}, paths...)
		if t.IsGoPtr() {
			// prepend
			result = append([][]GoDecl{paths}, result...)
			paths = []GoDecl{}
		}
		t = t.parent
	}
	return result
}
func (c goInfo) IsUnderlyingPtr() bool {
	if c.IsGoPtr() {
		return true
	}
	t := c.parent
	for t != nil {
		if t.IsGoPtr() {
			return true
		}
		t = t.parent
	}
	return false
}
func (c goInfo) IsNullable() bool {
	switch c.goType.(type) {
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

func mapToGoInfo(s *structField) *goInfo {
	t := goInfo{}
	t.goName = s.name
	t.goIndexes = append(t.goIndexes, s.index...)
	t.goType = s.t
	return &t
}

type BasicColumn struct {
	*goInfo
	Readonly bool
	name     string
	pos      int
}

func (c BasicColumn) GoSize() (int64, bool) {
	if v, ok := c.goType.(*types.Array); ok {
		return v.Len(), true
	}
	return 0, false
}
func (c BasicColumn) Name() string { return c.name }
func (c BasicColumn) Pos() int     { return c.pos }
func (BasicColumn) columnType()    {}

type GeneratedColumn struct {
	*goInfo
	pos  int
	name string
}

func (c GeneratedColumn) Name() string { return c.name }
func (c GeneratedColumn) Pos() int     { return c.pos }
func (GeneratedColumn) columnType()    {}
