package sequel

import "database/sql/driver"

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
	return column[T]{
		colName: columnName,
		v:       convert(value),
		convert: convert,
	}
}

type sqlCol[T any] struct {
	column[T]
	sqlValuer QueryFunc
}

func (c sqlCol[T]) SQLValue(placeholder string) string {
	return c.sqlValuer(placeholder)
}

func SQLColumn[T any](columnName string, value T, sqlValue QueryFunc, convert ConvertFunc[T]) SQLColumnValuer[T] {
	c := sqlCol[T]{}
	c.colName = columnName
	c.v = convert(value)
	c.convert = convert
	c.sqlValuer = sqlValue
	return c
}

type orderByCol struct {
	colName string
	asc     bool
}

func (c orderByCol) ColumnName() string {
	return c.colName
}

func (c orderByCol) Asc() bool {
	return c.asc
}

func OrderByColumn(columnName string, asc bool) OrderByClause {
	return orderByCol{
		colName: columnName,
		asc:     asc,
	}
}
