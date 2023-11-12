package types

import (
	"database/sql/driver"
	"encoding"
	"fmt"
)

type textMarshaler[T interface {
	encoding.TextMarshaler
}] struct {
	v T
}

func TextMarshaler[T interface {
	encoding.TextMarshaler
}](addr T) textMarshaler[T] {
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
}](addr Ptr) textUnmarshaler[T, Ptr] {
	return textUnmarshaler[T, Ptr]{v: addr}
}

func (b textUnmarshaler[T, Ptr]) Scan(v any) error {
	switch vi := v.(type) {
	case string:
		return b.v.UnmarshalText([]byte(vi))
	case []byte:
		return b.v.UnmarshalText(vi)
	default:
		return fmt.Errorf(`sqlgen: text must be []byte to unmarshal`)
	}
}
