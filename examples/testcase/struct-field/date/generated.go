// Code generated by sqlgen, version v1.1.0-alpha. DO NOT EDIT.

package date

import (
	"database/sql"
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
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
func (v User) InsertOneStmt() string {
	return "INSERT INTO " + v.TableName() + " (`id`,`birth_date`) VALUES (?,?);"
}
func (User) InsertVarQuery() string {
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
	return []any{(sql.Scanner)(&v.ID), types.Date(&v.BirthDate)}
}
func (v User) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column[uuid.UUID]("`id`", v.ID, func(vi uuid.UUID) driver.Value { return (driver.Valuer)(vi) })
}
func (v User) GetBirthDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column[civil.Date]("`birth_date`", v.BirthDate, func(vi civil.Date) driver.Value { return types.TextMarshaler(vi) })
}
