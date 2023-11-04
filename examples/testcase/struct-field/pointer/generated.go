// Code generated by sqlgen, version v1.0.0-beta. DO NOT EDIT.

package pointer

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Ptr) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `ptr` (`id` BIGINT NOT NULL AUTO_INCREMENT,`str` VARCHAR(255),`uint_8` TINYINT UNSIGNED,`uint_16` SMALLINT UNSIGNED,`uint_32` MEDIUMINT UNSIGNED,`uint_64` BIGINT UNSIGNED,`f_32` VARCHAR(255),`f_64` VARCHAR(255),`time` DATETIME(6),`bytes` BLOB,`bool` TINYINT,`int` INTEGER,`int_8` TINYINT,`int_16` SMALLINT,`int_32` MEDIUMINT,`int_64` BIGINT,`uint` INTEGER UNSIGNED,PRIMARY KEY (`id`));"
}
func (Ptr) AlterTableStmt() string {
	return "ALTER TABLE `ptr` MODIFY `id` BIGINT NOT NULL AUTO_INCREMENT,MODIFY `str` VARCHAR(255) AFTER `id`,MODIFY `uint_8` TINYINT UNSIGNED AFTER `str`,MODIFY `uint_16` SMALLINT UNSIGNED AFTER `uint_8`,MODIFY `uint_32` MEDIUMINT UNSIGNED AFTER `uint_16`,MODIFY `uint_64` BIGINT UNSIGNED AFTER `uint_32`,MODIFY `f_32` VARCHAR(255) AFTER `uint_64`,MODIFY `f_64` VARCHAR(255) AFTER `f_32`,MODIFY `time` DATETIME(6) AFTER `f_64`,MODIFY `bytes` BLOB AFTER `time`,MODIFY `bool` TINYINT AFTER `bytes`,MODIFY `int` INTEGER AFTER `bool`,MODIFY `int_8` TINYINT AFTER `int`,MODIFY `int_16` SMALLINT AFTER `int_8`,MODIFY `int_32` MEDIUMINT AFTER `int_16`,MODIFY `int_64` BIGINT AFTER `int_32`,MODIFY `uint` INTEGER UNSIGNED AFTER `int_64`;"
}
func (Ptr) TableName() string {
	return "`ptr`"
}
func (Ptr) Columns() []string {
	return []string{"`id`", "`str`", "`uint_8`", "`uint_16`", "`uint_32`", "`uint_64`", "`f_32`", "`f_64`", "`time`", "`bytes`", "`bool`", "`int`", "`int_8`", "`int_16`", "`int_32`", "`int_64`", "`uint`"}
}
func (v Ptr) IsAutoIncr() bool {
	return true
}
func (v Ptr) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, int64(v.ID)
}
func (v Ptr) Values() []any {
	return []any{int64(v.ID), types.String(v.Str), types.Integer(v.Uint8), types.Integer(v.Uint16), types.Integer(v.Uint32), types.Integer(v.Uint64), types.Float(v.F32), types.Float(v.F64), types.Time(v.Time), types.String(v.Bytes), types.Bool(v.Bool), types.Integer(v.Int), types.Integer(v.Int8), types.Integer(v.Int16), types.Integer(v.Int32), types.Integer(v.Int64), types.Integer(v.Uint)}
}
func (v *Ptr) Addrs() []any {
	return []any{types.Integer(&v.ID), types.PtrOfString(&v.Str), types.PtrOfInt(&v.Uint8), types.PtrOfInt(&v.Uint16), types.PtrOfInt(&v.Uint32), types.PtrOfInt(&v.Uint64), types.PtrOfFloat(&v.F32), types.PtrOfFloat(&v.F64), types.PtrOfTime(&v.Time), types.PtrOfString(&v.Bytes), types.PtrOfBool(&v.Bool), types.PtrOfInt(&v.Int), types.PtrOfInt(&v.Int8), types.PtrOfInt(&v.Int16), types.PtrOfInt(&v.Int32), types.PtrOfInt(&v.Int64), types.PtrOfInt(&v.Uint)}
}