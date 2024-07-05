package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Array) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`tuple`", Definition: "`tuple` VARCHAR(255) NOT NULL"},
		},
	}
}
func (Array) TableName() string {
	return "`array`"
}
func (Array) ColumnNames() []string {
	return []string{"`tuple`"}
}
func (v Array) Values() []any {
	return []any{types.JSONMarshaler(v.Tuple)}
}
func (v *Array) Addrs() []any {
	return []any{types.JSONUnmarshaler(&v.Tuple)}
}
func (Array) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v Array) InsertOneStmt() (string, []any) {
	return "INSERT INTO `array` (`tuple`) VALUES (?);", v.Values()
}
func (v Array) GetTuple() sequel.ColumnValuer[[2]string] {
	return sequel.Column("`tuple`", v.Tuple, func(val [2]string) driver.Value { return types.JSONMarshaler(val) })
}
