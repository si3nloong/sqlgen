package pk

import "database/sql/driver"

type LongText string

type PK int64

var _ driver.Valuer = (*PK)(nil)

func (pk PK) Value() (driver.Value, error) {
	return int64(pk), nil
}
