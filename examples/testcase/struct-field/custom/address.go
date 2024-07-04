package custom

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	o "github.com/paulmach/orb"
)

func empty() {}

type Address struct {
	Line1       string
	Line2       sql.NullString
	City        string
	PostCode    uint
	StateCode   StateCode
	GeoPoint    o.Point
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
