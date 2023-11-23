package postgres

import (
	"strconv"

	"github.com/si3nloong/sqlgen/sequel"
)

type postgresDriver struct{}

var (
	_ sequel.Dialect    = (*postgresDriver)(nil)
	_ sequel.DialectVar = (*postgresDriver)(nil)
)

func init() {
	sequel.RegisterDialect("postgres", &postgresDriver{})
}

func (*postgresDriver) Driver() string {
	return "postgres"
}

func (*postgresDriver) VarChar() string {
	return "$"
}

func (s postgresDriver) Var(n int) string {
	return s.VarChar() + strconv.Itoa(n)
}

func (*postgresDriver) Wrap(v string) string {
	return strconv.Quote(v)
}

func (*postgresDriver) QuoteChar() rune {
	return '"'
}
