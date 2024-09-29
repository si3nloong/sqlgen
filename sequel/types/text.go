package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"unsafe"
)

type textMarshaler[T interface {
	encoding.TextMarshaler
}] struct {
	v T
}

func TextMarshaler[T interface {
	encoding.TextMarshaler
}](addr T) driver.Valuer {
	return textMarshaler[T]{v: addr}
}

func (b textMarshaler[T]) Value() (driver.Value, error) {
	return b.v.MarshalText()
}

type textUnmarshaler[T any, Ptr interface {
	*T
	encoding.TextUnmarshaler
}] struct {
	v Ptr
}

func TextUnmarshaler[T any, Ptr interface {
	*T
	encoding.TextUnmarshaler
}](addr Ptr) sql.Scanner {
	return textUnmarshaler[T, Ptr]{v: addr}
}

func (b textUnmarshaler[T, Ptr]) Scan(v any) error {
	switch vi := v.(type) {
	case string:
		return b.v.UnmarshalText(unsafe.Slice(unsafe.StringData(vi), len(vi)))
	case []byte:
		return b.v.UnmarshalText(vi)
	default:
		return fmt.Errorf(`sequel/types: unable to unmarshal %T`, vi)
	}
}
