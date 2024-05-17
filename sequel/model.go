package sequel

import "database/sql/driver"

// For rename table name
type Table struct{}

type Keyer interface {
	PK() (columnName string, pos int, value driver.Value)
}

type AutoIncrKeyer interface {
	Keyer
	IsAutoIncr()
}

type DuplicateKeyer interface {
	OnDuplicateKey() string
}

type Tabler interface {
	TableName() string
}

type Columner interface {
	Columns() []string
}

type Valuer interface {
	Values() []any
}
