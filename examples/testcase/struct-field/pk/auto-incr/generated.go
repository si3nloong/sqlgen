package pkautoincr

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Model) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"`id`"},
			Definition: "PRIMARY KEY (`id`)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "`name`", Definition: "`name` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`f`", Definition: "`f` BOOL NOT NULL DEFAULT false"},
			{Name: "`id`", Definition: "`id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT"},
			{Name: "`n`", Definition: "`n` BIGINT NOT NULL DEFAULT 0"},
		},
	}
}
func (Model) TableName() string {
	return "`AutoIncrPK`"
}
func (Model) HasPK()      {}
func (Model) IsAutoIncr() {}
func (v Model) PK() (string, int, any) {
	return "`id`", 2, int64(v.ID)
}
func (Model) Columns() []string {
	return []string{"`name`", "`f`", "`id`", "`n`"}
}
func (v Model) Values() []any {
	return []any{string(v.Name), bool(v.F), int64(v.ID), int64(v.N)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name), types.Bool(&v.F), types.Integer(&v.ID), types.Integer(&v.N)}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO `AutoIncrPK` (`name`,`f`,`n`) VALUES (?,?,?);", []any{string(v.Name), bool(v.F), int64(v.N)}
}
func (v Model) FindOneByPKStmt() (string, []any) {
	return "SELECT `name`,`f`,`id`,`n` FROM `AutoIncrPK` WHERE `id` = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v Model) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `AutoIncrPK` SET `name` = ?,`f` = ?,`n` = ? WHERE `id` = ? LIMIT 1;", []any{string(v.Name), bool(v.F), int64(v.N), int64(v.ID)}
}
func (v Model) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column("`name`", v.Name, func(val LongText) driver.Value { return string(val) })
}
func (v Model) GetF() sequel.ColumnValuer[Flag] {
	return sequel.Column("`f`", v.F, func(val Flag) driver.Value { return bool(val) })
}
func (v Model) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column("`id`", v.ID, func(val uint) driver.Value { return int64(val) })
}
func (v Model) GetN() sequel.ColumnValuer[int64] {
	return sequel.Column("`n`", v.N, func(val int64) driver.Value { return int64(val) })
}
