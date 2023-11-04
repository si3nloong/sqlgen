// Code generated by sqlgen, version v1.0.0-beta. DO NOT EDIT.

package main

import (
	"database/sql"
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `user` (`id` VARCHAR(36) NOT NULL,`name` VARCHAR(255) NOT NULL);"
}
func (User) AlterTableStmt() string {
	return "ALTER TABLE `user` MODIFY `id` VARCHAR(36) NOT NULL,MODIFY `name` VARCHAR(255) NOT NULL AFTER `id`;"
}
func (User) TableName() string {
	return "`user`"
}
func (User) Columns() []string {
	return []string{"`id`", "`name`"}
}
func (v User) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.Name)}
}
func (v *User) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID), types.String(&v.Name)}
}