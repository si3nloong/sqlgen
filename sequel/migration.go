package sequel

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
	Columns    []string
	Type       string
	Definition string
}
