package postgres

import (
	"strconv"

	"github.com/si3nloong/sqlgen/sequel"
)

type postgresDriver struct{}

func init() {
	sequel.RegisterDialect("postgres", &postgresDriver{})
}

func (*postgresDriver) Var(n int) string {
	return "$" + strconv.Itoa(n)
}

func (*postgresDriver) Wrap(v string) string {
	if v[0] == '"' {
		return v
	}
	return strconv.Quote(v)
}
