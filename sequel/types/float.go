package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"
)

type float32Like[T ~float32] struct {
	addr       *T
	strictType bool
}

// Float returns a sql.Scanner
func Float32[T ~float32](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &float32Like[T]{addr: addr, strictType: strictType}
}

func (f float32Like[T]) Interface() T {
	if f.addr == nil {
		return *new(T)
	}
	return *f.addr
}

func (f float32Like[T]) Value() (driver.Value, error) {
	if f.addr == nil {
		return nil, nil
	}
	return float64(*f.addr), nil
}

func (f *float32Like[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case float64:
		val = T(vi)
	case int64:
		val = T(vi)
	case nil:
		f.addr = nil
		return nil
	default:
		if f.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~float32`, vi)
		}

		switch vi := v.(type) {
		case []byte:
			f, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(vi), len(vi)), 64)
			if err != nil {
				return err
			}
			val = T(f)
		case string:
			f, err := strconv.ParseFloat(vi, 32)
			if err != nil {
				return err
			}
			val = T(f)
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~float32`, vi)
		}
	}
	*f.addr = val
	return nil
}

type float64Like[T ~float64] struct {
	addr       *T
	strictType bool
}

// Float returns a sql.Scanner
func Float64[T ~float64](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &float64Like[T]{addr: addr, strictType: strictType}
}

func (f float64Like[T]) Interface() T {
	if f.addr == nil {
		return *new(T)
	}
	return *f.addr
}

func (f float64Like[T]) Value() (driver.Value, error) {
	if f.addr == nil {
		return nil, nil
	}
	return float64(*f.addr), nil
}

func (f *float64Like[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case float64:
		val = T(vi)
	case int64:
		val = T(vi)
	case nil:
		f.addr = nil
		return nil
	default:
		if f.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~float64`, vi)
		}

		switch vi := v.(type) {
		case []byte:
			f, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(vi), len(vi)), 64)
			if err != nil {
				return err
			}
			val = T(f)
		case string:
			f, err := strconv.ParseFloat(vi, 64)
			if err != nil {
				return err
			}
			val = T(f)
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~float64`, vi)
		}
	}
	*f.addr = val
	return nil
}
