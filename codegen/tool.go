package codegen

import (
	"fmt"
	"go/types"
	"regexp"
	"strconv"
)

type GoType interface {
	String() string
}

type GoArray interface {
	Len() int
}

type nNode struct {
	t string
}

func (n nNode) String() string {
	return n.t
}

type arrNode struct {
	t    string
	size int
}

func (n arrNode) String() string {
	return n.t
}
func (n arrNode) Len() int {
	return n.size
}

var (
	arrayRegexp = regexp.MustCompile(`^\[(\d+)\](rune|string|byte)$`)
)

func UnderlyingType(t types.Type) (*Mapping, GoType) {
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
			return v, &nNode{t: typeStr}
		}
		if prev == t {
			break loop
		}
		t = prev
	}
	if v, ok := typeMap[typeStr]; ok {
		return v, &nNode{t: typeStr}
	}
	// Find fixed size array mapper
	if matches := arrayRegexp.FindStringSubmatch(typeStr); len(matches) > 0 {
		if v, ok := typeMap["[...]"+matches[2]]; ok {
			len, _ := strconv.Atoi(matches[1])
			return v, &arrNode{t: typeStr, size: len}
		}
	}
	return nil, &nNode{t: typeStr}
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
