package sequel

type TableSchema interface {
	PK() (TablePK, bool)
	Columns() []string
	Indexes() []string
	Column(i int) TableColumn
	Index(i int) TableIndex
}

type TablePK interface {
	Columns() []string
	Definition() string
}

type TableColumn interface {
	Name() string
	Length() int64
	DataType() string
	Nullable() bool
	// Default() any
	Comment() string
	Definition() string
}

type TableIndex interface {
	Name() string
	Type() string
	Columns() []string
	Definition() string
}

type Migrator interface {
	Schemas() TableDefinition
}

type TableDefinition struct {
	PK      *PrimaryKeyDefinition
	Columns []ColumnDefinition
	Indexes []IndexDefinition
}

type PrimaryKeyDefinition struct {
	Columns    []string
	Definition string
}

type ColumnDefinition struct {
	Name       string
	Definition string
}

type IndexDefinition struct {
	Name       string
	Definition string
}
