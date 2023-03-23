package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type Scanner[T any] interface {
	*T
	Addrs() []any
}

type KeyValuer[T any] interface {
	Keyer
	Valuer[T]
}

type Keyer interface {
	// Primary key
	PK() (string, driver.Valuer)
}

type Valuer[T any] interface {
	Table() string
	Columns() []string
	Values() []any
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
