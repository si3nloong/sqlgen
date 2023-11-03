package custom

import (
	"database/sql/driver"
	"strconv"
	t "time"
)

type longText string

func (longText) OK() {}

func (l longText) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(l))), nil
}

func (l *longText) Scan(v any) error {
	return nil
}

func (l longText) Value() (driver.Value, error) {
	return string(l), nil
}

type Addresses []Address

var _ driver.Valuer = (*Addresses)(nil)

func (a Addresses) Value() (driver.Value, error) {
	return "", nil
}

type Customer struct {
	ID        int64 `sql:"id"`
	Age       uint8 `sql:"howOld"`
	Name      longText
	Address   Addresses
	Nicknames []longText
	Status    string `sql:"status"`
	JoinAt    t.Time
}

func newCustomer() *Customer {
	c := &Customer{}
	return c
}
