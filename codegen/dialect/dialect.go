package dialect

import (
	"io"
	"sync"
)

var (
	dialectMap     = new(sync.Map)
	defaultDialect = "mysql"
)

type Writer interface {
	io.StringWriter
	io.Writer
}

type Dialect interface {
	// SQL driver name
	Driver() string
	// Argument string to escape SQL injection
	QuoteVar(n int) string
	VarRune() rune
	Var() string
	// Character to escape table, column name
	QuoteIdentifier(v string) string
	// Quote rune can be ' or " or `
	QuoteRune() rune

	ColumnDataTypes() map[string]*ColumnType

	CreateTableStmt(Writer, Schema) error
	// AlterTableStmt(Schema) string
}

type Schema interface {
	DBName() string
	TableName() string
	Columns() []string
	Keys() []string
	ColumnGoType(i int) GoColumn
	Indexes() []string
}

type GoColumn interface {
	Name() string
	DataType() string
	GoName() string
	Size() int
	GoPath() string
	GoType() string
	AutoIncr() bool
	// Key is to identify whether column is primary or foreign key
	Key() bool
	// Type() types.Type
	Nullable() bool

	// Implements(*types.Interface) (*types.Func, bool)
}

func RegisterDialect(name string, d Dialect) {
	if d == nil {
		panic("sqlgen: cannot register nil as dialect")
	}
	dialectMap.Store(name, d)
}

func GetDialect(name string) (Dialect, bool) {
	v, ok := dialectMap.Load(name)
	if !ok {
		return nil, ok
	}
	return v.(Dialect), ok
}

func DefaultDialect() Dialect {
	v, _ := dialectMap.Load(defaultDialect)
	return v.(Dialect)
}
