package sequel

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
)

// For rename table name
type Table struct{}

type (
	ConvertFunc[T any] func(T) driver.Value
	QueryFunc          func(placeholder string) string
	WhereClause        func(StmtBuilder)
	SetClause          func(StmtBuilder)
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Databaser interface {
	DatabaseName() string
}

type Tabler interface {
	TableName() string
}

type Columner interface {
	Columns() []string
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
	io.Writer
	io.StringWriter
	io.ByteWriter
}

type StmtBuilder interface {
	StmtWriter
	Var(v any) string
	// Vars will group the valus in parenthesis
	Vars(vals []any) string
}

type Stmt interface {
	StmtBuilder
	fmt.Stringer
	fmt.Formatter
	Args() []any
	Reset()
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

type ColumnOrder interface {
	ColumnName() string
	Asc() bool
}
