package types

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/si3nloong/sqlgen/internal/strfmt"
	"golang.org/x/exp/constraints"
)

type floatList[T constraints.Float] struct {
	v *[]T
}

func FloatList[T constraints.Float](v *[]T) floatList[T] {
	return floatList[T]{v: v}
}

func (s floatList[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sqlgen: invalid value of %q to unmarshal to []~float`, vi)
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		var (
			paths  = bytes.Split(vi, []byte{','})
			values = make([]T, len(paths))
			b      []byte
		)
		for i := range paths {
			b = bytes.TrimSpace(paths[i])
			f64, err := strconv.ParseFloat(strfmt.B2s(b), 64)
			if err != nil {
				return err
			}
			values[i] = T(f64)
		}
		*s.v = values
	}
	return nil
}
