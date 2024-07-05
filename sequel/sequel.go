package sequel

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"go/types"
)

type (
	ConvertFunc[T any] func(T) driver.Value
	QueryFunc          func(placeholder string) string
)

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

	TableSchemas(model GoTableSchema) TableSchema
}

type GoTableSchema interface {
	GoName() string
	DatabaseName() string
	TableName() string
	AutoIncrKey() (GoColumnSchema, bool)
	Keys() []string
	Columns() []string
	Indexes() []string
	// Key(i int) GoColumnSchema
	Column(i int) GoColumnSchema
	Index(i int) GoIndexSchema
	Implements(*types.Interface) (*types.Func, bool)
}

type GoColumnSchema interface {
	GoName() string
	GoPath() string
	AutoIncr() bool
	DataType() (string, bool)
	// GoTag() reflect.StructTag
	ColumnName() string
	ColumnPos() int
	Size() int64
	Type() types.Type
	SQLValuer() QueryFunc
	SQLScanner() QueryFunc
}

type GoIndexSchema interface {
	Columns() []string
	Type() string
}
