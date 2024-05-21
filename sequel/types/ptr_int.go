package types

import (
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type ptrOfIntLike[T constraints.Integer] struct {
	addr **T
}

func PtrOfInt[T constraints.Integer](v **T) ptrOfIntLike[T] {
	return ptrOfIntLike[T]{addr: v}
}

func (p ptrOfIntLike[T]) Interface() *T {
	if p.addr == nil {
		return nil
	}
	return *p.addr
}

func (p ptrOfIntLike[T]) Scan(v any) error {
	if v == nil {
		(*p.addr) = nil
		return nil
	}

	switch vi := v.(type) {
	case []byte:
		i, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(vi), len(vi)), 10, 64)
		if err != nil {
			return err
		}
		val := T(i)
		*p.addr = &val
	case int64:
		val := T(vi)
		*p.addr = &val
	}
	return nil
}
