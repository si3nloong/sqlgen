package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Array) CreateTableStmt() string {
	return `CREATE TABLE IF NOT EXISTS ` + v.TableName() + ` ("tuple" VARCHAR(255) NOT NULL);`
}
func (Array) TableName() string {
	return `"array"`
}
func (Array) InsertOneStmt() string {
	return `INSERT INTO "array" ("tuple") VALUES ($1);`
}
func (Array) Columns() []string {
	return []string{`"tuple"`}
}
func (v Array) Values() []any {
	return []any{encoding.MarshalStringList(v.Tuple[:])}
}
func (v *Array) Addrs() []any {
	v1 := v.Tuple[:]
	return []any{types.StringList(&v1)}
}
func (v Array) GetTuple() sequel.ColumnValuer[[2]string] {
	return sequel.Column(`"tuple"`, v.Tuple, func(val [2]string) driver.Value { return encoding.MarshalStringList(val[:]) })
}
