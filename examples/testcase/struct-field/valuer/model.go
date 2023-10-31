package valuer

import "database/sql/driver"

type anyType struct{}

func (anyType) Value() (driver.Value, error) {
	return nil, nil
}

type B struct {
	ID    int64
	Value anyType
	N     string
}
