package sqltype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"
)

// Float32Slice represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float32Slice[T ~float32] []T

// Value implements the driver.Valuer interface.
func (a Float32Slice[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendFloat(b, (float64)(a[0]), 'f', -1, 32)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, (float64)(a[i]), 'f', -1, 32)
		}
		b = append(b, ']')
		return unsafe.String(unsafe.SliceData(b), len(b)), nil
	}
	return "[]", nil
}

// Scan implements the sql.Scanner interface.
func (a *Float32Slice[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("sqltype: cannot convert %T to Float32Slice", src)
}

func (a *Float32Slice[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float32Slice")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float32Slice[T], len(elems))
		for i, v := range elems {
			f, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(v), len(v)), 32)
			if err != nil {
				return fmt.Errorf("sqltype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(f)
		}
		*a = b
	}
	return nil
}

// Float64Slice represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float64Slice[T ~float64] []T

// Value implements the driver.Valuer interface.
func (a Float64Slice[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendFloat(b, (float64)(a[0]), 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, (float64)(a[i]), 'f', -1, 64)
		}

		return string(append(b, ']')), nil
	}
	return "[]", nil
}

// Scan implements the sql.Scanner interface.
func (a *Float64Slice[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("sqltype: cannot convert %T to Float64Slice", src)
}

func (a *Float64Slice[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float64Slice")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float64Slice[T], len(elems))
		for i, v := range elems {
			f, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(v), len(v)), 64)
			if err != nil {
				return fmt.Errorf("sqltype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(f)
		}
		*a = b
	}
	return nil
}
