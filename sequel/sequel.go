package sequel

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
)

type Name struct{ Local string }

type Scanner[T any] interface {
	*T
	Addrs() []any
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
	IsAutoIncr() bool
	PK() (columnName string, pos int, value driver.Value)
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
	Var(n int) string
	Wrap(v string) string
	Driver() string
}

type Migrator interface {
	Tabler
	CreateTableStmt() string
	AlterTableStmt() string
}

type SingleInserter interface {
	InsertOneStmt() string
}

type Inserter interface {
	Columner
	InsertVarStmt() string
}

type StmtWriter interface {
	io.StringWriter
	io.ByteWriter
}

type StmtBuilder interface {
	StmtWriter
	Var(query string, v ...any)
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
