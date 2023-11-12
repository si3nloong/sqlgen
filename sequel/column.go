package sequel

import "database/sql/driver"

type ConvertFunc[T any] func(T) driver.Value

type ColumnValuer[T any] interface {
	ColumnName() string
	Convert(T) driver.Value
	Value() driver.Value
}

type column[T any] struct {
	colName string
	convert ConvertFunc[T]
	v       driver.Value
}

func (c column[T]) ColumnName() string {
	return c.colName
}

func (c column[T]) Convert(v T) driver.Value {
	return c.convert(v)
}

func (c column[T]) Value() driver.Value {
	return c.v
}

func Column[T any](columnName string, value T, convert ConvertFunc[T]) ColumnValuer[T] {
	return column[T]{colName: columnName, v: convert(value), convert: convert}
}
