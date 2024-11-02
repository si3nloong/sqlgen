package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"
)

type boolLike[T ~bool] struct {
	addr       *T
	strictType bool
}

// Bool returns a sql.Scanner
func Bool[T ~bool](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &boolLike[T]{addr: addr, strictType: strictType}
}

func (b boolLike[T]) Interface() T {
	if b.addr == nil {
		return *new(T)
	}
	return *b.addr
}

// Value implements the driver.Valuer interface.
func (b boolLike[T]) Value() (driver.Value, error) {
	if b.addr == nil {
		return nil, nil
	}
	return (bool)(*b.addr), nil
}

// Scan implements the sql.Scanner interface.
func (b *boolLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case bool:
		val = T(vi)
	case int64:
		val = T(vi != 0)
	case nil:
		b.addr = nil
		return nil
	default:
		if b.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~bool`, vi)
		}

		switch vi := v.(type) {
		case []byte:
			f, err := strconv.ParseBool(unsafe.String(unsafe.SliceData(vi), len(vi)))
			if err != nil {
				return err
			}
			val = T(f)
		case string:
			f, err := strconv.ParseBool(vi)
			if err != nil {
				return err
			}
			val = T(f)
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~bool`, vi)
		}
	}
	*b.addr = val
	return nil
}
