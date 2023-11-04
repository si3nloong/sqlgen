package types

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/si3nloong/sqlgen/internal/strfmt"
)

type boolList[T ~bool] struct {
	v *[]T
}

func BoolList[T ~bool](v *[]T) boolList[T] {
	return boolList[T]{v: v}
}

func (s boolList[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sqlgen: invalid value of %q to unmarshal to %v`, vi, reflect.TypeOf(vi))
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		b := bytes.Split(vi, []byte{','})
		values := make([]T, len(b))
		for i := range b {
			flag, err := strconv.ParseBool(strfmt.B2s(b[i]))
			if err != nil {
				return err
			}
			values[i] = T(flag)
		}
		*s.v = values
	}
	return nil
}
