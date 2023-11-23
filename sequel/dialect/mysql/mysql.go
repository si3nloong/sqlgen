package mysql

import "github.com/si3nloong/sqlgen/sequel"

type mysqlDriver struct{}

var (
	_ sequel.Dialect = (*mysqlDriver)(nil)
)

func init() {
	sequel.RegisterDialect("mysql", &mysqlDriver{})
}

func (*mysqlDriver) Driver() string {
	return "mysql"
}

func (*mysqlDriver) Var(n int) string {
	return "?"
}

func (*mysqlDriver) Wrap(v string) string {
	return "`" + v + "`"
}

func (*mysqlDriver) QuoteChar() rune {
	return '`'
}
