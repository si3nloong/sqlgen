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

	TableSchemas(model TableSchema) TableDefinition
}

type TableSchema interface {
	GoName() string
	DatabaseName() string
	TableName() string
	AutoIncrKey() (ColumnSchema, bool)
	Keys() []ColumnSchema
	Columns() []ColumnSchema
	Indexes() []IndexSchema
	Implements(*types.Interface) (*types.Func, bool)
}

type ColumnSchema interface {
	GoName() string
	GoPath() string
	// GoTag() reflect.StructTag
	ColumnName() string
	ColumnPos() int
	Size() int
	Type() types.Type
	SQLValuer() QueryFunc
	SQLScanner() QueryFunc
}

type IndexSchema interface {
	Name() string
	Type() string
	ColumnNames() []string
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
	Schemas() TableDefinition
}
