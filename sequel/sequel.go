package sequel

import (
	"context"
	"database/sql"
	"fmt"
	"go/types"
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

type Migrator interface {
	Tabler
	CreateTableStmt() string
	AlterTableStmt() string
}

type MigratorV2 interface {
	Up(ctx context.Context, db DB) error
	Down(ctx context.Context, db DB) error
}

type Stmt interface {
	StmtBuilder
	fmt.Stringer
	Args() []any
	Reset()
}
