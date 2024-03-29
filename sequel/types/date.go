package types

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/internal/strfmt"
)

type localDate struct {
	addr       *civil.Date
	strictType bool
}

var (
	_ sql.Scanner   = (*localDate)(nil)
	_ driver.Valuer = (*localDate)(nil)
)

// Date returns a sql.Scanner
func Date(addr *civil.Date, strict ...bool) localDate {
	var strictType bool
	if len(strict) > 0 {
		strictType = strict[0]
	}
	return localDate{addr: addr, strictType: strictType}
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
		f, err := civil.ParseDate(strfmt.B2s(vi))
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
		f, err := civil.ParseDate(strfmt.B2s(vi))
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
