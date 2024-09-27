package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type intLike[T constraints.Integer] struct {
	addr       *T
	strictType bool
}

func Integer[T constraints.Integer](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &intLike[T]{addr: addr, strictType: strictType}
}

func (i intLike[T]) Interface() T {
	if i.addr == nil {
		return *new(T)
	}
	return *i.addr
}

func (i intLike[T]) Value() (driver.Value, error) {
	if i.addr == nil {
		return nil, nil
	}
	return int64(*i.addr), nil
}

func (i *intLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case []byte:
		m, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(vi), len(vi)), 10, 64)
		if err != nil {
			return err
		}
		val = T(m)
	case int64:
		val = T(vi)
	case nil:
		i.addr = nil
		return nil

	default:
		if i.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~int`, vi)
		}

		switch vi := v.(type) {
		case string:
			m, err := strconv.ParseInt(string(vi), 10, 64)
			if err != nil {
				return err
			}
			val = T(m)
		case uint64:
			val = T(vi)
		case uint32:
			val = T(vi)
		case uint16:
			val = T(vi)
		case uint8:
			val = T(vi)
		case uint:
			val = T(vi)
		case int32:
			val = T(vi)
		case int16:
			val = T(vi)
		case int8:
			val = T(vi)
		case int:
			val = T(vi)
		case float32:
			val = T(vi)
		case float64:
			val = T(vi)
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~int`, vi)
		}
	}
	*i.addr = val
	return nil
}
