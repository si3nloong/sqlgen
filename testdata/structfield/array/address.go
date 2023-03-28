package array

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

func empty() {}

type Address struct {
	Line1       string
	Line2       sql.NullString
	City        string
	PostCode    uint
	StateCode   StateCode
	CountryCode CountryCode
}

var (
	_ sql.Scanner = (*Address)(nil)
)

func (Address) ID() {}

func (a *Address) Scan(v interface{}) error {
	return nil
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}
