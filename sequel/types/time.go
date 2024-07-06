package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
	"unsafe"
)

type datetime[T time.Time] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*datetime[time.Time])(nil)
	_ driver.Valuer = (*datetime[time.Time])(nil)
)

func Time[T time.Time](addr *T, strict ...bool) datetime[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return datetime[T]{addr: addr, strictType: strictType}
}

func (dt datetime[T]) Interface() T {
	if dt.addr == nil {
		return *new(T)
	}
	return *dt.addr
}

func (dt datetime[T]) Value() (driver.Value, error) {
	if dt.addr == nil {
		return nil, nil
	}
	return time.Time(*dt.addr), nil
}

func (s datetime[T]) Scan(v any) error {
	var val T
	switch vi := v.(type) {
	case []byte:
		t, err := parseTime(unsafe.String(unsafe.SliceData(vi), len(vi)))
		if err != nil {
			return err
		}
		val = T(t)
	case string:
		t, err := parseTime(vi)
		if err != nil {
			return err
		}
		val = T(t)
	case time.Time:
		val = T(vi)
	case int64:
		val = T(time.Unix(vi, 0))
	default:
		return fmt.Errorf(`types: unsupported scan type %T for time.Time`, vi)
	}
	*s.addr = val
	return nil
}
