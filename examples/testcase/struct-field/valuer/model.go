package valuer

import (
	"database/sql/driver"
)

type anyType struct{ ptr bool }

func (a anyType) Value() (driver.Value, error) {
	if a.ptr {
		return "ptr", nil
	}
	return "any", nil
}

type B struct {
	ID       int64
	Value    anyType
	PtrValue *anyType
	N        string
}
