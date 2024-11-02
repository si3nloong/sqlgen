package pgtype

import (
	"database/sql/driver"
	"fmt"
	"unsafe"
)

// StringArray represents a one-dimensional array of the PostgreSQL character types.
type StringArray[T ~string] []T

// Value implements the driver.Valuer interface.
func (a StringArray[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'

		b = appendArrayQuotedBytes(b, []byte(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i]))
		}
		b = append(b, '}')
		return unsafe.String(unsafe.SliceData(b), len(b)), nil
	}
	return "{}", nil
}

// Scan implements the sql.Scanner interface.
func (a *StringArray[T]) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pgtype: cannot convert %T to StringArray", src)
}

func (a *StringArray[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "StringArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(StringArray[T], len(elems))
		for i, v := range elems {
			if v == nil {
				return fmt.Errorf("pgtype: parsing array element index %d: cannot convert nil to string", i)
			}
			b[i] = T(v)
		}
		*a = b
	}
	return nil
}
