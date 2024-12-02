package types

import (
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

type timestamp[T time.Time, Ptr *T | **T] struct {
	addr       Ptr
	strictType bool
}

func Time[T time.Time, Ptr *T | **T](addr Ptr, strict ...bool) ValueScanner[T] {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return &timestamp[T, Ptr]{addr: addr, strictType: strictType}
}

func (t timestamp[T, Ptr]) Interface() T {
	switch vi := any(t.addr).(type) {
	case **T:
		if vi == nil || *vi == nil {
			return *new(T)
		}
		return **vi
	case *T:
		if vi == nil {
			return *new(T)
		}
		return *vi
	default:
		panic("unreachable")
	}
}

func (t timestamp[T, Ptr]) Value() (driver.Value, error) {
	switch vi := any(t.addr).(type) {
	case **T:
		if vi == nil || *vi == nil {
			return nil, nil
		}
		return (**vi), nil
	case *T:
		if vi == nil {
			return nil, nil
		}
		return time.Time(*vi), nil
	default:
		panic("unreachable")
	}
}

func (t *timestamp[T, Ptr]) Scan(v any) error {
	switch vt := v.(type) {
	case nil:
		switch vi := any(t.addr).(type) {
		case **T:
			*vi = nil
			return nil
		case *T:
			*vi = T(time.Time{})
			return nil
		default:
			panic("unreachable")
		}
	case []byte:
		ts, err := parseTime(unsafe.String(unsafe.SliceData(vt), len(vt)))
		if err != nil {
			return err
		}
		val := T(ts)
		switch vi := any(t.addr).(type) {
		case **T:
			*vi = &val
			return nil
		case *T:
			*vi = val
			return nil
		default:
			panic("unreachable")
		}
	case string:
		ts, err := parseTime(vt)
		if err != nil {
			return err
		}
		val := T(ts)
		switch vi := any(t.addr).(type) {
		case **T:
			*vi = &val
			return nil
		case *T:
			*vi = val
			return nil
		default:
			panic("unreachable")
		}
	case time.Time:
		val := T(vt)
		switch vi := any(t.addr).(type) {
		case **T:
			*vi = &val
			return nil
		case *T:
			*vi = val
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := T(time.Unix(vt, 0))
		switch vi := any(t.addr).(type) {
		case **T:
			*vi = &val
			return nil
		case *T:
			*vi = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		return fmt.Errorf(`sequel/types: unsupported scan type %T for time.Time`, vt)
	}
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
