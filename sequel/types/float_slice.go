package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/si3nloong/sqlgen/sequel/encoding"
	"golang.org/x/exp/constraints"
)

type floatList[T constraints.Float] struct {
	v    *[]T
	prec int
}

var (
	_ driver.Valuer = (*floatList[float32])(nil)
	_ sql.Scanner   = (*floatList[float32])(nil)
	_ driver.Valuer = (*floatList[float64])(nil)
	_ sql.Scanner   = (*floatList[float64])(nil)
)

func Float32Slice[T constraints.Float](v *[]T) floatList[T] {
	return floatList[T]{v: v, prec: 32}
}

func Float64Slice[T constraints.Float](v *[]T) floatList[T] {
	return floatList[T]{v: v, prec: 64}
}

func (s floatList[T]) Value() (driver.Value, error) {
	if s.v == nil || *s.v == nil {
		return nil, nil
	}
	return encoding.MarshalFloatList(*s.v, 64), nil
}

func (s *floatList[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~float`, vi)
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
			f, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(b), len(b)), s.prec)
			if err != nil {
				return err
			}
			values[i] = T(f)
		}
		*s.v = values
	case string:
		if vi == nullStr {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~float`, vi)
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		var (
			paths  = strings.Split(vi, ",")
			values = make([]T, len(paths))
			b      string
		)
		for i := range paths {
			b = strings.TrimSpace(paths[i])
			f, err := strconv.ParseFloat(b, s.prec)
			if err != nil {
				return err
			}
			values[i] = T(f)
		}
		*s.v = values
	case nil:
		*s.v = nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for []~float`, vi)
	}
	return nil
}
