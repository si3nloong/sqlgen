package sqltype

import (
	"cmp"
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"
)

// Int64Slice represents a one-dimensional array of the PostgreSQL integer types.
type (
	IntSlice[T ~int]     []T
	Int8Slice[T ~int8]   []T
	Int16Slice[T ~int16] []T
	Int32Slice[T ~int32] []T
	Int64Slice[T ~int64] []T
)

// Value implements the driver.Valuer interface.
func (a IntSlice[T]) Value() (driver.Value, error) {
	return intSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *IntSlice[T]) Scan(src any) error {
	return arrayScan(a, src, "IntSlice")
}

// Value implements the driver.Valuer interface.
func (a Int8Slice[T]) Value() (driver.Value, error) {
	return intSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int8Slice[T]) Scan(src any) error {
	return arrayScan(a, src, "Int8Slice")
}

// Value implements the driver.Valuer interface.
func (a Int16Slice[T]) Value() (driver.Value, error) {
	return intSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int16Slice[T]) Scan(src any) error {
	return arrayScan(a, src, "Int16Slice")
}

// Value implements the driver.Valuer interface.
func (a Int32Slice[T]) Value() (driver.Value, error) {
	return intSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int32Slice[T]) Scan(src any) error {
	return arrayScan(a, src, "Int32Slice")
}

// Value implements the driver.Valuer interface.
func (a Int64Slice[T]) Value() (driver.Value, error) {
	return intSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int64Slice[T]) Scan(src any) error {
	return arrayScan(a, src, "Int64Slice")
}

func scanBytes[T cmp.Ordered, Arr interface{ ~[]T }](a *Arr, src []byte, t string) error {
	elems, err := scanLinearArray(src, []byte{','}, t)
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make([]T, len(elems))
		for i, v := range elems {
			n, err := strconv.ParseInt(unsafe.String(unsafe.SliceData(v), len(v)), 10, 64)
			if err != nil {
				return fmt.Errorf("sqltype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(n)
		}
		*a = b
	}
	return nil
}

func arrayScan[T cmp.Ordered, Arr interface{ ~[]T }](a *Arr, src any, t string) error {
	switch src := src.(type) {
	case []byte:
		return scanBytes(a, src, t)
	case string:
		return scanBytes(a, unsafe.Slice(unsafe.StringData(src), len(src)), t)
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("sqltype: cannot convert %T to %s", src, t)
}

func intSliceValue[T ~int | ~int8 | ~int16 | ~int32 | ~int64](a []T) (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendInt(b, (int64)(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, (int64)(a[i]), 10)
		}

		return string(append(b, ']')), nil
	}
	return "[]", nil
}
