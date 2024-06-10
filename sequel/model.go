package sequel

// For rename table name
type Table struct{ Name string }

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
	HasPK()
}

type PrimaryKeyer interface {
	Keyer
	PK() (string, int, any)
}

type AutoIncrKeyer interface {
	PrimaryKeyer
	AutoIncr()
}

type CompositeKeyer interface {
	Keyer
	CompositeKey() ([]string, []int, []any)
}

type DuplicateKeyer interface {
	// Allow to support composite key, this only applicable in postgres
	OnDuplicateKey() []string
}
