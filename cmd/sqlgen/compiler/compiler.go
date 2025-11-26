package compiler

import (
	"errors"
	"go/ast"
	"go/types"
	"iter"
	"reflect"
	"regexp"

	"golang.org/x/tools/go/packages"
)

var (
	ErrSkip = errors.New(`compiler: skip code generation`)

	nameRegex = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9\_]*$`)
)

type Matcher interface {
	Match(v string) bool
}

type PK interface {
	pk()
}

type Table struct {
	GoName string
	// +sql:database=name
	DbName string
	// +sql:table=name
	Name string
	// +sql:readonly
	Readonly bool
	Columns  []Column
	// pk can be either `AutoIncrPrimaryKey`, `PrimaryKey` or `CompositePrimaryKey`
	pk PK
	t  types.Type
}

func (t *Table) InsertColumns() []*BasicColumn {
	columns := make([]*BasicColumn, 0, len(t.Columns))
	switch k := t.pk.(type) {
	case *AutoIncrPrimaryKey:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				if k.BasicColumn == v {
					continue
				} else if v.Readonly {
					continue
				}
				columns = append(columns, v)
			}
		}
	case *PrimaryKey:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				if v.Readonly {
					continue
				}
				columns = append(columns, v)
			}
		}
	case *CompositePrimaryKey:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				if v.Readonly {
					continue
				}
				columns = append(columns, v)
			}
		}
	default:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				if v.Readonly {
					continue
				}
				columns = append(columns, v)
			}
		}
	}
	return columns
}

func (t *Table) ColumnsExceptPK() []*BasicColumn {
	columns := make([]*BasicColumn, 0, len(t.Columns))
	switch k := t.pk.(type) {
	case *AutoIncrPrimaryKey:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				if k.BasicColumn != v {
					columns = append(columns, v)
				}
			}
		}
	case *PrimaryKey:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				if k.BasicColumn != v {
					columns = append(columns, v)
				}
			}
		}
	case *CompositePrimaryKey:
	loop:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				for _, j := range k.Columns {
					if j == v {
						continue loop
					}
				}
				columns = append(columns, v)
			}
		}
	default:
		for _, c := range t.Columns {
			switch v := c.(type) {
			case *BasicColumn:
				columns = append(columns, v)
			}
		}
	}
	return columns
}

func (t *Table) ColumnGoPtrPaths() iter.Seq[GoDecl] {
	return func(yield func(GoDecl) bool) {
		uniqueMap := make(map[string]struct{})
		defer clear(uniqueMap)
		for _, c := range t.Columns {
			for _, paths := range c.GoPtrPaths() {
				for len(paths) > 0 {
					p := paths[0]
					paths = paths[1:]
					if _, ok := uniqueMap[p.GoPath()]; ok {
						continue
					}
					if !yield(p) {
						return
					}
					uniqueMap[p.GoPath()] = struct{}{}
				}
			}
		}
	}
}

type AutoIncrPrimaryKey struct {
	*BasicColumn
}

func (AutoIncrPrimaryKey) pk() {}

type PrimaryKey struct {
	*BasicColumn
}

func (PrimaryKey) pk() {}

type CompositePrimaryKey struct {
	Columns []*BasicColumn
}

func (CompositePrimaryKey) pk() {}

func (t *Table) PK() (PK, bool) {
	if t.pk == nil {
		return nil, false
	}
	return t.pk, true
}
func (t *Table) MustPK() PK {
	if t.pk == nil {
		panic("missing primary key")
	}
	return t.pk
}
func (t *Table) ReadColumns() []*BasicColumn {
	columns := make([]*BasicColumn, 0, len(t.Columns))
	for _, c := range t.Columns {
		switch v := c.(type) {
		case *BasicColumn:
			columns = append(columns, v)
		default:
			continue
		}
	}
	return columns
}
func (t *Table) Implements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(t.t, T, true)
}
func (t *Table) PtrImplements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(types.NewPointer(t.t), T, true)
}

// func (t *Table) BasicColumnsWithoutPK() []*BasicColumn {
// 	BasicColumns := make([]*BasicColumn, 0, len(t.BasicColumns))
// 	nameMap := make(map[string]struct{})
// 	for _, k := range t.Keys {
// 		nameMap[k.Name] = struct{}{}
// 	}
// 	for _, c := range t.BasicColumns {
// 		if _, ok := nameMap[c.Name]; ok || c.Readonly || c == t.autoIncrKey {
// 			continue
// 		}
// 		BasicColumns = append(BasicColumns, c)
// 	}
// 	return BasicColumns
// }

// func (t *Table) HasPK() bool {
// 	return t.autoIncrKey != nil || len(t.Keys) > 0
// }

// func (t *Table) AutoIncrKey() (*BasicColumn, bool) {
// 	return t.autoIncrKey, t.autoIncrKey != nil
// }

// func (t *Table) InsertBasicColumns() []*BasicColumn {
// 	return lo.Filter(t.BasicColumns, func(v *BasicColumn, _ int) bool {
// 		return !v.Readonly && v != t.autoIncrKey
// 	})
// }

type structCache struct {
	name *ast.Ident
	doc  *ast.CommentGroup
	t    *ast.StructType
	pkg  *packages.Package
}

type structField struct {
	name     string
	t        types.Type
	index    []int
	exported bool
	embedded bool
	parent   *structField
	tag      reflect.StructTag
}

type typeQueue struct {
	idx  []int
	prev *structField
	doc  *ast.CommentGroup
	t    *ast.StructType
	pkg  *packages.Package
}

func isPtr(v any) bool {
	_, ok := v.(*types.Pointer)
	return ok
}
