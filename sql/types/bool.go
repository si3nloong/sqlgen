package types

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

type boolLike[T ~bool] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*boolLike[bool])(nil)
	_ driver.Valuer = (*boolLike[bool])(nil)
)

// Bool returns a sql.Scanner
func Bool[T ~bool](addr *T, strict ...bool) boolLike[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return boolLike[T]{addr: addr, strictType: strictType}
}

func (b boolLike[T]) Interface() T {
	if b.addr == nil {
		return *new(T)
	}
	return *b.addr
}

func (b boolLike[T]) Value() (driver.Value, error) {
	if b.addr == nil {
		return nil, nil
	}
	return bool(*b.addr), nil
}

func (b boolLike[T]) Scan(v any) error {
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
