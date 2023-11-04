package types

import (
	"bytes"
	"fmt"
)

type strList[T ~string] struct {
	v *[]T
}

func StringList[T ~string](v *[]T) strList[T] {
	return strList[T]{v: v}
}

func (s strList[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sqlgen: invalid value of %q to unmarshal to []~string`, vi)
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		b := bytes.Split(vi, []byte{','})
		values := make([]T, len(b))
		for i := range b {
			values[i] = T(bytes.Trim(b[i], `"`))
		}
		*s.v = values
	}
	return nil
}
