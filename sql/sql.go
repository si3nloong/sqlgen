package sql

import (
	"context"
	"database/sql"
)

type Scanner[T any] interface {
	*T
	Addrs() []any
}

type KeyValuer[T any] interface {
	Keyer
	Valuer[T]
}

type KeyValueScanner[T any] interface {
	KeyValuer[T]
	Scanner[T]
}

type Keyer interface {
	IsAutoIncr() bool
	PK() (string, int, any)
}

type Valuer[T any] interface {
	Table() string
	Columns() []string
	Values() []any
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
