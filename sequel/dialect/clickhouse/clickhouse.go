package clickhouse

import "github.com/si3nloong/sqlgen/sequel"

type clickhouseDriver struct{}

var (
	_ sequel.Dialect = (*clickhouseDriver)(nil)
)

func init() {
	sequel.RegisterDialect("clickhouse", &clickhouseDriver{})
}

func (*clickhouseDriver) Driver() string {
	return "clickhouse"
}

func (*clickhouseDriver) Var(n int) string {
	return "?"
}

func (*clickhouseDriver) Wrap(v string) string {
	return "`" + v + "`"
}

func (*clickhouseDriver) QuoteChar() rune {
	return '\''
}
