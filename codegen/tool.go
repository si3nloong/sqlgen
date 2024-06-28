package codegen

import (
	"fmt"
	"go/types"
)

func UnderlyingType(t types.Type) (codec *Mapping, typeStr string) {
	var (
		prev = t
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
		case *types.Slice:
			typeStr += "[]"
			prev = v.Elem()
		case *types.Array:
			typeStr += fmt.Sprintf("[%d]", v.Len())
			prev = v.Elem()
		default:
			break loop
		}
		if v, ok := typeMap[typeStr]; ok {
			return v, typeStr
		}
		if prev == t {
			break loop
		}
		t = prev
	}
	if v, ok := typeMap[typeStr]; ok {
		return v, typeStr
	}
	return nil, typeStr
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
