package postgres

import (
	"strconv"

	"github.com/si3nloong/sqlgen/sequel"
)

type postgresDriver struct{}

var (
	_ sequel.Dialect = (*postgresDriver)(nil)
)

func init() {
	sequel.RegisterDialect("postgres", &postgresDriver{})
}

func (*postgresDriver) Driver() string {
	return "postgres"
}

func (*postgresDriver) VarRune() rune {
	return '$'
}

func (s postgresDriver) QuoteVar(n int) string {
	return string(s.VarRune()) + strconv.Itoa(n)
}

func (*postgresDriver) QuoteIdentifier(v string) string {
	return strconv.Quote(v)
}

func (*postgresDriver) QuoteRune() rune {
	return '"'
}
