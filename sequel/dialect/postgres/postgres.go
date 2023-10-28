package postgres

import "strconv"

type postgresDriver struct{}

func (*postgresDriver) Var(n int) string {
	return "$" + strconv.Itoa(n)
}

func (*postgresDriver) Wrap(v string) string {
	return strconv.Quote(v)
}
