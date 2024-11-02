package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type uintList[T constraints.Unsigned] struct {
	v *[]T
}

func UintSlice[T constraints.Unsigned](v *[]T) uintList[T] {
	return uintList[T]{v: v}
}

func (s *uintList[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~uint`, vi)
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
			u64, err := strconv.ParseUint(unsafe.String(unsafe.SliceData(b), len(b)), 10, 64)
			if err != nil {
				return err
			}
			values[i] = T(u64)
		}
		*s.v = values
	case string:
		if vi == nullStr {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~uint`, vi)
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
			u64, err := strconv.ParseUint(b, 10, 64)
			if err != nil {
				return err
			}
			values[i] = T(u64)
		}
		*s.v = values
	case nil:
		*s.v = nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for []~uint`, vi)
	}
	return nil
}
