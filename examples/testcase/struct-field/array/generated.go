// Code generated by sqlgen, version v1.0.0-alpha. DO NOT EDIT.

package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Array) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `array` (`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,`bool_list` JSON NOT NULL,`uint_8_list` JSON NOT NULL,`uint_16_list` JSON NOT NULL,`uint_32_list` JSON NOT NULL,`uint_64_list` JSON NOT NULL,`f_32_list` JSON NOT NULL,`f_64_list` JSON NOT NULL,`str_list` JSON NOT NULL,`custom_str_list` JSON NOT NULL,`int_list` JSON NOT NULL,`int_8_list` JSON NOT NULL,`int_16_list` JSON NOT NULL,`int_32_list` JSON NOT NULL,`int_64_list` JSON NOT NULL,`uint_list` JSON NOT NULL,PRIMARY KEY (`id`));"
}
func (Array) AlterTableStmt() string {
	return "ALTER TABLE `array` MODIFY `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,MODIFY `bool_list` JSON NOT NULL AFTER `id`,MODIFY `uint_8_list` JSON NOT NULL AFTER `bool_list`,MODIFY `uint_16_list` JSON NOT NULL AFTER `uint_8_list`,MODIFY `uint_32_list` JSON NOT NULL AFTER `uint_16_list`,MODIFY `uint_64_list` JSON NOT NULL AFTER `uint_32_list`,MODIFY `f_32_list` JSON NOT NULL AFTER `uint_64_list`,MODIFY `f_64_list` JSON NOT NULL AFTER `f_32_list`,MODIFY `str_list` JSON NOT NULL AFTER `f_64_list`,MODIFY `custom_str_list` JSON NOT NULL AFTER `str_list`,MODIFY `int_list` JSON NOT NULL AFTER `custom_str_list`,MODIFY `int_8_list` JSON NOT NULL AFTER `int_list`,MODIFY `int_16_list` JSON NOT NULL AFTER `int_8_list`,MODIFY `int_32_list` JSON NOT NULL AFTER `int_16_list`,MODIFY `int_64_list` JSON NOT NULL AFTER `int_32_list`,MODIFY `uint_list` JSON NOT NULL AFTER `int_64_list`;"
}
func (Array) TableName() string {
	return "`array`"
}
func (Array) Columns() []string {
	return []string{"`id`", "`bool_list`", "`uint_8_list`", "`uint_16_list`", "`uint_32_list`", "`uint_64_list`", "`f_32_list`", "`f_64_list`", "`str_list`", "`custom_str_list`", "`int_list`", "`int_8_list`", "`int_16_list`", "`int_32_list`", "`int_64_list`", "`uint_list`"}
}
func (v Array) IsAutoIncr() bool {
	return true
}
func (v Array) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, int64(v.ID)
}
func (v Array) Values() []any {
	return []any{int64(v.ID), encoding.MarshalBoolList(v.BoolList), encoding.MarshalUnsignedIntList(v.Uint8List), encoding.MarshalUnsignedIntList(v.Uint16List), encoding.MarshalUnsignedIntList(v.Uint32List), encoding.MarshalUnsignedIntList(v.Uint64List), encoding.MarshalFloatList(v.F32List), encoding.MarshalFloatList(v.F64List), encoding.MarshalStringList(v.StrList), encoding.MarshalStringList(v.CustomStrList), encoding.MarshalSignedIntList(v.IntList), encoding.MarshalSignedIntList(v.Int8List), encoding.MarshalSignedIntList(v.Int16List), encoding.MarshalSignedIntList(v.Int32List), encoding.MarshalSignedIntList(v.Int64List), encoding.MarshalUnsignedIntList(v.UintList)}
}
func (v *Array) Addrs() []any {
	return []any{types.Integer(&v.ID), types.BoolList(&v.BoolList), types.UintList(&v.Uint8List), types.UintList(&v.Uint16List), types.UintList(&v.Uint32List), types.UintList(&v.Uint64List), types.FloatList(&v.F32List), types.FloatList(&v.F64List), types.StringList(&v.StrList), types.StringList(&v.CustomStrList), types.IntList(&v.IntList), types.IntList(&v.Int8List), types.IntList(&v.Int16List), types.IntList(&v.Int32List), types.IntList(&v.Int64List), types.UintList(&v.UintList)}
}
