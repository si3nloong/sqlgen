package pgtype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"
)

// Int64Array represents a one-dimensional array of the PostgreSQL integer types.
type (
	IntArray[T ~int]     []T
	Int8Array[T ~int8]   []T
	Int16Array[T ~int16] []T
	Int32Array[T ~int32] []T
	Int64Array[T ~int64] []T
)

// Value implements the driver.Valuer interface.
func (a IntArray[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *IntArray[T]) Scan(src any) error {
	return arrayScan(a, src, "IntArray")
}

// Value implements the driver.Valuer interface.
func (a Int8Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int8Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int8Array")
}

// Value implements the driver.Valuer interface.
func (a Int16Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int16Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int16Array")
}

// Value implements the driver.Valuer interface.
func (a Int32Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int32Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int32Array")
}

// Value implements the driver.Valuer interface.
func (a Int64Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int64Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int64Array")
}

func scanBytes[T ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr, Arr interface{ ~[]T }](a *Arr, src []byte, t string) error {
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
				return fmt.Errorf("pgtype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(n)
		}
		*a = b
	}
	return nil
}
