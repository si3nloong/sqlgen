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

func (*mysqlDriver) VarRune() rune {
	return '?'
}

func (*mysqlDriver) QuoteVar(_ int) string {
	return "?"
}

func (*mysqlDriver) QuoteIdentifier(v string) string {
	return "`" + v + "`"
}

func (*mysqlDriver) QuoteRune() rune {
	return '`'
}
