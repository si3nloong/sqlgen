package dialect

type Kind int

const (
	Int8 Kind = iota
	Int16
	Int32
	Int64
	Int
	Uint8
	Uint16
	Uint32
	Uint64
	Uint
)

type ColumnType struct {
	DataType   func(col GoColumn) string
	Scanner    string
	Valuer     string
	SQLScanner string
	SQLValuer  string
}
