package types

import (
	"fmt"
	"strconv"

	"github.com/si3nloong/sqlgen/internal/strfmt"
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
		f, err := strconv.ParseFloat(strfmt.B2s(vi), 64)
		if err != nil {
			return err
		}
		val := T(f)
		*p.addr = &val
	case string:
		f, err := strconv.ParseFloat(vi, 64)
		if err != nil {
			return err
		}
		val := T(f)
		*p.addr = &val
	case float32:
		val := T(vi)
		*p.addr = &val
	case float64:
		val := T(vi)
		*p.addr = &val
	case int64:
		val := T(vi)
		*p.addr = &val
	default:
		return fmt.Errorf(`sqlgen: unable to scan to float`)
	}
	return nil
}
