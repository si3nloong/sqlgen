package types

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type floatLike[T constraints.Float] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*floatLike[float32])(nil)
	_ driver.Valuer = (*floatLike[float32])(nil)
)

// Float returns a sql.Scanner
func Float[T constraints.Float](addr *T, strict ...bool) floatLike[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return floatLike[T]{addr: addr, strictType: strictType}
}

func (f floatLike[T]) Interface() T {
	if f.addr == nil {
		return *new(T)
	}
	return *f.addr
}

func (f floatLike[T]) Value() (driver.Value, error) {
	if f.addr == nil {
		return nil, nil
	}
	return float64(*f.addr), nil
}

func (f floatLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case []byte:
		f, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(vi), len(vi)), 64)
		if err != nil {
			return err
		}
		val = T(f)
	case float64:
		val = T(vi)
	case int64:
		val = T(vi)
	case uint64:
		val = T(vi)
	default:
		if !f.strictType {
			switch vi := v.(type) {
			case string:
				f, err := strconv.ParseFloat(vi, 64)
				if err != nil {
					return err
				}
				val = T(f)
			}
		}
	}
	*f.addr = val
	return nil
}
