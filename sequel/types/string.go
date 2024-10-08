package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
)

type StringLikeType interface {
	~string | ~[]byte
}

type strLike[T StringLikeType] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*strLike[string])(nil)
	_ driver.Valuer = (*strLike[string])(nil)
)

func String[T StringLikeType](addr *T, strict ...bool) strLike[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return strLike[T]{addr: addr, strictType: strictType}
}

func (s strLike[T]) Interface() T {
	if s.addr == nil {
		return *new(T)
	}
	return *s.addr
}

func (s strLike[T]) Value() (driver.Value, error) {
	if s.addr == nil {
		return nil, nil
	}
	return string(*s.addr), nil
}

func (s strLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case string:
		val = T(vi)
	case []byte:
		val = T(vi)
	default:
		if s.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~string`, vi)
		}

		switch vi := v.(type) {
		case bool:
			val = T(strconv.FormatBool(vi))
		case int64:
			val = T(strconv.FormatInt(vi, 10))
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~string`, vi)
		}
	}
	*s.addr = val
	return nil
}
