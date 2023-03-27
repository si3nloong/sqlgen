package types

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

type stringLikeType interface {
	~string | ~[]byte
}

type StringLike[T stringLikeType] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*StringLike[string])(nil)
	_ driver.Valuer = (*StringLike[string])(nil)
)

func String[T stringLikeType](addr *T, strict ...bool) *StringLike[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &StringLike[T]{addr: addr, strictType: strictType}
}

func (s StringLike[T]) Interface() T {
	if s.addr == nil {
		return *new(T)
	}
	return *s.addr
}

func (s StringLike[T]) Value() (driver.Value, error) {
	if s.addr == nil {
		return nil, nil
	}
	return string(*s.addr), nil
}

func (s *StringLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case string:
		val = T(vi)
	case []byte:
		val = T(vi)
	default:
		if !s.strictType {
			switch vi := v.(type) {
			case bool:
				val = T(strconv.FormatBool(vi))
			case int:
				val = T(strconv.Itoa(vi))
			case int64:
				val = T(strconv.FormatInt(vi, 10))
			case uint64:
				val = T(strconv.FormatUint(vi, 10))
			}
		}
	}
	*s.addr = val
	return nil
}
