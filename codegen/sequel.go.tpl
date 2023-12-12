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

type Scanner[T any] interface {
	*T
	Addrs() []any
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
	Scanner[T]
}

type Keyer interface {
	PK() (columnName string, pos int, value driver.Value)
}

type AutoIncrKeyer interface {
	Keyer
	IsAutoIncr()
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

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Dialect interface {
	Driver() string
	// argument string to escape SQL injection
	Var(n int) string
	Wrap(v string) string
	// character to escape table, column name
	QuoteChar() rune
}

type DialectVar interface {
	VarChar() string
}

type Migrator interface {
	Tabler
	CreateTableStmt() string
	AlterTableStmt() string
}

type SingleInserter interface {
	InsertOneStmt() string
}

type KeyFinder interface {
	Keyer
	FindByPKStmt() string
}

type KeyUpdater interface {
	Keyer
	UpdateByPKStmt() string
}

type Inserter interface {
	Columner
	InsertVarQuery() string
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

type Stmt interface {
	StmtBuilder
	fmt.Stringer
	Args() []any
	Reset()
}

type (
	WhereClause   func(StmtBuilder)
	SetClause     func(StmtBuilder)
	OrderByClause func(StmtWriter)
)
