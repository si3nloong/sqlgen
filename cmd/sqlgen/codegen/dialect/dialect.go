package dialect

import (
	"database/sql/driver"
	"errors"
	"go/types"
	"io"
	"sync"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
)

var (
	dialectMap     sync.Map
	defaultDialect = "mysql"

	ErrSkipMigration = errors.New("no migration required")
)

type UpFunc func(w io.Writer) error
type DownFunc func(w io.Writer) error

type GoExpr string

type Dialect interface {
	// SQL driver name
	Driver() string
	// Argument string to escape SQL injection
	QuoteVar(n int) string
	VarRune() rune
	Var() string

	// Character to escape table, column name
	QuoteIdentifier(v string) string

	// Quote rune can be ' or " or `
	QuoteRune() rune

	// Column data types
	ColumnDataTypes() map[string]*ColumnType

	// To create migration
	// Migrate(ctx context.Context, t *compiler.Table) (UpFunc, DownFunc)
}

type Column interface {
	DataType(column compiler.Column) string
	Scanner() GoExpr
	Valuer() GoExpr
	SQLScanner() (GoExpr, bool)
	SQLValuer() (GoExpr, bool)
}

type ColumnType struct {
	DataType   func(col GoColumn) string
	Scanner    string
	Valuer     string
	SQLScanner *string
	SQLValuer  *string
}

type GoColumn interface {
	// Go field name
	GoName() string

	// Go field path
	//	eg. A.nested.Field
	GoPath() string

	// Go actual type
	GoType() types.Type

	// Is it a go nullable type? such as
	// pointer, slice, map or chan
	GoNullable() bool

	// SQL column name
	ColumnName() string

	// SQL data type
	DataType() string

	// SQL default value, this can be
	// 	string, []byte, bool, int64, float64, sql.RawBytes
	Default() (driver.Value, bool)

	// Determine whether this column is auto increment or not
	AutoIncr() bool

	// Key is to identify whether column is primary or foreign key
	Key() bool

	// Column size that declared by user
	Size() int

	// CharacterMaxLength() (int64, bool)

	// NumericPrecision() (int64, bool)

	// DatetimePrecision() (int64, bool)
}

func RegisterDialect(name string, d Dialect) {
	if d == nil {
		panic("sqlgen: cannot register nil as dialect")
	}
	dialectMap.Store(name, d)
}

func GetDialect(name string) (Dialect, bool) {
	v, ok := dialectMap.Load(name)
	if !ok {
		return nil, ok
	}
	return v.(Dialect), ok
}

func DefaultDialect() Dialect {
	v, _ := dialectMap.Load(defaultDialect)
	return v.(Dialect)
}
