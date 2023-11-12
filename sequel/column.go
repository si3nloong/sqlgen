package sequel

import "database/sql/driver"

type ConvertFunc[T any] func(T) driver.Value

type ColumnValuer[T any] interface {
	ColumnName() string
	Convert(T) driver.Value
	Value() driver.Value
}

type colValuer[T any] struct {
	colName string
	convert ConvertFunc[T]
	v       driver.Value
}

func (c colValuer[T]) ColumnName() string {
	return c.colName
}

func (c colValuer[T]) Convert(v T) driver.Value {
	return c.convert(v)
}

func (c colValuer[T]) Value() driver.Value {
	return c.v
}

func Column[T any](columnName string, value T, convert ConvertFunc[T]) ColumnValuer[T] {
	return colValuer[T]{colName: columnName, v: convert(value), convert: convert}
}
