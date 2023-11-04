package sqlite

import (
	"strconv"

	"github.com/si3nloong/sqlgen/sequel"
)

type sqliteDriver struct{}

func init() {
	sequel.RegisterDialect("sqlite", &sqliteDriver{})
}

func (*sqliteDriver) Var(n int) string {
	return "?"
}

func (*sqliteDriver) Wrap(v string) string {
	if v[0] == '"' {
		return v
	}
	return strconv.Quote(v)
}
