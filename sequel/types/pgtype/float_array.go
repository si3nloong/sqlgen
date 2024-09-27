package pgtype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// Float64Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float64Array[T ~float64] []T

// Scan implements the sql.Scanner interface.
func (a *Float64Array[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pgtype: cannot convert %T to Float64Array", src)
}

func (a *Float64Array[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float64Array[T], len(elems))
		for i, v := range elems {
			f, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				return fmt.Errorf("pgtype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(f)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float64Array[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, (float64)(a[0]), 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, (float64)(a[i]), 'f', -1, 64)
		}

		return string(append(b, '}')), nil
	}
	return "{}", nil
}

// Float32Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float32Array[T ~float32] []T

// Scan implements the sql.Scanner interface.
func (a *Float32Array[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pgtype: cannot convert %T to Float32Array", src)
}

func (a *Float32Array[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float32Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float32Array[T], len(elems))
		for i, v := range elems {
			f, err := strconv.ParseFloat(string(v), 32)
			if err != nil {
				return fmt.Errorf("pgtype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(f)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float32Array[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, (float64)(a[0]), 'f', -1, 32)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, (float64)(a[i]), 'f', -1, 32)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}
