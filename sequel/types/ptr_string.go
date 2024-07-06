package types

import "fmt"

type ptrOfStrLike[T StringLikeType] struct {
	addr **T
}

func PtrOfString[T StringLikeType](v **T) ptrOfStrLike[T] {
	return ptrOfStrLike[T]{addr: v}
}

func (p ptrOfStrLike[T]) Interface() *T {
	if p.addr == nil {
		return nil
	}
	return *p.addr
}

func (p ptrOfStrLike[T]) Scan(v any) error {
	if v == nil {
		(*p.addr) = nil
		return nil
	}

	switch vi := v.(type) {
	case string:
		val := T(vi)
		*p.addr = &val
	case []byte:
		val := T(vi)
		*p.addr = &val
	default:
		return fmt.Errorf(`types: unable to scan %T to *string`, vi)
	}
	return nil
}
