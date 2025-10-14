package compiler

import (
	"fmt"
	"go/types"
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
	switch v := any(g.Type).(type) {
	case *types.Basic:
		return v.String()
	case *types.Named:
		return ""
	case *types.Array:
		return "[...]"
	case nil:
		return "*"
	default:
		return g.Type.String()
	}
}

var (
	Rune         = GoType{types.Typ[types.Rune]}
	Byte         = GoType{types.Typ[types.Byte]}
	String       = GoType{types.Typ[types.String]}
	Bool         = GoType{types.Typ[types.Bool]}
	Int          = GoType{types.Typ[types.Int]}
	Int8         = GoType{types.Typ[types.Int8]}
	Int16        = GoType{types.Typ[types.Int16]}
	Int32        = GoType{types.Typ[types.Int32]}
	Int64        = GoType{types.Typ[types.Int64]}
	Uint         = GoType{types.Typ[types.Uint]}
	Uint8        = GoType{types.Typ[types.Uint8]}
	Uint16       = GoType{types.Typ[types.Uint16]}
	Uint32       = GoType{types.Typ[types.Uint32]}
	Uint64       = GoType{types.Typ[types.Uint64]}
	Float32      = GoType{types.Typ[types.Float32]}
	Float64      = GoType{types.Typ[types.Float64]}
	Time         = GoType{types.Typ[types.Float32]}
	RuneArray    = GoType{types.NewArray(types.Typ[types.Rune], 1)}
	ByteArray    = GoType{types.NewArray(types.Typ[types.Byte], 1)}
	StringArray  = GoType{types.NewArray(types.Typ[types.String], 1)}
	BoolArray    = GoType{types.NewArray(types.Typ[types.Bool], 1)}
	IntArray     = GoType{types.NewArray(types.Typ[types.Int], 1)}
	Int8Array    = GoType{types.NewArray(types.Typ[types.Int8], 1)}
	Int16Array   = GoType{types.NewArray(types.Typ[types.Int16], 1)}
	Int32Array   = GoType{types.NewArray(types.Typ[types.Int32], 1)}
	Int64Array   = GoType{types.NewArray(types.Typ[types.Int64], 1)}
	UintArray    = GoType{types.NewArray(types.Typ[types.Uint], 1)}
	Uint8Array   = GoType{types.NewArray(types.Typ[types.Uint8], 1)}
	Uint16Array  = GoType{types.NewArray(types.Typ[types.Uint16], 1)}
	Uint32Array  = GoType{types.NewArray(types.Typ[types.Uint32], 1)}
	Uint64Array  = GoType{types.NewArray(types.Typ[types.Uint64], 1)}
	Float32Array = GoType{types.NewArray(types.Typ[types.Float32], 1)}
	Float64Array = GoType{types.NewArray(types.Typ[types.Float64], 1)}
	// TimeArray    GoType = "[...]time.Time"
	Runes        = GoType{types.NewSlice(types.Typ[types.Rune])}
	Bytes        = GoType{types.NewSlice(types.Typ[types.Byte])}
	StringSlice  = GoType{types.NewSlice(types.Typ[types.String])}
	BoolSlice    = GoType{types.NewSlice(types.Typ[types.Bool])}
	IntSlice     = GoType{types.NewSlice(types.Typ[types.Int])}
	Int8Slice    = GoType{types.NewSlice(types.Typ[types.Int8])}
	Int16Slice   = GoType{types.NewSlice(types.Typ[types.Int16])}
	Int32Slice   = GoType{types.NewSlice(types.Typ[types.Int32])}
	Int64Slice   = GoType{types.NewSlice(types.Typ[types.Int64])}
	UintSlice    = GoType{types.NewSlice(types.Typ[types.Uint])}
	Uint8Slice   = GoType{types.NewSlice(types.Typ[types.Uint8])}
	Uint16Slice  = GoType{types.NewSlice(types.Typ[types.Uint16])}
	Uint32Slice  = GoType{types.NewSlice(types.Typ[types.Uint32])}
	Uint64Slice  = GoType{types.NewSlice(types.Typ[types.Uint64])}
	Float32Slice = GoType{types.NewSlice(types.Typ[types.Float32])}
	Float64Slice = GoType{types.NewSlice(types.Typ[types.Float64])}
	Any          = GoType{}
	// TimeSlice    GoType = "[]time.Time"
)
