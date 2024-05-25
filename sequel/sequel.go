package sequel

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/si3nloong/sqlgen/codegen/templates"
)

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

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Dialect interface {
	// SQL driver name
	Driver() string
	// Argument string to escape SQL injection
	QuoteVar(n int) string
	VarRune() rune
	// Character to escape table, column name
	QuoteIdentifier(v string) string
	QuoteRune() rune

	CreateTableStmt(n string, model *templates.Model) string
	AlterTableStmt(n string, model *templates.Model) string
}

type Migrator interface {
	Tabler
	CreateTableStmt() string
	AlterTableStmt() string
}

type MigratorV2 interface {
	Up(ctx context.Context, db DB) error
	Down(ctx context.Context, db DB) error
}

type SingleInserter interface {
	InsertOneStmt() string
}

type SingleUpserter interface {
	UpsertOneStmt() string
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
