package aliasstruct

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`date` DATE NOT NULL,`time` VARCHAR(255) NOT NULL);"
}
func (A) AlterTableStmt() string {
	return "ALTER TABLE `a` MODIFY `date` DATE NOT NULL,MODIFY `time` VARCHAR(255) NOT NULL AFTER `date`;"
}
func (A) InsertVarQuery() string {
	return "(?,?)"
}
func (A) Columns() []string {
	return []string{"`date`", "`time`"}
}
func (v A) Values() []any {
	return []any{types.TextMarshaler(v.Date), types.TextMarshaler(v.Time)}
}
func (v *A) Addrs() []any {
	return []any{types.Date(&v.Date), types.TextUnmarshaler(&v.Time)}
}
func (v A) GetDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column[civil.Date]("`date`", v.Date, func(vi civil.Date) driver.Value { return types.TextMarshaler(vi) })
}
func (v A) GetTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column[civil.Time]("`time`", v.Time, func(vi civil.Time) driver.Value { return types.TextMarshaler(vi) })
}

func (v C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`string` VARCHAR(255) NOT NULL,`valid` TINYINT NOT NULL);"
}
func (C) AlterTableStmt() string {
	return "ALTER TABLE `c` MODIFY `string` VARCHAR(255) NOT NULL,MODIFY `valid` TINYINT NOT NULL AFTER `string`;"
}
func (C) TableName() string {
	return "`c`"
}
func (v C) InsertOneStmt() string {
	return "INSERT INTO `c` (`string`,`valid`) VALUES (?,?);"
}
func (C) InsertVarQuery() string {
	return "(?,?)"
}
func (C) Columns() []string {
	return []string{"`string`", "`valid`"}
}
func (v C) Values() []any {
	return []any{string(v.String), bool(v.Valid)}
}
func (v *C) Addrs() []any {
	return []any{types.String(&v.String), types.Bool(&v.Valid)}
}
func (v C) GetString() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`string`", v.String, func(vi string) driver.Value { return string(vi) })
}
func (v C) GetValid() sequel.ColumnValuer[bool] {
	return sequel.Column[bool]("`valid`", v.Valid, func(vi bool) driver.Value { return bool(vi) })
}
