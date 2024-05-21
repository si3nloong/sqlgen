package types

import (
	"strconv"
	"unsafe"
)

type ptrOfBoolLike[T ~bool] struct {
	addr **T
}

func PtrOfBool[T ~bool](v **T) ptrOfBoolLike[T] {
	return ptrOfBoolLike[T]{addr: v}
}

func (p ptrOfBoolLike[T]) Interface() *T {
	if p.addr == nil {
		return nil
	}
	return *p.addr
}

func (p ptrOfBoolLike[T]) Scan(v any) error {
	if v == nil {
		(*p.addr) = nil
		return nil
	}

	switch vi := v.(type) {
	case []byte:
		b, err := strconv.ParseBool(unsafe.String(unsafe.SliceData(vi), len(vi)))
		if err != nil {
			return err
		}
		val := T(b)
		*p.addr = &val
	case string:
		b, err := strconv.ParseBool(vi)
		if err != nil {
			return err
		}
		val := T(b)
		*p.addr = &val
	case bool:
		val := T(vi)
		*p.addr = &val
	case int64:
		val := T(vi != 0)
		*p.addr = &val
	}
	return nil
}
