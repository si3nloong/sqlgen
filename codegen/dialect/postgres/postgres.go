package postgres

import (
	"strconv"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

type postgresDriver struct{}

var (
	_ dialect.Dialect = (*postgresDriver)(nil)
)

func init() {
	dialect.RegisterDialect("postgres", &postgresDriver{})
}

func (*postgresDriver) Driver() string {
	return "postgres"
}

func (*postgresDriver) Var() string {
	return "$%d"
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
