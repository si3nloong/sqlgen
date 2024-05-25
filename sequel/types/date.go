package types

import (
	"database/sql"
	"database/sql/driver"
	"time"
	"unsafe"

	"cloud.google.com/go/civil"
)

type localDate struct {
	addr *civil.Date
}

var (
	_ sql.Scanner   = (*localDate)(nil)
	_ driver.Valuer = (*localDate)(nil)
)

// Date returns a sql.Scanner
func Date(addr *civil.Date) localDate {
	return localDate{addr: addr}
}

func (b localDate) Value() (driver.Value, error) {
	if b.addr == nil {
		return nil, nil
	}
	return b.addr.MarshalText()
}

func (b localDate) Scan(v any) error {
	var val civil.Date
	switch vi := v.(type) {
	case []byte:
		f, err := civil.ParseDate(unsafe.String(unsafe.SliceData(vi), len(vi)))
		if err != nil {
			return err
		}
		val = f
	case string:
		f, err := civil.ParseDate(vi)
		if err != nil {
			return err
		}
		val = f
	case time.Time:
		val = civil.DateOf(vi)
	}
	*b.addr = val
	return nil
}

type ptrOfLocalDate struct {
	addr **civil.Date
}

// Date returns a sql.Scanner
func PtrOfDate(addr **civil.Date) ptrOfLocalDate {
	return ptrOfLocalDate{addr: addr}
}

func (b ptrOfLocalDate) Value() (driver.Value, error) {
	if b.addr == nil {
		return nil, nil
	}
	return (*b.addr).MarshalText()
}

func (b ptrOfLocalDate) Scan(v any) error {
	var val civil.Date
	switch vi := v.(type) {
	case []byte:
		f, err := civil.ParseDate(unsafe.String(unsafe.SliceData(vi), len(vi)))
		if err != nil {
			return err
		}
		val = f
	case string:
		f, err := civil.ParseDate(vi)
		if err != nil {
			return err
		}
		val = f
	case time.Time:
		val = civil.DateOf(vi)
	}
	*b.addr = &val
	return nil
}
