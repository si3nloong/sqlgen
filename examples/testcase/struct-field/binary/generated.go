package binary

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Binary) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BINARY(16) NOT NULL,`str` VARCHAR(255) NOT NULL,`time` DATETIME NOT NULL,PRIMARY KEY (`id`));"
}
func (Binary) AlterTableStmt() string {
	return "ALTER TABLE `binary` MODIFY `id` BINARY(16) NOT NULL,MODIFY `str` VARCHAR(255) NOT NULL AFTER `id`,MODIFY `time` DATETIME NOT NULL AFTER `str`;"
}
func (Binary) TableName() string {
	return "`binary`"
}
func (v Binary) InsertOneStmt() string {
	return "INSERT INTO `binary` (`id`,`str`,`time`) VALUES (?,?,?);"
}
func (Binary) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Binary) Columns() []string {
	return []string{"`id`", "`str`", "`time`"}
}
func (v Binary) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, types.BinaryMarshaler(v.ID)
}
func (v Binary) FindByPKStmt() string {
	return "SELECT `id`,`str`,`time` FROM `binary` WHERE `id` = ? LIMIT 1;"
}
func (v Binary) UpdateByPKStmt() string {
	return "UPDATE `binary` SET `str` = ?,`time` = ? WHERE `id` = ? LIMIT 1;"
}
func (v Binary) Values() []any {
	return []any{types.BinaryMarshaler(v.ID), string(v.Str), time.Time(v.Time)}
}
func (v *Binary) Addrs() []any {
	return []any{types.BinaryUnmarshaler(&v.ID), types.String(&v.Str), (*time.Time)(&v.Time)}
}
func (v Binary) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column[uuid.UUID]("`id`", v.ID, func(vi uuid.UUID) driver.Value { return types.BinaryMarshaler(vi) })
}
func (v Binary) GetStr() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`str`", v.Str, func(vi string) driver.Value { return string(vi) })
}
func (v Binary) GetTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("`time`", v.Time, func(vi time.Time) driver.Value { return time.Time(vi) })
}
