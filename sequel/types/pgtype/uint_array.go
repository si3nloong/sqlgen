package pgtype

import (
	"database/sql/driver"

	"golang.org/x/exp/constraints"
)

// Uint64Array represents a one-dimensional array of the PostgreSQL integer types.
type (
	UintArray[T constraints.Unsigned] []T
	Uint8Array[T ~uint8]              []T
	Uint16Array[T ~uint16]            []T
	Uint32Array[T ~uint32]            []T
	Uint64Array[T ~uint64]            []T
)

// Value implements the driver.Valuer interface.
func (a UintArray[T]) Value() (driver.Value, error) {
	return uintArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *UintArray[T]) Scan(src any) error {
	return arrayScan(a, src, "UintArray")
}

// Value implements the driver.Valuer interface.
func (a Uint8Array[T]) Value() (driver.Value, error) {
	return uintArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint8Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Uint8Array")
}

// Value implements the driver.Valuer interface.
func (a Uint16Array[T]) Value() (driver.Value, error) {
	return uintArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint16Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Uint16Array")
}

// Value implements the driver.Valuer interface.
func (a Uint32Array[T]) Value() (driver.Value, error) {
	return uintArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint32Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Uint32Array")
}

// Value implements the driver.Valuer interface.
func (a Uint64Array[T]) Value() (driver.Value, error) {
	return uintArrayValue(a)
}

// Scan implements the sql.Scanner interface.
func (a *Uint64Array[T]) Scan(src any) error {
	return arrayScan(a, src, "Uint64Array")
}
