package types

import (
	"fmt"
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
	case []byte:
		if len(b) > s.size {
			return fmt.Errorf(`types: byte array overflow, should be %d, but it is %d`, s.size, len(b))
		}
		for i := range b {
			s.v[i] = T(b[i])
		}
	default:
		return fmt.Errorf(`types: unsupported scan type %T for [%d]~byte`, b, s.size)
	}
	return nil
}
