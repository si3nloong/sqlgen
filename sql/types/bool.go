package types

import (
	"database/sql"
	"strconv"
)

type BoolLike[T ~bool] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner = (*BoolLike[bool])(nil)
)

// Bool returns a sql.Scanner
func Bool[T ~bool](addr *T, strict ...bool) *BoolLike[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &BoolLike[T]{addr: addr, strictType: strictType}
}

func (b *BoolLike[T]) Interface() T {
	if b.addr == nil {
		return *new(T)
	}
	return *b.addr
}

func (b *BoolLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case bool:
		val = T(vi)
	default:
		if !b.strictType {
			switch vi := v.(type) {
			case string:
				f, err := strconv.ParseBool(vi)
				if err != nil {
					return err
				}
				val = T(f)
			case []byte:
				f, err := strconv.ParseBool(string(vi))
				if err != nil {
					return err
				}
				val = T(f)
			case int64:
				val = T(vi != 0)
			}
		}
	}
	*b.addr = val
	return nil
}
