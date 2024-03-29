package primitive

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Primitive) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`str` VARCHAR(255) NOT NULL,`bytes` BLOB,`uint_16` SMALLINT UNSIGNED NOT NULL,`uint_32` MEDIUMINT UNSIGNED NOT NULL,`uint_64` BIGINT UNSIGNED NOT NULL,`f_32` FLOAT NOT NULL,`f_64` FLOAT NOT NULL,`time` DATETIME NOT NULL,`bool` TINYINT NOT NULL,`int` INTEGER NOT NULL,`int_8` TINYINT NOT NULL,`int_16` SMALLINT NOT NULL,`int_32` MEDIUMINT NOT NULL,`int_64` BIGINT NOT NULL,`uint` INTEGER UNSIGNED NOT NULL,`uint_8` TINYINT UNSIGNED NOT NULL);"
}
func (Primitive) AlterTableStmt() string {
	return "ALTER TABLE `primitive` MODIFY `str` VARCHAR(255) NOT NULL,MODIFY `bytes` BLOB AFTER `str`,MODIFY `uint_16` SMALLINT UNSIGNED NOT NULL AFTER `bytes`,MODIFY `uint_32` MEDIUMINT UNSIGNED NOT NULL AFTER `uint_16`,MODIFY `uint_64` BIGINT UNSIGNED NOT NULL AFTER `uint_32`,MODIFY `f_32` FLOAT NOT NULL AFTER `uint_64`,MODIFY `f_64` FLOAT NOT NULL AFTER `f_32`,MODIFY `time` DATETIME NOT NULL AFTER `f_64`,MODIFY `bool` TINYINT NOT NULL AFTER `time`,MODIFY `int` INTEGER NOT NULL AFTER `bool`,MODIFY `int_8` TINYINT NOT NULL AFTER `int`,MODIFY `int_16` SMALLINT NOT NULL AFTER `int_8`,MODIFY `int_32` MEDIUMINT NOT NULL AFTER `int_16`,MODIFY `int_64` BIGINT NOT NULL AFTER `int_32`,MODIFY `uint` INTEGER UNSIGNED NOT NULL AFTER `int_64`,MODIFY `uint_8` TINYINT UNSIGNED NOT NULL AFTER `uint`;"
}
func (Primitive) TableName() string {
	return "`primitive`"
}
func (v Primitive) InsertOneStmt() string {
	return "INSERT INTO `primitive` (`str`,`bytes`,`uint_16`,`uint_32`,`uint_64`,`f_32`,`f_64`,`time`,`bool`,`int`,`int_8`,`int_16`,`int_32`,`int_64`,`uint`,`uint_8`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
}
func (Primitive) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
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
func (v Primitive) GetStr() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`str`", v.Str, func(vi string) driver.Value { return string(vi) })
}
func (v Primitive) GetBytes() sequel.ColumnValuer[[]byte] {
	return sequel.Column[[]byte]("`bytes`", v.Bytes, func(vi []byte) driver.Value { return string(vi) })
}
func (v Primitive) GetUint16() sequel.ColumnValuer[uint16] {
	return sequel.Column[uint16]("`uint_16`", v.Uint16, func(vi uint16) driver.Value { return int64(vi) })
}
func (v Primitive) GetUint32() sequel.ColumnValuer[uint32] {
	return sequel.Column[uint32]("`uint_32`", v.Uint32, func(vi uint32) driver.Value { return int64(vi) })
}
func (v Primitive) GetUint64() sequel.ColumnValuer[uint64] {
	return sequel.Column[uint64]("`uint_64`", v.Uint64, func(vi uint64) driver.Value { return int64(vi) })
}
func (v Primitive) GetF32() sequel.ColumnValuer[float32] {
	return sequel.Column[float32]("`f_32`", v.F32, func(vi float32) driver.Value { return float64(vi) })
}
func (v Primitive) GetF64() sequel.ColumnValuer[float64] {
	return sequel.Column[float64]("`f_64`", v.F64, func(vi float64) driver.Value { return float64(vi) })
}
func (v Primitive) GetTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("`time`", v.Time, func(vi time.Time) driver.Value { return time.Time(vi) })
}
func (v Primitive) GetBool() sequel.ColumnValuer[bool] {
	return sequel.Column[bool]("`bool`", v.Bool, func(vi bool) driver.Value { return bool(vi) })
}
func (v Primitive) GetInt() sequel.ColumnValuer[int] {
	return sequel.Column[int]("`int`", v.Int, func(vi int) driver.Value { return int64(vi) })
}
func (v Primitive) GetInt8() sequel.ColumnValuer[int8] {
	return sequel.Column[int8]("`int_8`", v.Int8, func(vi int8) driver.Value { return int64(vi) })
}
func (v Primitive) GetInt16() sequel.ColumnValuer[int16] {
	return sequel.Column[int16]("`int_16`", v.Int16, func(vi int16) driver.Value { return int64(vi) })
}
func (v Primitive) GetInt32() sequel.ColumnValuer[int32] {
	return sequel.Column[int32]("`int_32`", v.Int32, func(vi int32) driver.Value { return int64(vi) })
}
func (v Primitive) GetInt64() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("`int_64`", v.Int64, func(vi int64) driver.Value { return int64(vi) })
}
func (v Primitive) GetUint() sequel.ColumnValuer[uint] {
	return sequel.Column[uint]("`uint`", v.Uint, func(vi uint) driver.Value { return int64(vi) })
}
func (v Primitive) GetUint8() sequel.ColumnValuer[uint8] {
	return sequel.Column[uint8]("`uint_8`", v.Uint8, func(vi uint8) driver.Value { return int64(vi) })
}
