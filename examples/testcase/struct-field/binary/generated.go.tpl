package binary

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Binary) TableName() string {
	return "`binary`"
}
func (Binary) HasPK() {}
func (v Binary) PK() (string, int, any) {
	return "`id`", 0, (driver.Valuer)(v.ID)
}
func (Binary) ColumnNames() []string {
	return []string{"`id`", "`str`", "`time`"}
}
func (v Binary) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.Str), time.Time(v.Time)}
}
func (v *Binary) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID), types.String(&v.Str), (*time.Time)(&v.Time)}
}
func (Binary) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Binary) InsertOneStmt() (string, []any) {
	return "INSERT INTO `binary` (`id`,`str`,`time`) VALUES (?,?,?);", v.Values()
}
func (v Binary) FindByPKStmt() (string, []any) {
	return "SELECT `id`,`str`,`time` FROM `binary` WHERE `id` = ? LIMIT 1;", []any{(driver.Valuer)(v.ID)}
}
func (v Binary) UpdateByPKStmt() (string, []any) {
	return "UPDATE `binary` SET `str` = ?,`time` = ? WHERE `id` = ? LIMIT 1;", []any{string(v.Str), time.Time(v.Time), (driver.Valuer)(v.ID)}
}
func (v Binary) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("`id`", v.ID, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
func (v Binary) GetStr() sequel.ColumnValuer[string] {
	return sequel.Column("`str`", v.Str, func(val string) driver.Value { return string(val) })
}
func (v Binary) GetTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("`time`", v.Time, func(val time.Time) driver.Value { return time.Time(val) })
}
