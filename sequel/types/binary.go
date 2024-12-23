package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
)

type binaryMarshaler[T interface {
	encoding.BinaryMarshaler
}] struct {
	v T
}

func BinaryMarshaler[T interface {
	encoding.BinaryMarshaler
}](addr T) driver.Valuer {
	return binaryMarshaler[T]{v: addr}
}

func (b binaryMarshaler[T]) Value() (driver.Value, error) {
	return b.v.MarshalBinary()
}

type binaryUnmarshaler[T any, Ptr interface {
	*T
	encoding.BinaryUnmarshaler
}] struct {
	v Ptr
}

func BinaryUnmarshaler[T any, Ptr interface {
	*T
	encoding.BinaryUnmarshaler
}](addr Ptr) sql.Scanner {
	return binaryUnmarshaler[T, Ptr]{v: addr}
}

func (b binaryUnmarshaler[T, Ptr]) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		return b.v.UnmarshalBinary(vi)
	default:
		return fmt.Errorf(`sequel/types: binary must be []byte`)
	}
}
