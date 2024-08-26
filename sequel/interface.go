package sequel

import (
	"database/sql/driver"
	"fmt"
	"io"
)

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
	FindOneByPKStmt() (string, []any)
}

type KeyUpdater interface {
	Keyer
	UpdateOneByPKStmt() (string, []any)
}

type KeyDeleter interface {
	Keyer
	DeleteOneByPKStmt() (string, []any)
}

type SingleInserter interface {
	InsertOneStmt() (string, []any)
}

type Inserter interface {
	TableColumnValuer
	InsertPlaceholders(row int) string
}

type TableColumnValuer interface {
	Tabler
	Columner
	Valuer
}

type KeyValuer interface {
	Keyer
	Tabler
	Columner
	Valuer
}

type KeyValueScanner[T any] interface {
	KeyValuer
	PtrScanner[T]
}

type StmtWriter interface {
	io.StringWriter
	io.ByteWriter
	io.Writer
}

type Stmt interface {
	StmtBuilder
	fmt.Stringer
	Args() []any
	Reset()
}

type StmtBuilder interface {
	StmtWriter
	Var(v any) string
	Vars(vals []any) string
}

type ColumnValuer[T any] interface {
	ColumnName() string
	Convert(T) driver.Value
	Value() driver.Value
}

type SQLColumnValuer[T any] interface {
	ColumnName() string
	Convert(T) driver.Value
	Value() driver.Value
	SQLValue(placeholder string) string
}
