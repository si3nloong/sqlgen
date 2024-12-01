package encoding

import (
	"database/sql"
	"fmt"
	"log"
	"unsafe"
)

type byteArrScanner[T ~byte] struct {
	v    []T
	size int
}

func ByteArrayScanner[T ~byte](v []T, size int) sql.Scanner {
	return &byteArrScanner[T]{v: v, size: size}
}

func (s *byteArrScanner[T]) Scan(v any) error {
	switch b := v.(type) {
	case nil:
		var v T
		for i := range s.size {
			s.v[i] = v
		}
		return nil
	case string:
		bytes := unsafe.Slice(unsafe.StringData(b), len(b))
		if len(bytes) > s.size {
			return fmt.Errorf(`sequel/types: byte array overflow, should be %d, but it is %d`, s.size, len(b))
		}
		for i := range bytes {
			s.v[i] = T(bytes[i])
		}
	case []byte:
		if len(b) > s.size {
			return fmt.Errorf(`sequel/types: byte array overflow, should be %d, but it is %d`, s.size, len(b))
		}
		log.Println("HERE")
		for i := range b {
			s.v[i] = T(b[i])
		}
		return nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for [%d]~byte`, b, s.size)
	}
	return nil
}
