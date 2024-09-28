package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func Uint8[T ~uint8](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeUint(addr, 8, strict...)
}

func Uint16[T ~uint16](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeUint(addr, 16, strict...)
}

func Uint32[T ~uint32](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeUint(addr, 32, strict...)
}

func Uint64[T ~uint64](addr *T, strict ...bool) ValueScanner[T] {
	return newFixedSizeUint(addr, 64, strict...)
}

func newFixedSizeUint[T constraints.Unsigned](addr *T, bitSize int, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &fixedSizeUintLike[T]{addr: addr, bitSize: bitSize, strictType: strictType}
}

type fixedSizeUintLike[T constraints.Unsigned] struct {
	addr       *T
	bitSize    int
	strictType bool
}

func (i fixedSizeUintLike[T]) Interface() T {
	if i.addr == nil {
		return *new(T)
	}
	return *i.addr
}

func (i fixedSizeUintLike[T]) Value() (driver.Value, error) {
	if i.addr == nil {
		return nil, nil
	}
	return int64(*i.addr), nil
}

func (i *fixedSizeUintLike[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case []byte:
		m, err := strconv.ParseUint(unsafe.String(unsafe.SliceData(vi), len(vi)), 10, i.bitSize)
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
			return fmt.Errorf(`sequel/types: unable to scan %T to ~uint`, vi)
		}

		switch vi := v.(type) {
		case string:
			m, err := strconv.ParseUint(vi, 10, i.bitSize)
			if err != nil {
				return err
			}
			val = T(m)
		case float64:
			val = T(vi)
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~uint`, vi)
		}
	}
	*i.addr = val
	return nil
}
