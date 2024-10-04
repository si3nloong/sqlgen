package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type intLike[T constraints.Integer] struct {
	addr       *T
	strictType bool
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
		m, err := strconv.Atoi(unsafe.String(unsafe.SliceData(vi), len(vi)))
		if err != nil {
			return err
		}
		val = T(m)
	case int64:
		val = T(vi)
	case uint64:
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
			m, err := strconv.Atoi(vi)
			if err != nil {
				return err
			}
			val = T(m)
		case float64:
			val = T(vi)
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~int`, vi)
		}
	}
	*i.addr = val
	return nil
}

func Integer[T constraints.Integer](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &intLike[T]{addr: addr, strictType: strictType}
}

func Int8[T ~int8](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeInt(addr, 8, strict...)
}

func Int16[T ~int16](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeInt(addr, 16, strict...)
}

func Int32[T ~int32](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeInt(addr, 32, strict...)
}

func Int64[T ~int64](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeInt(addr, 64, strict...)
}

func newFixedSizeInt[T constraints.Signed](addr *T, bitSize int, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &fixedSizeIntLike[T]{addr: addr, bitSize: bitSize, strictType: strictType}
}

type fixedSizeIntLike[T constraints.Signed] struct {
	addr       *T
	bitSize    int
	strictType bool
}

func (i fixedSizeIntLike[T]) Interface() T {
	if i.addr == nil {
		return *new(T)
	}
	return *i.addr
}

func (i fixedSizeIntLike[T]) Value() (driver.Value, error) {
	if i.addr == nil {
		return nil, nil
	}
	return int64(*i.addr), nil
}

func (i *fixedSizeIntLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case nil:
		i.addr = nil
		return nil
	case []byte:
		m, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(vi), len(vi)), 10, i.bitSize)
		if err != nil {
			return err
		}
		val = T(m)
	case int64:
		val = T(vi)
	case uint64:
		val = T(vi)
	default:
		if i.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~int`, vi)
		}

		switch vi := v.(type) {
		case string:
			m, err := strconv.ParseInt(vi, 10, i.bitSize)
			if err != nil {
				return err
			}
			val = T(m)
		case float64:
			val = T(vi)
		case time.Time:
			val = T(vi.Unix())
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~int`, vi)
		}
	}
	*i.addr = val
	return nil
}
