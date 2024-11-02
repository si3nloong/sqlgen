package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

type strSlice[T ~string] struct {
	v *[]T
}

var (
	_ driver.Valuer = (*strSlice[string])(nil)
	_ sql.Scanner   = (*strSlice[string])(nil)
)

func StringSlice[T ~string](v *[]T) strSlice[T] {
	return strSlice[T]{v: v}
}

func (s strSlice[T]) Value() (driver.Value, error) {
	if s.v == nil || *s.v == nil {
		return nil, nil
	}
	return encoding.MarshalStringSlice(*s.v), nil
}

func (s *strSlice[T]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if bytes.Equal(vi, nullBytes) {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~string`, vi)
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		b := bytes.Split(vi, []byte{','})
		values := make([]T, len(b))
		for i := range b {
			values[i] = (T)(bytes.Trim(b[i], `"`))
		}
		*s.v = values
	case string:
		if vi == nullStr {
			*s.v = nil
			return nil
		}
		length := len(vi)
		if length < 2 || vi[0] != '[' || vi[length-1] != ']' {
			return fmt.Errorf(`sequel/types: invalid value of %q to unmarshal to []~string`, vi)
		}
		vi = vi[1 : length-1]
		if len(vi) == 0 {
			return nil
		}
		b := strings.Split(vi, ",")
		values := make([]T, len(b))
		for i := range b {
			values[i] = (T)(strings.Trim(b[i], `"`))
		}
		*s.v = values
	case nil:
		*s.v = nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for []~string`, vi)
	}
	return nil
}
