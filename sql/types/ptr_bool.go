package types

import (
	"strconv"

	"github.com/si3nloong/sqlgen/internal/strfmt"
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
		b, err := strconv.ParseBool(strfmt.B2s(vi))
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
	}
	return nil
}
