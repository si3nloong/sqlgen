package sequel

import (
	"context"
	"database/sql"
	"io"
)

// For rename table name
type TableName struct{}

type (
	ConvertFunc[T any] func(T) any
	QueryFunc          func(placeholder string) string
	WhereClause        func(StmtWriter)
	SetClause          func(StmtWriter)
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

type InsertColumner interface {
	InsertColumns() []string
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

type AutoIncrKeyer interface {
	PrimaryKeyer
	IsAutoIncr()
}

type PrimaryKeyer interface {
	Keyer
	PK() (string, int, any)
}

type CompositeKeyer interface {
	Keyer
	CompositeKey() ([]string, []int, []any)
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
	Tabler
	Columner
	Valuer
	InsertPlaceholders(row int) string
}

type ColumnValuer interface {
	Tabler
	Columner
	Valuer
}

type KeyScanner interface {
	Keyer
	Tabler
	Columner
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

type KeyPtrScanner[T any] interface {
	KeyScanner
	PtrScanner[T]
}

type RowLevelLocker interface {
	LockMode() string
}

type StmtWriter interface {
	io.Writer
	io.StringWriter
	Quote(v string) string
	Var(v any) string
	// Vars will group the valus in parenthesis
	Vars(vals []any) string
}

type Stmt interface {
	StmtWriter
	Query() string
	Args() []any
	Reset()
}

type ColumnClause[T any] interface {
	ColumnName() string
	Value() T
}

type ColumnConvertClause[T any] interface {
	ColumnClause[T]
	Convert(T) any
}

type SQLColumnClause[T any] interface {
	ColumnConvertClause[T]
	SQLColumn(placeholder string) string
}

type OrderByClause interface {
	ColumnName() string
	Asc() bool
	Desc() bool
}
