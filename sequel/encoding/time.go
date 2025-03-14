package encoding

import (
	"database/sql"
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

type tsScanner[T addrOrPtr[time.Time]] struct {
	addr T
}

// Time is a function return a SQL valuer nor scanner, it support
// pointer or pointer of value
func TimeScanner[T addrOrPtr[time.Time]](addr T) sql.Scanner {
	return &tsScanner[T]{addr: addr}
}

func (t *tsScanner[T]) Scan(v any) error {
	switch vt := v.(type) {
	case nil:
		switch vi := any(t.addr).(type) {
		case **time.Time:
			*vi = nil
			return nil
		case *time.Time:
			*vi = time.Time{}
			return nil
		default:
			panic("unreachable")
		}
	case string:
		ts, err := parseTime(vt)
		if err != nil {
			return err
		}
		switch vi := any(t.addr).(type) {
		case **time.Time:
			*vi = &ts
			return nil
		case *time.Time:
			*vi = ts
			return nil
		default:
			panic("unreachable")
		}
	case []byte:
		ts, err := parseTime(unsafe.String(unsafe.SliceData(vt), len(vt)))
		if err != nil {
			return err
		}
		switch vi := any(t.addr).(type) {
		case **time.Time:
			*vi = &ts
			return nil
		case *time.Time:
			*vi = ts
			return nil
		default:
			panic("unreachable")
		}
	case time.Time:
		switch vi := any(t.addr).(type) {
		case **time.Time:
			*vi = &vt
			return nil
		case *time.Time:
			*vi = vt
			return nil
		default:
			panic("unreachable")
		}
	case int64:
		val := time.Unix(vt, 0)
		switch vi := any(t.addr).(type) {
		case **time.Time:
			*vi = &val
			return nil
		case *time.Time:
			*vi = val
			return nil
		default:
			panic("unreachable")
		}
	default:
		return fmt.Errorf(`sequel/encoding: unsupported scan type %T for time.Time`, vt)
	}
}

func parseTime(str string) (time.Time, error) {
	switch {
	case ddmmyyyy.MatchString(str):
		return time.Parse("2006-01-02", str)
	case ddmmyyyyhhmmss.MatchString(str):
		return time.Parse("2006-01-02 15:04:05", str)
	case ddmmyyyyhhmmsstz.MatchString(str):
		return time.Parse("2006-01-02 15:04:05.999999", str)
	default:
		return time.Parse(time.RFC3339Nano, str)
	}
}
