package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

type boolList[T ~bool] struct {
	v *[]T
}

var (
	_ driver.Valuer = (*boolList[bool])(nil)
	_ sql.Scanner   = (*boolList[bool])(nil)
)

func BoolSlice[T ~bool](v *[]T) boolList[T] {
	return boolList[T]{v: v}
}

func (s boolList[T]) Value() (driver.Value, error) {
	if s.v == nil || *s.v == nil {
		return nil, nil
	}
	return encoding.MarshalBoolSlice(*s.v), nil
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
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to %v`, vi, reflect.TypeOf(vi))
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
			flag, err := strconv.ParseBool(unsafe.String(unsafe.SliceData(b), len(b)))
			if err != nil {
				return err
			}
			values[i] = T(flag)
		}
		*s.v = values
	case string:
		if vi == nullStr {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to %v`, vi, reflect.TypeOf(vi))
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		var (
			paths  = strings.Split(vi, ",")
			values = make([]T, len(paths))
		)
		for i := range paths {
			flag, err := strconv.ParseBool(strings.TrimSpace(paths[i]))
			if err != nil {
				return err
			}
			values[i] = T(flag)
		}
		*s.v = values
	case nil:
		*s.v = nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for []~bool`, vi)
	}
	return nil
}
