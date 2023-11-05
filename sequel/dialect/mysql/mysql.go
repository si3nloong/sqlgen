package mysql

import "github.com/si3nloong/sqlgen/sequel"

type mysqlDriver struct{}

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
	if v[0] == '`' {
		return v
	}
	return "`" + v + "`"
}
