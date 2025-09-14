package encoding

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

type stringType interface {
	~string | ~[]byte
}

func StringScanner[T stringType, Ptr addrOrPtr[T]](addr Ptr, strict ...bool) sql.Scanner {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &strScanner[T, Ptr]{
		addr:       addr,
		strictType: strictType,
	}
}

type strScanner[T stringType, Ptr addrOrPtr[T]] struct {
	addr       Ptr
	strictType bool
}

func (s *strScanner[T, Ptr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch any(s.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(s.addr)) = nil
			return nil
		case *T:
			*(*T)(unsafe.Pointer(s.addr)) = T("")
			return nil
		default:
			panic("unreachable")
		}
	case string:
		val := T(vi)
		switch any(s.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(s.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(s.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	case []byte:
		val := T(vi)
		switch any(s.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(s.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(s.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		if s.strictType {
			return fmt.Errorf(`sequel/encoding: unable to scan %T to ~string`, vi)
		}

		switch vi := v.(type) {
		case bool:
			val := T(strconv.FormatBool(vi))
			switch any(s.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(s.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(s.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		case int64:
			val := T(strconv.FormatInt(vi, 10))
			switch any(s.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(s.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(s.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		case float64:
			val := T(strconv.FormatFloat(vi, 'f', -1, 64))
			switch any(s.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(s.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(s.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		case time.Time:
			val := T(vi.String())
			switch any(s.addr).(type) {
			case **T:
				*(**T)(unsafe.Pointer(s.addr)) = &val
				return nil
			case *T:
				*(*T)(unsafe.Pointer(s.addr)) = val
				return nil
			default:
				panic("unreachable")
			}
		default:
			return fmt.Errorf(`sequel/encoding: unable to scan %T to ~string`, vi)
		}
	}
}
