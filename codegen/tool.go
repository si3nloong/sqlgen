package codegen

import (
	"go/types"
)

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

func assertAsPtr[T any](v any) *T {
	t, ok := v.(*T)
	if ok {
		return t
	}
	return nil
}

func isImplemented(t types.Type, iv *types.Interface) bool {
	method, wrongType := types.MissingMethod(t, iv, true)
	return method == nil && !wrongType
}

func newPointer(t types.Type) *types.Pointer {
	v, ok := t.(*types.Pointer)
	if ok {
		return v
	}
	return types.NewPointer(t)
}
