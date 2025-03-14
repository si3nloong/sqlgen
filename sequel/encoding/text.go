package encoding

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

type textValuer[T encoding.TextMarshaler] struct {
	v T
}

func TextValue[T encoding.TextMarshaler](addr T) driver.Valuer {
	return textValuer[T]{v: addr}
}

func (b textValuer[T]) Value() (driver.Value, error) {
	return b.v.MarshalText()
}

type textScanner[T any, Ptr interface {
	*T
	encoding.TextUnmarshaler
}, Addr interface{ *T | **T }] struct {
	v Addr
}

func TextScanner[T any, Ptr interface {
	*T
	encoding.TextUnmarshaler
}, Addr interface{ *T | **T }](addr Addr) sql.Scanner {
	return &textScanner[T, Ptr, Addr]{v: addr}
}

func (t *textScanner[T, Ptr, Addr]) Scan(v any) error {
	switch vt := v.(type) {
	case nil:
		switch any(t.v).(type) {
		case **T:
			*(**T)(unsafe.Pointer(t.v)) = nil
			return nil
		case *T:
			var val T
			*(*T)(unsafe.Pointer(t.v)) = val
			return nil
		default:
			panic("unreachable")
		}
	case string:
		switch vi := any(t.v).(type) {
		case *Ptr:
			return (*vi).UnmarshalText(unsafe.Slice(unsafe.StringData(vt), len(vt)))
		case Ptr:
			return vi.UnmarshalText(unsafe.Slice(unsafe.StringData(vt), len(vt)))
		default:
			panic("unreachable")
		}
	case []byte:
		switch vi := any(t.v).(type) {
		case *Ptr:
			return (*vi).UnmarshalText(vt)
		case Ptr:
			return vi.UnmarshalText(vt)
		default:
			panic("unreachable")
		}
	case bool:
		b := []byte{0}
		if vt {
			b = []byte{1}
		}
		switch vi := any(t.v).(type) {
		case *Ptr:
			return (*vi).UnmarshalText(b)
		case Ptr:
			return vi.UnmarshalText(b)
		default:
			panic("unreachable")
		}
	case int64:
		val := strconv.FormatInt(vt, 10)
		switch vi := any(t.v).(type) {
		case *Ptr:
			return (*vi).UnmarshalText(unsafe.Slice(unsafe.StringData(val), len(val)))
		case Ptr:
			return vi.UnmarshalText(unsafe.Slice(unsafe.StringData(val), len(val)))
		default:
			panic("unreachable")
		}
	case float64:
		val := strconv.FormatFloat(vt, 'f', -1, 64)
		switch vi := any(t.v).(type) {
		case *Ptr:
			return (*vi).UnmarshalText(unsafe.Slice(unsafe.StringData(val), len(val)))
		case Ptr:
			return vi.UnmarshalText(unsafe.Slice(unsafe.StringData(val), len(val)))

		default:
			panic("unreachable")
		}
	case time.Time:
		b, err := vt.MarshalText()
		if err != nil {
			return err
		}
		switch vi := any(t.v).(type) {
		case *Ptr:
			return (*vi).UnmarshalText(b)
		case Ptr:
			return vi.UnmarshalText(b)
		default:
			panic("unreachable")
		}
	default:
		return fmt.Errorf(`sequel/encoding: unable to TextUnmarshal %T`, vt)
	}
}
