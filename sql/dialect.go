package sql

import "strconv"

type Dialect interface {
	Var(n int) string
	Wrap(v string) string
}

var (
	dialect Dialect = &mysql{}
)

func SetDialect(driver string) {
	switch driver {
	case "mysql":
		dialect = &mysql{}
	case "postgres":
		dialect = &postgres{}
	}
}

type mysql struct{}

func (*mysql) Var(n int) string {
	return "?"
}

func (*mysql) Wrap(v string) string {
	return "`" + v + "`"
}

type postgres struct{}

func (*postgres) Var(n int) string {
	return "$" + strconv.Itoa(n)
}

func (*postgres) Wrap(v string) string {
	return strconv.Quote(v)
}
