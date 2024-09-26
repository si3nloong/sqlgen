package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/si3nloong/sqlgen/sequel/encoding"
	"golang.org/x/exp/constraints"
)

type intList[T constraints.Signed] struct {
	v *[]T
}

func IntSlice[T constraints.Signed](v *[]T) intList[T] {
	return intList[T]{v: v}
}

func (s intList[T]) Value() (driver.Value, error) {
	if s.v == nil || *s.v == nil {
		return nil, nil
	}
	return encoding.MarshalIntSlice(*s.v), nil
}

func (s *intList[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~int`, vi)
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
			i64, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(b), len(b)), 10, 64)
			if err != nil {
				return err
			}
			values[i] = T(i64)
		}
		*s.v = values
	case string:
		if vi == nullStr {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~int`, vi)
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
			i64, err := strconv.ParseInt(b, 10, 64)
			if err != nil {
				return err
			}
			values[i] = T(i64)
		}
		*s.v = values
	case nil:
		*s.v = nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for []~int`, vi)
	}
	return nil
}
