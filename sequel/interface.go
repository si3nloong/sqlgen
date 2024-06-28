package sequel

import "io"

// For rename table name
type Table struct{}

type (
	WhereClause   func(StmtBuilder)
	SetClause     func(StmtBuilder)
	OrderByClause func(StmtWriter)
)

type Databaser interface {
	DatabaseName() string
}

type Tabler interface {
	TableName() string
}

type Columner interface {
	ColumnNames() []string
}

type SQLColumner interface {
	SQLColumns() []string
}

type Valuer interface {
	Values() []any
}

type Scanner interface {
	Addrs() []any
}

type PtrScanner[T any] interface {
	*T
	Scanner
}

type Keyer interface {
	HasPK()
}

type PrimaryKeyer interface {
	Keyer
	PK() (string, int, any)
}

type CompositeKeyer interface {
	Keyer
	CompositeKey() ([]string, []int, []any)
}

type AutoIncrKeyer interface {
	PrimaryKeyer
	IsAutoIncr()
}

type KeyFinder interface {
	Keyer
	FindByPKStmt() (string, []any)
}

type KeyUpdater interface {
	Keyer
	UpdateByPKStmt() (string, []any)
}

type KeyDeleter interface {
	Keyer
	DeleteByPKStmt() (string, []any)
}

type SingleInserter interface {
	InsertOneStmt() (string, []any)
}

type Inserter interface {
	Columner
	InsertPlaceholders(row int) string
}

type TableColumnValuer[T any] interface {
	Tabler
	Columner
	Valuer
}

type KeyValuer[T any] interface {
	Keyer
	Tabler
	Columner
	Valuer
}

type KeyValueScanner[T any] interface {
	KeyValuer[T]
	PtrScanner[T]
}

type StmtWriter interface {
	io.StringWriter
	io.ByteWriter
}

type StmtBuilder interface {
	StmtWriter
	Var(query string, v any)
	Vars(query string, v []any)
}
