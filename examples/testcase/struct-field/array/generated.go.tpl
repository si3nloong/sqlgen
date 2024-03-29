package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Array) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,`bool_list` JSON NOT NULL,`uint_8_list` JSON NOT NULL,`uint_16_list` JSON NOT NULL,`uint_32_list` JSON NOT NULL,`uint_64_list` JSON NOT NULL,`f_32_list` JSON NOT NULL,`f_64_list` JSON NOT NULL,`str_list` JSON NOT NULL,`custom_str_list` JSON NOT NULL,`int_list` JSON NOT NULL,`int_8_list` JSON NOT NULL,`int_16_list` JSON NOT NULL,`int_32_list` JSON NOT NULL,`int_64_list` JSON NOT NULL,`uint_list` JSON NOT NULL,PRIMARY KEY (`id`));"
}
func (Array) AlterTableStmt() string {
	return "ALTER TABLE `array` MODIFY `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,MODIFY `bool_list` JSON NOT NULL AFTER `id`,MODIFY `uint_8_list` JSON NOT NULL AFTER `bool_list`,MODIFY `uint_16_list` JSON NOT NULL AFTER `uint_8_list`,MODIFY `uint_32_list` JSON NOT NULL AFTER `uint_16_list`,MODIFY `uint_64_list` JSON NOT NULL AFTER `uint_32_list`,MODIFY `f_32_list` JSON NOT NULL AFTER `uint_64_list`,MODIFY `f_64_list` JSON NOT NULL AFTER `f_32_list`,MODIFY `str_list` JSON NOT NULL AFTER `f_64_list`,MODIFY `custom_str_list` JSON NOT NULL AFTER `str_list`,MODIFY `int_list` JSON NOT NULL AFTER `custom_str_list`,MODIFY `int_8_list` JSON NOT NULL AFTER `int_list`,MODIFY `int_16_list` JSON NOT NULL AFTER `int_8_list`,MODIFY `int_32_list` JSON NOT NULL AFTER `int_16_list`,MODIFY `int_64_list` JSON NOT NULL AFTER `int_32_list`,MODIFY `uint_list` JSON NOT NULL AFTER `int_64_list`;"
}
func (Array) TableName() string {
	return "`array`"
}
func (v Array) InsertOneStmt() string {
	return "INSERT INTO `array` (`bool_list`,`uint_8_list`,`uint_16_list`,`uint_32_list`,`uint_64_list`,`f_32_list`,`f_64_list`,`str_list`,`custom_str_list`,`int_list`,`int_8_list`,`int_16_list`,`int_32_list`,`int_64_list`,`uint_list`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
}
func (Array) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (Array) Columns() []string {
	return []string{"`id`", "`bool_list`", "`uint_8_list`", "`uint_16_list`", "`uint_32_list`", "`uint_64_list`", "`f_32_list`", "`f_64_list`", "`str_list`", "`custom_str_list`", "`int_list`", "`int_8_list`", "`int_16_list`", "`int_32_list`", "`int_64_list`", "`uint_list`"}
}
func (Array) IsAutoIncr() {}
func (v Array) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, int64(v.ID)
}
func (v Array) FindByPKStmt() string {
	return "SELECT `id`,`bool_list`,`uint_8_list`,`uint_16_list`,`uint_32_list`,`uint_64_list`,`f_32_list`,`f_64_list`,`str_list`,`custom_str_list`,`int_list`,`int_8_list`,`int_16_list`,`int_32_list`,`int_64_list`,`uint_list` FROM `array` WHERE `id` = ? LIMIT 1;"
}
func (v Array) UpdateByPKStmt() string {
	return "UPDATE `array` SET `bool_list` = ?,`uint_8_list` = ?,`uint_16_list` = ?,`uint_32_list` = ?,`uint_64_list` = ?,`f_32_list` = ?,`f_64_list` = ?,`str_list` = ?,`custom_str_list` = ?,`int_list` = ?,`int_8_list` = ?,`int_16_list` = ?,`int_32_list` = ?,`int_64_list` = ?,`uint_list` = ? WHERE `id` = ? LIMIT 1;"
}
func (v Array) Values() []any {
	return []any{int64(v.ID), encoding.MarshalBoolList(v.BoolList), encoding.MarshalUnsignedIntList(v.Uint8List), encoding.MarshalUnsignedIntList(v.Uint16List), encoding.MarshalUnsignedIntList(v.Uint32List), encoding.MarshalUnsignedIntList(v.Uint64List), encoding.MarshalFloatList(v.F32List), encoding.MarshalFloatList(v.F64List), encoding.MarshalStringList(v.StrList), encoding.MarshalStringList(v.CustomStrList), encoding.MarshalSignedIntList(v.IntList), encoding.MarshalSignedIntList(v.Int8List), encoding.MarshalSignedIntList(v.Int16List), encoding.MarshalSignedIntList(v.Int32List), encoding.MarshalSignedIntList(v.Int64List), encoding.MarshalUnsignedIntList(v.UintList)}
}
func (v *Array) Addrs() []any {
	return []any{types.Integer(&v.ID), types.BoolList(&v.BoolList), types.UintList(&v.Uint8List), types.UintList(&v.Uint16List), types.UintList(&v.Uint32List), types.UintList(&v.Uint64List), types.FloatList(&v.F32List), types.FloatList(&v.F64List), types.StringList(&v.StrList), types.StringList(&v.CustomStrList), types.IntList(&v.IntList), types.IntList(&v.Int8List), types.IntList(&v.Int16List), types.IntList(&v.Int32List), types.IntList(&v.Int64List), types.UintList(&v.UintList)}
}
func (v Array) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column[uint64]("`id`", v.ID, func(vi uint64) driver.Value { return int64(vi) })
}
func (v Array) GetBoolList() sequel.ColumnValuer[[]bool] {
	return sequel.Column[[]bool]("`bool_list`", v.BoolList, func(vi []bool) driver.Value { return encoding.MarshalBoolList(vi) })
}
func (v Array) GetUint8List() sequel.ColumnValuer[[]uint8] {
	return sequel.Column[[]uint8]("`uint_8_list`", v.Uint8List, func(vi []uint8) driver.Value { return encoding.MarshalUnsignedIntList(vi) })
}
func (v Array) GetUint16List() sequel.ColumnValuer[[]uint16] {
	return sequel.Column[[]uint16]("`uint_16_list`", v.Uint16List, func(vi []uint16) driver.Value { return encoding.MarshalUnsignedIntList(vi) })
}
func (v Array) GetUint32List() sequel.ColumnValuer[[]uint32] {
	return sequel.Column[[]uint32]("`uint_32_list`", v.Uint32List, func(vi []uint32) driver.Value { return encoding.MarshalUnsignedIntList(vi) })
}
func (v Array) GetUint64List() sequel.ColumnValuer[[]uint64] {
	return sequel.Column[[]uint64]("`uint_64_list`", v.Uint64List, func(vi []uint64) driver.Value { return encoding.MarshalUnsignedIntList(vi) })
}
func (v Array) GetF32List() sequel.ColumnValuer[[]float32] {
	return sequel.Column[[]float32]("`f_32_list`", v.F32List, func(vi []float32) driver.Value { return encoding.MarshalFloatList(vi) })
}
func (v Array) GetF64List() sequel.ColumnValuer[[]float64] {
	return sequel.Column[[]float64]("`f_64_list`", v.F64List, func(vi []float64) driver.Value { return encoding.MarshalFloatList(vi) })
}
func (v Array) GetStrList() sequel.ColumnValuer[[]string] {
	return sequel.Column[[]string]("`str_list`", v.StrList, func(vi []string) driver.Value { return encoding.MarshalStringList(vi) })
}
func (v Array) GetCustomStrList() sequel.ColumnValuer[[]customStr] {
	return sequel.Column[[]customStr]("`custom_str_list`", v.CustomStrList, func(vi []customStr) driver.Value { return encoding.MarshalStringList(vi) })
}
func (v Array) GetIntList() sequel.ColumnValuer[[]int] {
	return sequel.Column[[]int]("`int_list`", v.IntList, func(vi []int) driver.Value { return encoding.MarshalSignedIntList(vi) })
}
func (v Array) GetInt8List() sequel.ColumnValuer[[]int8] {
	return sequel.Column[[]int8]("`int_8_list`", v.Int8List, func(vi []int8) driver.Value { return encoding.MarshalSignedIntList(vi) })
}
func (v Array) GetInt16List() sequel.ColumnValuer[[]int16] {
	return sequel.Column[[]int16]("`int_16_list`", v.Int16List, func(vi []int16) driver.Value { return encoding.MarshalSignedIntList(vi) })
}
func (v Array) GetInt32List() sequel.ColumnValuer[[]int32] {
	return sequel.Column[[]int32]("`int_32_list`", v.Int32List, func(vi []int32) driver.Value { return encoding.MarshalSignedIntList(vi) })
}
func (v Array) GetInt64List() sequel.ColumnValuer[[]int64] {
	return sequel.Column[[]int64]("`int_64_list`", v.Int64List, func(vi []int64) driver.Value { return encoding.MarshalSignedIntList(vi) })
}
func (v Array) GetUintList() sequel.ColumnValuer[[]uint] {
	return sequel.Column[[]uint]("`uint_list`", v.UintList, func(vi []uint) driver.Value { return encoding.MarshalUnsignedIntList(vi) })
}
