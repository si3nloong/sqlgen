package imported

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`str` VARCHAR(255) NOT NULL,`bool` VARCHAR(255) NOT NULL,`raw_bytes` BLOB,`int_16` VARCHAR(255) NOT NULL,`int_32` VARCHAR(255) NOT NULL,`int_64` VARCHAR(255) NOT NULL,`time` VARCHAR(255) NOT NULL);"
}
func (Model) AlterTableStmt() string {
	return "ALTER TABLE `model` MODIFY `str` VARCHAR(255) NOT NULL,MODIFY `bool` VARCHAR(255) NOT NULL AFTER `str`,MODIFY `raw_bytes` BLOB AFTER `bool`,MODIFY `int_16` VARCHAR(255) NOT NULL AFTER `raw_bytes`,MODIFY `int_32` VARCHAR(255) NOT NULL AFTER `int_16`,MODIFY `int_64` VARCHAR(255) NOT NULL AFTER `int_32`,MODIFY `time` VARCHAR(255) NOT NULL AFTER `int_64`;"
}
func (Model) TableName() string {
	return "`model`"
}
func (v Model) InsertOneStmt() string {
	return "INSERT INTO `model` (`str`,`bool`,`raw_bytes`,`int_16`,`int_32`,`int_64`,`time`) VALUES (?,?,?,?,?,?,?);"
}
func (Model) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?)"
}
func (Model) Columns() []string {
	return []string{"`str`", "`bool`", "`raw_bytes`", "`int_16`", "`int_32`", "`int_64`", "`time`"}
}
func (v Model) Values() []any {
	return []any{(driver.Valuer)(v.Str), (driver.Valuer)(v.Bool), string(v.RawBytes), (driver.Valuer)(v.Int16), (driver.Valuer)(v.Int32), (driver.Valuer)(v.Int64), (driver.Valuer)(v.Time)}
}
func (v *Model) Addrs() []any {
	return []any{(sql.Scanner)(&v.Str), (sql.Scanner)(&v.Bool), types.String(&v.RawBytes), (sql.Scanner)(&v.Int16), (sql.Scanner)(&v.Int32), (sql.Scanner)(&v.Int64), (sql.Scanner)(&v.Time)}
}
func (v Model) GetStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column[sql.NullString]("`str`", v.Str, func(vi sql.NullString) driver.Value { return (driver.Valuer)(vi) })
}
func (v Model) GetBool() sequel.ColumnValuer[sql.NullBool] {
	return sequel.Column[sql.NullBool]("`bool`", v.Bool, func(vi sql.NullBool) driver.Value { return (driver.Valuer)(vi) })
}
func (v Model) GetRawBytes() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column[sql.RawBytes]("`raw_bytes`", v.RawBytes, func(vi sql.RawBytes) driver.Value { return string(vi) })
}
func (v Model) GetInt16() sequel.ColumnValuer[sql.NullInt16] {
	return sequel.Column[sql.NullInt16]("`int_16`", v.Int16, func(vi sql.NullInt16) driver.Value { return (driver.Valuer)(vi) })
}
func (v Model) GetInt32() sequel.ColumnValuer[sql.NullInt32] {
	return sequel.Column[sql.NullInt32]("`int_32`", v.Int32, func(vi sql.NullInt32) driver.Value { return (driver.Valuer)(vi) })
}
func (v Model) GetInt64() sequel.ColumnValuer[sql.NullInt64] {
	return sequel.Column[sql.NullInt64]("`int_64`", v.Int64, func(vi sql.NullInt64) driver.Value { return (driver.Valuer)(vi) })
}
func (v Model) GetTime() sequel.ColumnValuer[sql.NullTime] {
	return sequel.Column[sql.NullTime]("`time`", v.Time, func(vi sql.NullTime) driver.Value { return (driver.Valuer)(vi) })
}

func (v Some) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` VARCHAR(36) NOT NULL);"
}
func (Some) AlterTableStmt() string {
	return "ALTER TABLE `some` MODIFY `id` VARCHAR(36) NOT NULL;"
}
func (Some) TableName() string {
	return "`some`"
}
func (v Some) InsertOneStmt() string {
	return "INSERT INTO `some` (`id`) VALUES (?);"
}
func (Some) InsertVarQuery() string {
	return "(?)"
}
func (Some) Columns() []string {
	return []string{"`id`"}
}
func (v Some) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *Some) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (v Some) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column[uuid.UUID]("`id`", v.ID, func(vi uuid.UUID) driver.Value { return (driver.Valuer)(vi) })
}
