package codegen

import (
	"fmt"
	"go/types"
	"regexp"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

var (
	arrayRegexp = regexp.MustCompile(`^\[(\d+)\](rune|string|byte)$`)
)

func (g *Generator) columnDataType(t types.Type) (*dialect.ColumnType, bool) {
	var (
		prev    = t
		typeStr string
	)

loop:
	for prev != nil {
		switch v := prev.(type) {
		case *types.Basic:
			typeStr += v.String()
			break loop
		case *types.Named:
			if _, ok := v.Underlying().(*types.Struct); ok {
				typeStr += v.String()
				break loop
			}
			prev = v.Underlying()
		case *types.Pointer:
			typeStr += "*"
			prev = v.Elem()
			continue
		case *types.Slice:
			typeStr += "[]"
			prev = v.Elem()
			continue
		case *types.Array:
			typeStr += fmt.Sprintf("[%d]", v.Len())
			prev = v.Elem()
			continue
		case *types.Alias:
			prev = v.Rhs()
			continue
		default:
			break loop
		}
		if v, ok := g.defaultColumnTypes[typeStr]; ok {
			return v, true
		}
		if prev == t {
			break loop
		}
		t = prev
	}
	if v, ok := g.defaultColumnTypes[typeStr]; ok {
		return v, true
	}
	// Find fixed size array mapper
	if matches := arrayRegexp.FindStringSubmatch(typeStr); len(matches) > 0 {
		if v, ok := g.defaultColumnTypes["[...]"+matches[2]]; ok {
			// len, _ := strconv.Atoi(matches[1])
			return v, true
		}
	}
	return nil, false
}

func assertAsPtr[T any](v any) *T {
	t, ok := v.(*T)
	if ok {
		return t
	}
	return nil
}

func arraySize(t types.Type) int64 {
	for t != nil {
		switch vt := t.(type) {
		case *types.Pointer:
			t = vt.Elem()
		case *types.Array:
			return vt.Len()
		default:
			return -1
		}
	}
	return -1
}

func isImplemented(t types.Type, iv *types.Interface) bool {
	method, wrongType := types.MissingMethod(t, iv, true)
	return method == nil && !wrongType
}

func pointerType(t types.Type) (*types.Pointer, bool) {
	v, ok := t.(*types.Pointer)
	if ok {
		return v, true
	}
	return types.NewPointer(t), false
}

func underlyingType(t types.Type) (types.Type, bool) {
	v, ok := t.(*types.Pointer)
	if ok {
		return v.Elem(), true
	}
	return t, false
}
