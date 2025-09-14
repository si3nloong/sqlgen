package encoding

import (
	"database/sql"
	"fmt"
	"strconv"
	"unsafe"
)

// Float returns a sql.Scanner
func Float32Scanner[T ~float32, Addr addrOrPtr[T]](addr Addr, strict ...bool) sql.Scanner {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &float32Scanner[T, Addr]{
		addr:       addr,
		strictType: strictType,
	}
}

// Float returns a sql.Scanner
func Float64Scanner[T ~float64, Addr addrOrPtr[T]](addr Addr, strict ...bool) sql.Scanner {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &float64Scanner[T, Addr]{
		addr:       addr,
		strictType: strictType,
	}
}

type float32Scanner[T ~float32, Addr addrOrPtr[T]] struct {
	addr       Addr
	strictType bool
}

func (f *float32Scanner[T, Addr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch any(f.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(f.addr)) = nil
			return nil
		case *T:
			var v T
			*(*T)(unsafe.Pointer(f.addr)) = v
			return nil
		default:
			panic("unreachable")
		}
	case float64:
		val := T(vi)
		switch any(f.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(f.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(f.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := T(vi)
		switch any(f.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(f.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(f.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		if f.strictType {
			return fmt.Errorf(`sequel/encoding: unable to scan %T to ~float32`, vi)
		}

		switch vi := v.(type) {
		case []byte:
			f32, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(vi), len(vi)), 32)
			if err != nil {
				return err
			}
			val := T(f32)
			switch any(f.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(f.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(f.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		case string:
			f32, err := strconv.ParseFloat(vi, 32)
			if err != nil {
				return err
			}
			val := T(f32)
			switch any(f.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(f.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(f.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		default:
			return fmt.Errorf(`sequel/encoding: unable to scan %T to ~float32`, vi)
		}
	}
}

type float64Scanner[T ~float64, Addr addrOrPtr[T]] struct {
	addr       Addr
	strictType bool
}

func (f *float64Scanner[T, Addr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch any(f.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(f.addr)) = nil
			return nil
		case *T:
			var v T
			*(*T)(unsafe.Pointer(f.addr)) = v
			return nil
		default:
			panic("unreachable")
		}
	case float64:
		val := T(vi)
		switch any(f.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(f.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(f.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := T(vi)
		switch any(f.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(f.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(f.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		if f.strictType {
			return fmt.Errorf(`sequel/encoding: unable to scan %T to ~float64`, vi)
		}

		switch vi := v.(type) {
		case []byte:
			f64, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(vi), len(vi)), 64)
			if err != nil {
				return err
			}
			val := T(f64)
			switch any(f.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(f.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(f.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		case string:
			f64, err := strconv.ParseFloat(vi, 64)
			if err != nil {
				return err
			}
			val := T(f64)
			switch any(f.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(f.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(f.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		default:
			return fmt.Errorf(`sequel/encoding: unable to scan %T to ~float64`, vi)
		}
	}
}
