package types

import (
	"fmt"
	"unsafe"
)

type byteArray[T ~byte] struct {
	v    []T
	size int
}

func FixedSizeBytes[T ~byte](v []T, size int) *byteArray[T] {
	return &byteArray[T]{v: v, size: size}
}

func (s *byteArray[T]) Scan(v any) error {
	switch b := v.(type) {
	case nil:
		s.v = make([]T, s.size)
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
		for i := range b {
			s.v[i] = T(b[i])
		}
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for [%d]~byte`, b, s.size)
	}
	return nil
}
