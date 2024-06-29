package date

import (
	"database/sql"
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (User) TableName() string {
	return "`user`"
}
func (User) HasPK() {}
func (v User) PK() (string, int, any) {
	return "`id`", 0, (driver.Valuer)(v.ID)
}
func (User) ColumnNames() []string {
	return []string{"`id`", "`birth_date`"}
}
func (v User) Values() []any {
	return []any{(driver.Valuer)(v.ID), types.TextMarshaler(v.BirthDate)}
}
func (v *User) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID), types.Date(&v.BirthDate)}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO `user` (`id`,`birth_date`) VALUES (?,?);", v.Values()
}
func (v User) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`birth_date` FROM `user` WHERE `id` = ? LIMIT 1;", []any{(driver.Valuer)(v.ID)}
}
func (v User) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `user` SET `birth_date` = ? WHERE `id` = ? LIMIT 1;", []any{types.TextMarshaler(v.BirthDate), (driver.Valuer)(v.ID)}
}
func (v User) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("`id`", v.ID, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
func (v User) GetBirthDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("`birth_date`", v.BirthDate, func(val civil.Date) driver.Value { return types.TextMarshaler(val) })
}
