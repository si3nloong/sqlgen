// Code generated by sqlgen, version v1.0.0-alpha.1. DO NOT EDIT.

package tabler

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`name` VARCHAR(255) NOT NULL);"
}
func (Model) AlterTableStmt() string {
	return "ALTER TABLE `model` MODIFY `name` VARCHAR(255) NOT NULL;"
}
func (v Model) InsertOneStmt() string {
	return "INSERT INTO `model` (`name`) VALUES (?);"
}
func (Model) InsertVarQuery() string {
	return "(?)"
}
func (Model) Columns() []string {
	return []string{"`name`"}
}
func (v Model) Values() []any {
	return []any{string(v.Name)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name)}
}
func (v Model) GetName() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`name`", v.Name, func(vi string) driver.Value { return string(vi) })
}
