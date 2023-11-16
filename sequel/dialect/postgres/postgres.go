package postgres

import (
	"strconv"

	"github.com/si3nloong/sqlgen/sequel"
)

type postgresDriver struct{}

func init() {
	sequel.RegisterDialect("postgres", &postgresDriver{})
}

func (*postgresDriver) Driver() string {
	return "postgres"
}

func (*postgresDriver) Var(n int) string {
	return "$" + strconv.Itoa(n)
}

func (*postgresDriver) Wrap(v string) string {
	return strconv.Quote(v)
}
