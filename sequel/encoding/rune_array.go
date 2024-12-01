package encoding

import (
	"database/sql"
	"fmt"
	"unicode/utf8"
)

type runeArrScanner[T ~rune] struct {
	v    []T
	size int
}

func RuneArrayScanner[T ~rune](v []T, size int) sql.Scanner {
	return &runeArrScanner[T]{v: v, size: size}
}

func (s *runeArrScanner[T]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		s.v = make([]T, s.size)
		return nil
	case []byte:
		i := int(0)
		for len(vi) > 0 {
			if i >= s.size {
				return fmt.Errorf(`sequel/types: rune array overflow, should be %d, but it is %d`, s.size, i)
			}
			r, size := utf8.DecodeRune(vi)
			s.v[i] = T(r)
			i++
			vi = vi[size:]
		}
	case string:
		if len(vi) > s.size {
			return fmt.Errorf(`sequel/types: rune array overflow, should be %d, but it is %d`, s.size, len(vi))
		}
		for i := range vi {
			s.v[i] = T(vi[i])
		}
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for [%d]~rune`, vi, s.size)
	}
	return nil
}