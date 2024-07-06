package types

import (
	"fmt"
	"unicode/utf8"
)

type runeArray[T ~rune] struct {
	v    []T
	size int
}

func FixedSizeRunes[T ~rune](v []T, size int) *runeArray[T] {
	return &runeArray[T]{v: v, size: size}
}

func (s *runeArray[T]) Scan(v any) error {
	switch b := v.(type) {
	case []byte:
		i := int(0)
		for len(b) > 0 {
			if i >= s.size {
				return fmt.Errorf(`types: rune array overflow, should be %d, but it is %d`, s.size, i)
			}
			r, size := utf8.DecodeRune(b)
			s.v[i] = T(r)
			i++
			b = b[size:]
		}
	default:
		return fmt.Errorf(`types: unsupported scan type %T for [%d]~rune`, b, s.size)
	}
	return nil
}
