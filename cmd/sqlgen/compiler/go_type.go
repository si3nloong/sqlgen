package compiler

import (
	"fmt"
	"go/importer"
	"go/types"
	"strings"
)

type GoType struct {
	Type types.Type
}

var _ (fmt.GoStringer) = (*GoType)(nil)

// Return the actual size of array else it will panic
func (g GoType) GoSize() int64 {
	prev := g.Type
	for prev != nil {
		switch v := any(g.Type).(type) {
		case *types.Pointer:
			prev = v.Elem()
		case *types.Array:
			return v.Len()
		default:
			panic(fmt.Sprintf("invalid type %s", v))
		}
	}
	panic(fmt.Sprintf("invalid type %s", prev))
}

func (g GoType) GoString() string {
	var (
		prev    = g.Type
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
			typeStr += "[...]"
			prev = v.Elem()
			continue
		case *types.Alias:
			prev = v.Rhs()
			continue
		default:
			break loop
		}
	}
	return typeStr
}

var (
	Rune         = types.Typ[types.Rune].String()
	Byte         = types.Typ[types.Byte].String()
	Bool         = types.Typ[types.Bool].String()
	String       = types.Typ[types.String].String()
	Int          = types.Typ[types.Int].String()
	Int8         = types.Typ[types.Int8].String()
	Int16        = types.Typ[types.Int16].String()
	Int32        = types.Typ[types.Int32].String()
	Int64        = types.Typ[types.Int64].String()
	Uint         = types.Typ[types.Uint].String()
	Uint8        = types.Typ[types.Uint8].String()
	Uint16       = types.Typ[types.Uint16].String()
	Uint32       = types.Typ[types.Uint32].String()
	Uint64       = types.Typ[types.Uint64].String()
	Float32      = types.Typ[types.Float32].String()
	Float64      = types.Typ[types.Float64].String()
	Complex64    = types.Typ[types.Complex64].String()
	Complex128   = types.Typ[types.Complex128].String()
	Time         = typeNamed("time.Time").String()
	RuneSlice    = types.NewSlice(types.Typ[types.Rune]).String()
	ByteSlice    = types.NewSlice(types.Typ[types.Byte]).String()
	StringSlice  = types.NewSlice(types.Typ[types.String]).String()
	BoolSlice    = types.NewSlice(types.Typ[types.Bool]).String()
	IntSlice     = types.NewSlice(types.Typ[types.Int]).String()
	Int8Slice    = types.NewSlice(types.Typ[types.Int8]).String()
	Int16Slice   = types.NewSlice(types.Typ[types.Int16]).String()
	Int32Slice   = types.NewSlice(types.Typ[types.Int32]).String()
	Int64Slice   = types.NewSlice(types.Typ[types.Int64]).String()
	UintSlice    = types.NewSlice(types.Typ[types.Uint]).String()
	Uint8Slice   = types.NewSlice(types.Typ[types.Uint8]).String()
	Uint16Slice  = types.NewSlice(types.Typ[types.Uint16]).String()
	Uint32Slice  = types.NewSlice(types.Typ[types.Uint32]).String()
	Uint64Slice  = types.NewSlice(types.Typ[types.Uint64]).String()
	Float32Slice = types.NewSlice(types.Typ[types.Float32]).String()
	Float64Slice = types.NewSlice(types.Typ[types.Float64]).String()
	Any          = "any"
)

func typeNamed(path string) *types.Named {
	paths := strings.SplitN(path, ".", 2)
	if len(paths) != 2 {
		panic(`invalid package path`)
	}
	pkg, err := importer.Default().Import(paths[0])
	if err != nil {
		panic(err)
	}
	return pkg.Scope().Lookup(paths[1]).Type().(*types.Named)
}
