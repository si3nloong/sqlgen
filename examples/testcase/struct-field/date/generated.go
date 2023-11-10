// Code generated by sqlgen, version v1.0.0-alpha. DO NOT EDIT.

package date

import (
	"database/sql"
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` VARCHAR(36) NOT NULL,`birth_date` DATE NOT NULL,PRIMARY KEY (`id`));"
}
func (User) AlterTableStmt() string {
	return "ALTER TABLE `user` MODIFY `id` VARCHAR(36) NOT NULL,MODIFY `birth_date` DATE NOT NULL AFTER `id`;"
}
func (User) TableName() string {
	return "`user`"
}
func (User) InsertVarStmt() string {
	return "(?,?)"
}
func (User) Columns() []string {
	return []string{"`id`", "`birth_date`"}
}
func (v User) IsAutoIncr() bool {
	return false
}
func (v User) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, (driver.Valuer)(v.ID)
}
func (v User) Values() []any {
	return []any{(driver.Valuer)(v.ID), types.TextMarshaler(v.BirthDate)}
}
func (v *User) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID), types.TextUnmarshaler(&v.BirthDate)}
}