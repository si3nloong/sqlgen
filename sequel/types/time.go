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

type timestamp[T time.Time] struct {
	addr       *T
	strictType bool
}

var (
	_ sql.Scanner   = (*timestamp[time.Time])(nil)
	_ driver.Valuer = (*timestamp[time.Time])(nil)
)

func Time[T time.Time](addr *T, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &timestamp[T]{addr: addr, strictType: strictType}
}

func (t timestamp[T]) Interface() T {
	if t.addr == nil {
		return *new(T)
	}
	return *t.addr
}

func (t timestamp[T]) Value() (driver.Value, error) {
	if t.addr == nil {
		return nil, nil
	}
	return time.Time(*t.addr), nil
}

func (t *timestamp[T]) Scan(v any) error {
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
	case nil:
		t.addr = nil
		return nil
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for time.Time`, vi)
	}
	*t.addr = val
	return nil
}

func parseTime(str string) (t time.Time, err error) {
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
