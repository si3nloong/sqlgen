package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"time"
	"unsafe"
)

var (
	ddmmyyyy         = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}$`)
	ddmmyyyyhhmmss   = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}$`)
	ddmmyyyyhhmmsstz = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}\.\d+$`)
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
		return fmt.Errorf(`sequel/types: unsupported scan type %T for time.Time`, vi)
	}
	*s.addr = val
	return nil
}

func parseTime(str string) (time.Time, error) {
	var (
		t   time.Time
		err error
	)
	switch {
	case ddmmyyyy.MatchString(str):
		t, err = time.Parse("2006-01-02", str)
	case ddmmyyyyhhmmss.MatchString(str):
		t, err = time.Parse("2006-01-02 15:04:05", str)
	case ddmmyyyyhhmmsstz.MatchString(str):
		t, err = time.Parse("2006-01-02 15:04:05.999999", str)
	default:
		t, err = time.Parse(time.RFC3339Nano, str)
	}
	return t, err
}
