package valuer

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (B) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` BIGINT NOT NULL DEFAULT 0"},
			{Name: "`value`", Definition: "`value` VARCHAR(255) NOT NULL"},
			{Name: "`n`", Definition: "`n` VARCHAR(255) NOT NULL DEFAULT ''"},
		},
	}
}
func (B) TableName() string {
	return "`b`"
}
func (B) ColumnNames() []string {
	return []string{"`id`", "`value`", "`n`"}
}
func (v B) Values() []any {
	return []any{int64(v.ID), (driver.Valuer)(v.Value), string(v.N)}
}
func (v *B) Addrs() []any {
	return []any{types.Integer(&v.ID), types.JSONUnmarshaler(&v.Value), types.String(&v.N)}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO `b` (`id`,`value`,`n`) VALUES (?,?,?);", v.Values()
}
func (v B) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("`id`", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v B) GetValue() sequel.ColumnValuer[anyType] {
	return sequel.Column("`value`", v.Value, func(val anyType) driver.Value { return (driver.Valuer)(val) })
}
func (v B) GetN() sequel.ColumnValuer[string] {
	return sequel.Column("`n`", v.N, func(val string) driver.Value { return string(val) })
}
