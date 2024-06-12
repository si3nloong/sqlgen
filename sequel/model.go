package sequel

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
	HasPK()
}

type PrimaryKeyer interface {
	Keyer
	PK() (string, int, any)
}

type AutoIncrKeyer interface {
	PrimaryKeyer
	IsAutoIncr()
}

type CompositeKeyer interface {
	Keyer
	CompositeKey() ([]string, []int, []any)
}
