package pgtype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// BoolArray represents a one-dimensional array of the PostgreSQL boolean type.
type BoolArray[T ~bool] []T

var (
	_ sql.Scanner   = (*BoolArray[bool])(nil)
	_ driver.Valuer = (*BoolArray[bool])(nil)
)

// Value implements the driver.Valuer interface.
func (a BoolArray[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be exactly two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1+2*n)

		for i := 0; i < n; i++ {
			b[2*i] = ','
			if a[i] {
				b[1+2*i] = 't'
			} else {
				b[1+2*i] = 'f'
			}
		}

		b[0] = '{'
		b[2*n] = '}'

		return string(b), nil
	}
	return "{}", nil
}

// Scan implements the sql.Scanner interface.
func (a *BoolArray[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pgtype: cannot convert %T to BoolArray", src)
}

func (a *BoolArray[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BoolArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BoolArray[T], len(elems))
		for i, v := range elems {
			if len(v) != 1 {
				return fmt.Errorf("pgtype: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
			switch v[0] {
			case 't':
				b[i] = true
			case 'f':
				b[i] = false
			default:
				return fmt.Errorf("pgtype: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
		}
		*a = b
	}
	return nil
}
