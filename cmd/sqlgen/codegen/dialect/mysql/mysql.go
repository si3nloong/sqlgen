//go:build !mysql
// +build !mysql

package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
)

type mysqlDriver struct{}

var (
	_       dialect.Dialect = (*mysqlDriver)(nil)
	typeMap                 = map[string]string{
		compiler.Byte:    "CHAR",
		compiler.Rune:    "CHAR",
		compiler.String:  "VARCHAR(255)",
		compiler.Bool:    "BOOL",
		compiler.Int:     "INTEGER",
		compiler.Int8:    "TINYINT",
		compiler.Int16:   "SMALLINT",
		compiler.Int32:   "MEDIUMINT",
		compiler.Int64:   "BIGINT",
		compiler.Uint:    "INTEGER UNSIGNED",
		compiler.Uint8:   "TINYINT UNSIGNED",
		compiler.Uint16:  "SMALLINT UNSIGNED",
		compiler.Uint32:  "MEDIUMINT UNSIGNED",
		compiler.Uint64:  "BIGINT UNSIGNED",
		compiler.Float32: "FLOAT",
		compiler.Float64: "FLOAT",
		compiler.Time:    "DATETIME(6)",
	}
)

func init() {
	dialect.RegisterDialect("mysql", &mysqlDriver{})
}

func (mysqlDriver) Driver() string {
	return "mysql"
}

func (mysqlDriver) Var() string {
	return "?"
}

func (mysqlDriver) VarRune() rune {
	return '?'
}

func (mysqlDriver) QuoteVar(_ int) string {
	return "?"
}

func (mysqlDriver) QuoteIdentifier(v string) string {
	return "`" + v + "`"
}

func (mysqlDriver) QuoteRune() rune {
	return '`'
}
