package encoding

import (
	"database/sql"
	"fmt"
	"strconv"
	"unsafe"
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

type signedIntScanner[T ~int | ~int8 | ~int16 | ~int32 | ~int64, Addr addrOrPtr[T]] struct {
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
		default: // The destination is *T; the constraint permits only *T or **T.
			var v T
			*(*T)(unsafe.Pointer(i.addr)) = v
			return nil
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
		default: // The destination is *T; the constraint permits only *T or **T.
			*(*T)(unsafe.Pointer(i.addr)) = val
			return nil
		}
	case int64:
		val := T(vi)
		switch any(i.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(i.addr)) = &val
			return nil
		default: // The destination is *T; the constraint permits only *T or **T.
			*(*T)(unsafe.Pointer(i.addr)) = val
			return nil
		}
	default:
		return fmt.Errorf(`sequel/encoding: unable to scan %T to ~int`, vi)
	}
}
