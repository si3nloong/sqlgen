package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type StringLikeType interface {
	~string | ~[]byte
}

type strLike[T StringLikeType] struct {
	addr       *T
	strictType bool
}

func String[T StringLikeType](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &strLike[T]{addr: addr, strictType: strictType}
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

func (s *strLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case nil:
		s.addr = nil
		return nil
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
		case float64:
			val = T(strconv.FormatFloat(vi, 'f', -1, 64))
		case time.Time:
			val = T(vi.String())
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~string`, vi)
		}
	}
	*s.addr = val
	return nil
}
