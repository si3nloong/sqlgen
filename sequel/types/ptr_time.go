package types

import (
	"regexp"
	"time"
	"unsafe"
)

var (
	ddmmyyyy         = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}$`)
	ddmmyyyyhhmmss   = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}$`)
	ddmmyyyyhhmmsstz = regexp.MustCompile(`^\d{4}\-\d{2}\-\d{2}\s\d{2}\:\d{2}:\d{2}\.\d+$`)
)

type ptrOfTime[T time.Time] struct {
	addr **T
}

func PtrOfTime[T time.Time](v **T) ptrOfTime[T] {
	return ptrOfTime[T]{addr: v}
}

func (p ptrOfTime[T]) Interface() *T {
	if p.addr == nil {
		return nil
	}
	return *p.addr
}

func (p ptrOfTime[T]) Scan(v any) error {
	if v == nil {
		(*p.addr) = nil
		return nil
	}

	switch vi := v.(type) {
	case []byte:
		t, err := parseTime(unsafe.String(unsafe.SliceData(vi), len(vi)))
		if err != nil {
			return err
		}
		val := T(t)
		*p.addr = &val
	case string:
		t, err := parseTime(vi)
		if err != nil {
			return err
		}
		val := T(t)
		*p.addr = &val
	case time.Time:
		val := T(vi)
		*p.addr = &val
	}
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
