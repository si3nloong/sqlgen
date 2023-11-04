// Code generated by sqlgen, version v1.0.0-beta. DO NOT EDIT.

package primitive

import (
	"time"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Primitive) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `primitive` (`str` VARCHAR(255) NOT NULL,`bytes` BLOB,`uint_16` SMALLINT UNSIGNED NOT NULL,`uint_32` MEDIUMINT UNSIGNED NOT NULL,`uint_64` BIGINT UNSIGNED NOT NULL,`f_32` VARCHAR(255) NOT NULL,`f_64` VARCHAR(255) NOT NULL,`time` DATETIME NOT NULL,`bool` TINYINT NOT NULL,`int` INTEGER NOT NULL,`int_8` TINYINT NOT NULL,`int_16` SMALLINT NOT NULL,`int_32` MEDIUMINT NOT NULL,`int_64` BIGINT NOT NULL,`uint` INTEGER UNSIGNED NOT NULL,`uint_8` TINYINT UNSIGNED NOT NULL);"
}
func (Primitive) AlterTableStmt() string {
	return "ALTER TABLE `primitive` MODIFY `str` VARCHAR(255) NOT NULL,MODIFY `bytes` BLOB AFTER `str`,MODIFY `uint_16` SMALLINT UNSIGNED NOT NULL AFTER `bytes`,MODIFY `uint_32` MEDIUMINT UNSIGNED NOT NULL AFTER `uint_16`,MODIFY `uint_64` BIGINT UNSIGNED NOT NULL AFTER `uint_32`,MODIFY `f_32` VARCHAR(255) NOT NULL AFTER `uint_64`,MODIFY `f_64` VARCHAR(255) NOT NULL AFTER `f_32`,MODIFY `time` DATETIME NOT NULL AFTER `f_64`,MODIFY `bool` TINYINT NOT NULL AFTER `time`,MODIFY `int` INTEGER NOT NULL AFTER `bool`,MODIFY `int_8` TINYINT NOT NULL AFTER `int`,MODIFY `int_16` SMALLINT NOT NULL AFTER `int_8`,MODIFY `int_32` MEDIUMINT NOT NULL AFTER `int_16`,MODIFY `int_64` BIGINT NOT NULL AFTER `int_32`,MODIFY `uint` INTEGER UNSIGNED NOT NULL AFTER `int_64`,MODIFY `uint_8` TINYINT UNSIGNED NOT NULL AFTER `uint`;"
}
func (Primitive) TableName() string {
	return "`primitive`"
}
func (Primitive) Columns() []string {
	return []string{"`str`", "`bytes`", "`uint_16`", "`uint_32`", "`uint_64`", "`f_32`", "`f_64`", "`time`", "`bool`", "`int`", "`int_8`", "`int_16`", "`int_32`", "`int_64`", "`uint`", "`uint_8`"}
}
func (v Primitive) Values() []any {
	return []any{string(v.Str), string(v.Bytes), int64(v.Uint16), int64(v.Uint32), int64(v.Uint64), float64(v.F32), float64(v.F64), time.Time(v.Time), bool(v.Bool), int64(v.Int), int64(v.Int8), int64(v.Int16), int64(v.Int32), int64(v.Int64), int64(v.Uint), int64(v.Uint8)}
}
func (v *Primitive) Addrs() []any {
	return []any{types.String(&v.Str), types.String(&v.Bytes), types.Integer(&v.Uint16), types.Integer(&v.Uint32), types.Integer(&v.Uint64), types.Float(&v.F32), types.Float(&v.F64), (*time.Time)(&v.Time), types.Bool(&v.Bool), types.Integer(&v.Int), types.Integer(&v.Int8), types.Integer(&v.Int16), types.Integer(&v.Int32), types.Integer(&v.Int64), types.Integer(&v.Uint), types.Integer(&v.Uint8)}
}
