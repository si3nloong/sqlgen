//go:build !mysql
// +build !mysql

package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
)

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
