package encoding

import (
	"database/sql"
	"fmt"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func IntScanner[T ~int, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &signedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: strconv.IntSize,
	}
}

func Int8Scanner[T ~int8, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &signedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 8,
	}
}

func Int16Scanner[T ~int16, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &signedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 16,
	}
}

func Int32Scanner[T ~int32, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &signedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 32,
	}
}

func Int64Scanner[T ~int64, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &signedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 64,
	}
}

type signedIntScanner[T constraints.Signed, Addr addrOrPtr[T]] struct {
	addr    Addr
	bitSize int
}

func (i *signedIntScanner[T, Addr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch any(i.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(i.addr)) = nil
			return nil
		case *T:
			var v T
			*(*T)(unsafe.Pointer(i.addr)) = v
			return nil
		default:
			panic("unreachable")
		}
	case []byte:
		m, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(vi), len(vi)), 10, i.bitSize)
		if err != nil {
			return err
		}
		val := T(m)
		switch any(i.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(i.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(i.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := T(vi)
		switch any(i.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(i.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(i.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		return fmt.Errorf(`sequel/encoding: unable to scan %T to ~int`, vi)
	}
}
