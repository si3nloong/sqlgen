//go:build !mysql
// +build !mysql

package mysql

import (
	"fmt"
	"go/types"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
)

type mysqlDriver struct {
	typeMap map[fmt.GoStringer]*dialect.ColumnType
}

var (
	_ dialect.Dialect = (*mysqlDriver)(nil)
)

func init() {
	dialect.RegisterDialect("mysql", &mysqlDriver{
		typeMap: map[fmt.GoStringer]*dialect.ColumnType{
			compiler.Byte:    &dialect.ColumnType{},
			compiler.Rune:    nil,
			compiler.String:  nil,
			compiler.Bool:    nil,
			compiler.Int:     nil,
			compiler.Int8:    nil,
			compiler.Int16:   nil,
			compiler.Int32:   nil,
			compiler.Int64:   nil,
			compiler.Uint:    nil,
			compiler.Uint8:   nil,
			compiler.Uint16:  nil,
			compiler.Uint32:  nil,
			compiler.Uint64:  nil,
			compiler.Float32: nil,
			compiler.Float64: nil,
			compiler.Any:     nil,
		},
	})
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

func (s *mysqlDriver) TypeMapper(t types.Type) {
	v, ok := s.typeMap[compiler.GoType{Type: t}]
	if !ok {
		return
	}
	log.Println(v)
	return
}
