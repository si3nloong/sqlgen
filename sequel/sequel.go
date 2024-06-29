package sequel

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"go/types"
)

type (
	ConvertFunc[T any] func(T) driver.Value
	SQLFunc            func(placeholder string) string
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

	CreateTableStmt(n string, model TableSchema) string
}

type TableSchema interface {
	GoName() string
	DatabaseName() string
	TableName() string
	AutoIncrKey() (ColumnSchema, bool)
	Implements(*types.Interface) (wrongType bool)
	Keys() []ColumnSchema
	Columns() []ColumnSchema
}

type ColumnSchema interface {
	GoName() string
	GoPath() string
	// GoTag() reflect.StructTag
	ColumnName() string
	ColumnPos() int
	Type() types.Type
	SQLValuer() SQLFunc
	SQLScanner() SQLFunc
	// ActualType() string
	Implements(*types.Interface) (wrongType bool)
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

type Migrator interface {
	// [0] is column name
	// [1] is column data type
	Schemas() [][2]string
}
