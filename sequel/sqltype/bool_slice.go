package sqltype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"unsafe"
)

// BoolSlice represents a one-dimensional array of the PostgreSQL boolean type.
type BoolSlice[T ~bool] []T

var (
	_ sql.Scanner   = (*BoolSlice[bool])(nil)
	_ driver.Valuer = (*BoolSlice[bool])(nil)
)

// Value implements the driver.Valuer interface.
func (a BoolSlice[T]) Value() (driver.Value, error) {
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
		return unsafe.String(unsafe.SliceData(b), len(b)), nil
	}
	return "{}", nil
}

// Scan implements the sql.Scanner interface.
func (a *BoolSlice[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(unsafe.Slice(unsafe.StringData(src), len(src)))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pgtype: cannot convert %T to BoolSlice", src)
}

func (a *BoolSlice[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BoolSlice")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BoolSlice[T], len(elems))
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