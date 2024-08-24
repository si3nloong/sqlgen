package sqlite

import (
	"strconv"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

type sqliteDriver struct{}

var (
	_ dialect.Dialect = (*sqliteDriver)(nil)
)

func init() {
	dialect.RegisterDialect("sqlite", &sqliteDriver{})
}

func (*sqliteDriver) Driver() string {
	return "sqlite"
}

func (*sqliteDriver) Var() string {
	return "?"
}

func (*sqliteDriver) VarRune() rune {
	return '?'
}

func (*sqliteDriver) QuoteVar(_ int) string {
	return "?"
}

func (*sqliteDriver) QuoteIdentifier(v string) string {
	return strconv.Quote(v)
}

func (*sqliteDriver) QuoteRune() rune {
	return '"'
}
