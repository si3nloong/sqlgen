// Code generated by sqlgen, version v1.1.0-alpha. DO NOT EDIT.

package pkautoincr

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`name` VARCHAR(255) NOT NULL,`f` TINYINT NOT NULL,`id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,`n` BIGINT NOT NULL,PRIMARY KEY (`id`));"
}
func (Model) AlterTableStmt() string {
	return "ALTER TABLE `model` MODIFY `name` VARCHAR(255) NOT NULL,MODIFY `f` TINYINT NOT NULL AFTER `name`,MODIFY `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT AFTER `f`,MODIFY `n` BIGINT NOT NULL AFTER `id`;"
}
func (Model) TableName() string {
	return "`model`"
}
func (Model) InsertVarStmt() string {
	return "(?,?,?,?)"
}
func (Model) Columns() []string {
	return []string{"`name`", "`f`", "`id`", "`n`"}
}
func (v Model) IsAutoIncr() bool {
	return true
}
func (v Model) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 2, int64(v.ID)
}
func (v Model) Values() []any {
	return []any{string(v.Name), bool(v.F), int64(v.ID), int64(v.N)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name), types.Bool(&v.F), types.Integer(&v.ID), types.Integer(&v.N)}
}
func (v Model) Get_Name() sequel.ColumnValuer[LongText] {
	return sequel.Column[LongText]("`name`", v.Name, func(vi LongText) driver.Value { return string(vi) })
}
func (v Model) Get_F() sequel.ColumnValuer[Flag] {
	return sequel.Column[Flag]("`f`", v.F, func(vi Flag) driver.Value { return bool(vi) })
}
func (v Model) Get_ID() sequel.ColumnValuer[uint] {
	return sequel.Column[uint]("`id`", v.ID, func(vi uint) driver.Value { return int64(vi) })
}
func (v Model) Get_N() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("`n`", v.N, func(vi int64) driver.Value { return int64(vi) })
}
