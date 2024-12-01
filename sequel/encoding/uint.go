package encoding

import (
	"database/sql"
	"fmt"
	"math/bits"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func UintScanner[T ~uint, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &unsignedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: bits.UintSize,
	}
}

func Uint8Scanner[T ~uint8, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &unsignedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 8,
	}
}

func Uint16Scanner[T ~uint16, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &unsignedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 16,
	}
}

func Uint32Scanner[T ~uint32, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &unsignedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 32,
	}
}

func Uint64Scanner[T ~uint64, Addr addrOrPtr[T]](addr Addr) sql.Scanner {
	return &unsignedIntScanner[T, Addr]{
		addr:    addr,
		bitSize: 64,
	}
}

type unsignedIntScanner[T constraints.Unsigned, Addr addrOrPtr[T]] struct {
	addr    Addr
	bitSize int
}

func (u *unsignedIntScanner[T, Addr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch any(u.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(u.addr)) = nil
			return nil
		case *T:
			var v T
			*(*T)(unsafe.Pointer(u.addr)) = v
			return nil
		default:
			panic("unreachable")
		}
	case []byte:
		m, err := strconv.ParseUint(unsafe.String(unsafe.SliceData(vi), len(vi)), 10, u.bitSize)
		if err != nil {
			return err
		}
		val := T(m)
		switch any(u.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(u.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(u.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := T(vi)
		switch any(u.addr).(type) {
		case **T:
			*(**T)(unsafe.Pointer(u.addr)) = &val
			return nil
		case *T:
			*(*T)(unsafe.Pointer(u.addr)) = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		return fmt.Errorf(`sequel/types: unable to scan %T to ~uint`, vi)
	}
}
