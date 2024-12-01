package encoding

import (
	"database/sql"
	"fmt"
	"strconv"
	"unsafe"
)

func BoolScanner[T ~bool, Ptr addrOrPtr[T]](addr Ptr, strict ...bool) sql.Scanner {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &boolScanner[T, Ptr]{
		addr:       addr,
		strictType: strictType,
	}
}

type boolScanner[T ~bool, Ptr addrOrPtr[T]] struct {
	addr       Ptr
	strictType bool
}

// Scan implements the sql.Scanner interface.
func (b *boolScanner[T, Ptr]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case nil:
		switch any(b.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(b.addr)) = nil
			return nil
		case *T:
			*(*T)(unsafe.Pointer(b.addr)) = T(false)
			return nil
		default:
			panic("unreachable")
		}
	case bool:
		val := T(vi)
		switch any(b.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(b.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(b.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := T(vi != 0)
		switch any(b.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(b.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(b.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		if b.strictType {
			return fmt.Errorf(`sequel/types: unable to scan %T to ~bool`, vi)
		}

		switch vi := v.(type) {
		case []byte:
			f, err := strconv.ParseBool(unsafe.String(unsafe.SliceData(vi), len(vi)))
			if err != nil {
				return err
			}
			val = T(f)
			switch any(b.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(b.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(b.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		case string:
			f, err := strconv.ParseBool(vi)
			if err != nil {
				return err
			}
			val := T(f)
			switch any(b.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(b.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(b.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		default:
			return fmt.Errorf(`sequel/types: unable to scan %T to ~bool`, vi)
		}
	}
}
