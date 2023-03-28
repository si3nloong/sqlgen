package types

import (
	"time"
)

type ptrOfTime[T time.Time] struct {
	addr **T
}

func PtrOfTime[T time.Time](v **T) ptrOfTime[T] {
	return ptrOfTime[T]{addr: v}
}

func (p ptrOfTime[T]) Interface() *T {
	if p.addr == nil {
		return nil
	}
	return *p.addr
}

func (p ptrOfTime[T]) Scan(v any) error {
	if v == nil {
		(*p.addr) = nil
		return nil
	}

	switch vi := v.(type) {
	case time.Time:
		val := T(vi)
		*p.addr = &val
	case string:
		b, err := time.Parse(time.RFC3339Nano, vi)
		if err != nil {
			return err
		}
		val := T(b)
		*p.addr = &val
	}
	return nil
}
