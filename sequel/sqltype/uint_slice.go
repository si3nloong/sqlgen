package sqltype

import (
	"database/sql/driver"
	"fmt"
	"math/bits"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Uint64Slice represents a one-dimensional array of the PostgreSQL integer types.
type (
	UintSlice[T ~uint]     []T
	Uint8Slice[T ~uint8]   []T
	Uint16Slice[T ~uint16] []T
	Uint32Slice[T ~uint32] []T
	Uint64Slice[T ~uint64] []T
)

// Value implements the driver.Valuer interface.
func (a UintSlice[T]) Value() (driver.Value, error) {
	return uintSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *UintSlice[T]) Scan(src any) error {
	return arrayUScan(bits.UintSize, a, src, "UintSlice")
}

// Value implements the driver.Valuer interface.
func (a Uint8Slice[T]) Value() (driver.Value, error) {
	return uintSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint8Slice[T]) Scan(src any) error {
	return arrayUScan(8, a, src, "Uint8Slice")
}

// Value implements the driver.Valuer interface.
func (a Uint16Slice[T]) Value() (driver.Value, error) {
	return uintSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint16Slice[T]) Scan(src any) error {
	return arrayUScan(16, a, src, "Uint16Slice")
}

// Value implements the driver.Valuer interface.
func (a Uint32Slice[T]) Value() (driver.Value, error) {
	return uintSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint32Slice[T]) Scan(src any) error {
	return arrayUScan(32, a, src, "Uint32Slice")
}

// Value implements the driver.Valuer interface.
func (a Uint64Slice[T]) Value() (driver.Value, error) {
	return uintSliceValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint64Slice[T]) Scan(src any) error {
	return arrayUScan(64, a, src, "Uint64Slice")
}

func scanUBytes[T constraints.Unsigned, Arr interface{ ~[]T }](bitSize int, a *Arr, src []byte, t string) error {
	elems, err := scanLinearArray(src, []byte{','}, t)
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make([]T, len(elems))
		for i, v := range elems {
			n, err := strconv.ParseUint(unsafe.String(unsafe.SliceData(v), len(v)), 10, bitSize)
			if err != nil {
				return fmt.Errorf("pgtype: parsing array element index %d: %v", i, err)
			}
			b[i] = T(n)
		}
		*a = b
	}
	return nil
}

func arrayUScan[T constraints.Unsigned, Arr interface{ ~[]T }](bitSize int, a *Arr, src any, t string) error {
	switch src := src.(type) {
	case []byte:
		return scanUBytes(bitSize, a, src, t)
	case string:
		return scanUBytes(bitSize, a, unsafe.Slice(unsafe.StringData(src), len(src)), t)
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pgtype: cannot convert %T to %s", src, t)
}

func uintSliceValue[T constraints.Unsigned](a []T) (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendUint(b, (uint64)(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendUint(b, (uint64)(a[i]), 10)
		}

		return string(append(b, ']')), nil
	}
	return "[]", nil
}
