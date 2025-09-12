package mysql

import "github.com/si3nloong/sqlgen/cmd/codegen/dialect"

type mysqlDriver struct{}

var (
	_ dialect.Dialect = (*mysqlDriver)(nil)
)

func init() {
	dialect.RegisterDialect("mysql", &mysqlDriver{})
}

func (*mysqlDriver) Driver() string {
	return "mysql"
}

func (*mysqlDriver) Var() string {
	return "?"
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
