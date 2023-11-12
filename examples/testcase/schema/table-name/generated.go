// Code generated by sqlgen, version v1.1.0-alpha. DO NOT EDIT.

package tablename

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v CustomTableName1) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`text` VARCHAR(255) NOT NULL);"
}
func (CustomTableName1) AlterTableStmt() string {
	return "ALTER TABLE `custom_table_name_1` MODIFY `text` VARCHAR(255) NOT NULL;"
}
func (CustomTableName1) TableName() string {
	return "`custom_table_name_1`"
}
func (CustomTableName1) InsertVarStmt() string {
	return "(?)"
}
func (CustomTableName1) Columns() []string {
	return []string{"`text`"}
}
func (v CustomTableName1) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName1) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (v CustomTableName1) Get_Text() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`text`", v.Text, func(vi string) driver.Value { return string(vi) })
}

func (v CustomTableName2) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`text` VARCHAR(255) NOT NULL);"
}
func (CustomTableName2) AlterTableStmt() string {
	return "ALTER TABLE `custom_table_2` MODIFY `text` VARCHAR(255) NOT NULL;"
}
func (CustomTableName2) TableName() string {
	return "`custom_table_2`"
}
func (CustomTableName2) InsertVarStmt() string {
	return "(?)"
}
func (CustomTableName2) Columns() []string {
	return []string{"`text`"}
}
func (v CustomTableName2) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName2) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (v CustomTableName2) Get_Text() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`text`", v.Text, func(vi string) driver.Value { return string(vi) })
}

func (v CustomTableName3) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`text` VARCHAR(255) NOT NULL);"
}
func (CustomTableName3) AlterTableStmt() string {
	return "ALTER TABLE `custom_table3` MODIFY `text` VARCHAR(255) NOT NULL;"
}
func (CustomTableName3) TableName() string {
	return "`custom_table3`"
}
func (CustomTableName3) InsertVarStmt() string {
	return "(?)"
}
func (CustomTableName3) Columns() []string {
	return []string{"`text`"}
}
func (v CustomTableName3) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName3) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (v CustomTableName3) Get_Text() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`text`", v.Text, func(vi string) driver.Value { return string(vi) })
}
