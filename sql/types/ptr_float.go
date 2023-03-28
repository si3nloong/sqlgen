package types

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type ptrOfFloatLike[T constraints.Float] struct {
	addr **T
}

func PtrOfFloat[T constraints.Float](v **T) ptrOfFloatLike[T] {
	return ptrOfFloatLike[T]{addr: v}
}

func (p ptrOfFloatLike[T]) Interface() *T {
	if p.addr == nil {
		return nil
	}
	return *p.addr
}

func (p ptrOfFloatLike[T]) Scan(v any) error {
	if v == nil {
		(*p.addr) = nil
		return nil
	}

	switch vi := v.(type) {
	case []byte:
		f, err := strconv.ParseFloat(string(vi), 64)
		if err != nil {
			return err
		}
		val := T(f)
		*p.addr = &val
	case float64:
		val := T(vi)
		*p.addr = &val
	case int64:
		val := T(vi)
		*p.addr = &val
	}
	return nil
}
