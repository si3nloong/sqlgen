package sequel

import (
	"time"
)

type BasicTypes interface {
	string | bool | int64 | float64 | []byte | time.Time
}

func Column[T any](columnName string, value T, convert ConvertFunc[T]) ColumnConvertClause[T] {
	return column[T]{
		name:    columnName,
		value:   value,
		convert: convert,
	}
}

func BasicColumn[T BasicTypes](columnName string, value T) ColumnClause[T] {
	return basicColumn[T]{
		name:  columnName,
		value: value,
	}
}

func SQLColumn[T any](columnName string, value T, sqlValue QueryFunc, convert ConvertFunc[T]) SQLColumnClause[T] {
	c := sqlColumn[T]{}
	c.name = columnName
	c.value = value
	c.convert = convert
	c.sqlValuer = sqlValue
	return c
}

func OrderByColumn(columnName string, asc bool) OrderByClause {
	return orderByColumn{
		column: columnName,
		asc:    asc,
	}
}

type basicColumn[T BasicTypes] struct {
	name  string
	value T
}

func (c basicColumn[T]) ColumnName() string {
	return c.name
}

func (c basicColumn[T]) Value() T {
	return c.value
}

type column[T any] struct {
	name    string
	convert ConvertFunc[T]
	value   T
}

func (c column[T]) ColumnName() string {
	return c.name
}

func (c column[T]) Convert(v T) any {
	return c.convert(v)
}

func (c column[T]) Value() T {
	return c.value
}

type sqlColumn[T any] struct {
	column[T]
	sqlValuer QueryFunc
}

func (c sqlColumn[T]) SQLColumn(placeholder string) string {
	return c.sqlValuer(placeholder)
}

type orderByColumn struct {
	column string
	asc    bool
}

func (c orderByColumn) ColumnName() string {
	return c.column
}

func (c orderByColumn) Asc() bool {
	return c.asc
}
