package main

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` VARCHAR(36) NOT NULL,`name` VARCHAR(255) NOT NULL);"
}
func (User) AlterTableStmt() string {
	return "ALTER TABLE `user` MODIFY `id` VARCHAR(36) NOT NULL,MODIFY `name` VARCHAR(255) NOT NULL AFTER `id`;"
}
func (User) TableName() string {
	return "`user`"
}
func (v User) InsertOneStmt() string {
	return "INSERT INTO `user` (`id`,`name`) VALUES (?,?);"
}
func (User) InsertVarQuery() string {
	return "(?,?)"
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
func (v User) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column[uuid.UUID]("`id`", v.ID, func(vi uuid.UUID) driver.Value { return (driver.Valuer)(vi) })
}
func (v User) GetName() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`name`", v.Name, func(vi string) driver.Value { return string(vi) })
}
