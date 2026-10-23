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
		clear(s.v[:s.size])
		return nil
	case []byte:
		i := int(0)
		for len(vi) > 0 {
			if i >= s.size {
				return fmt.Errorf(`sequel/encoding: rune array overflow, should be %d, but it is %d`, s.size, i)
			}
			r, size := utf8.DecodeRune(vi)
			s.v[i] = T(r)
			i++
			vi = vi[size:]
		}
	case string:
		if n := utf8.RuneCountInString(vi); n > s.size {
			return fmt.Errorf(`sequel/encoding: rune array overflow, should be %d, but it is %d`, s.size, n)
		}
		i := 0
		for _, r := range vi {
			s.v[i] = T(r)
			i++
		}
	default:
		return fmt.Errorf(`sequel/encoding: unsupported scan type %T for [%d]~rune`, vi, s.size)
	}
	return nil
}
