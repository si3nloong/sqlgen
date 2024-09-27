package pgtype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Int64Array represents a one-dimensional array of the PostgreSQL integer types.
type (
	IntArray[T ~int]     []T
	Int8Array[T ~int8]   []T
	Int16Array[T ~int16] []T
	Int32Array[T ~int32] []T
	Int64Array[T ~int64] []T
)

func IntArrayValue[T constraints.Signed](a []T) driver.Valuer {
	return intArray[T]{v: &a}
}

func IntArrayScanner[T constraints.Signed](a *[]T) sql.Scanner {
	return intArray[T]{v: a}
}

type intArray[T constraints.Signed] struct {
	v *[]T
}

// Value implements the driver.Valuer interface.
func (a intArray[T]) Value() (driver.Value, error) {
	return intArrayValue(*a.v)
}

// Scan implements the sql.Scanner interface.
func (a intArray[T]) Scan(src any) error {
	return arrayScan(a.v, src, "IntArray")
}

// Scan implements the sql.Scanner interface.
func (a *IntArray[T]) Scan(src any) error {
	return arrayScan(a, src, "IntArray")
}

// Value implements the driver.Valuer interface.
func (a IntArray[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int8Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int8Array")
}

// Value implements the driver.Valuer interface.
func (a Int8Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int16Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int16Array")
}

// Value implements the driver.Valuer interface.
func (a Int16Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int32Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int32Array")
}

// Value implements the driver.Valuer interface.
func (a Int32Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Int64Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Int64Array")
}

// Value implements the driver.Valuer interface.
func (a Int64Array[T]) Value() (driver.Value, error) {
	return intArrayValue(a)
}

func scanBytes[T constraints.Integer, Arr interface{ ~[]T }](a *Arr, src []byte, t string) error {
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
