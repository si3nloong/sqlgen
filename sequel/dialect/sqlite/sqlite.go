package sqlite

import (
	"strconv"

	"github.com/si3nloong/sqlgen/sequel"
)

type sqliteDriver struct{}

var (
	_ sequel.Dialect = (*sqliteDriver)(nil)
)

func init() {
	sequel.RegisterDialect("sqlite", &sqliteDriver{})
}

func (*sqliteDriver) Driver() string {
	return "sqlite"
}

func (*sqliteDriver) Var(n int) string {
	return "?"
}

func (*sqliteDriver) Wrap(v string) string {
	return strconv.Quote(v)
}

func (*sqliteDriver) QuoteChar() rune {
	return '"'
}
