package mysql

type mysqlDriver struct{}

func (*mysqlDriver) Var(n int) string {
	return "?"
}

func (*mysqlDriver) Wrap(v string) string {
	return "`" + v + "`"
}
