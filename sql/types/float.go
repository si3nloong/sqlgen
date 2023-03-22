package types

import (
	"database/sql"
	"database/sql/driver"
	"strconv"

	"golang.org/x/exp/constraints"
)

type FloatLike[T constraints.Float] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*FloatLike[float32])(nil)
	_ driver.Valuer = (*FloatLike[float32])(nil)
)

// Float returns a sql.Scanner
func Float[T constraints.Float](addr *T, strict ...bool) *FloatLike[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &FloatLike[T]{addr: addr, strictType: strictType}
}

func (f *FloatLike[T]) Value() (driver.Value, error) {
	if f.addr == nil {
		return nil, nil
	}
	return float64(*f.addr), nil
}

func (f *FloatLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case float64:
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
			case []byte:
				f, err := strconv.ParseFloat(string(vi), 64)
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
