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
