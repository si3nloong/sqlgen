package sql

func Var(n int) string {
	return "?"
}

func Wrap(v string) string {
	return "`" + v + "`"
}
