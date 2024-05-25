package sequel

import "database/sql/driver"

// For rename table name
type Table struct{}

type DatabaseNamer interface {
	DatabaseName() string
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

type Keyer interface {
	PK() (columnName string, pos int, value driver.Value)
}

type AutoIncrKeyer interface {
	Keyer
	IsAutoIncr()
}

type DuplicateKeyer interface {
	// Allow to support composite key, this only applicable in postgres
	OnDuplicateKey() []string
}
