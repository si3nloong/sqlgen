package sql

type Dialect interface {
	Var(n int) string
	Wrap(v string) string
}

var (
	dialect Dialect = mysql{}
)

type mysql struct{}

func (mysql) Var(n int) string {
	return "?"
}

func (mysql) Wrap(v string) string {
	return "`" + v + "`"
}
