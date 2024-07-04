package slice

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Slice) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"`id`"},
			Definition: "PRIMARY KEY (`id`)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT"},
			{Name: "`bool_list`", Definition: "`bool_list` JSON NOT NULL"},
			{Name: "`str_list`", Definition: "`str_list` JSON NOT NULL"},
			{Name: "`custom_str_list`", Definition: "`custom_str_list` JSON NOT NULL"},
			{Name: "`int_list`", Definition: "`int_list` JSON NOT NULL"},
			{Name: "`int_8_list`", Definition: "`int_8_list` JSON NOT NULL"},
			{Name: "`int_16_list`", Definition: "`int_16_list` JSON NOT NULL"},
			{Name: "`int_32_list`", Definition: "`int_32_list` JSON NOT NULL"},
			{Name: "`int_64_list`", Definition: "`int_64_list` JSON NOT NULL"},
			{Name: "`uint_list`", Definition: "`uint_list` JSON NOT NULL"},
			{Name: "`uint_8_list`", Definition: "`uint_8_list` JSON NOT NULL"},
			{Name: "`uint_16_list`", Definition: "`uint_16_list` JSON NOT NULL"},
			{Name: "`uint_32_list`", Definition: "`uint_32_list` JSON NOT NULL"},
			{Name: "`uint_64_list`", Definition: "`uint_64_list` JSON NOT NULL"},
			{Name: "`f_32_list`", Definition: "`f_32_list` JSON NOT NULL"},
			{Name: "`f_64_list`", Definition: "`f_64_list` JSON NOT NULL"},
		},
	}
}
func (Slice) TableName() string {
	return "`slice`"
}
func (Slice) HasPK()      {}
func (Slice) IsAutoIncr() {}
func (v Slice) PK() (string, int, any) {
	return "`id`", 0, int64(v.ID)
}
func (Slice) ColumnNames() []string {
	return []string{"`id`", "`bool_list`", "`str_list`", "`custom_str_list`", "`int_list`", "`int_8_list`", "`int_16_list`", "`int_32_list`", "`int_64_list`", "`uint_list`", "`uint_8_list`", "`uint_16_list`", "`uint_32_list`", "`uint_64_list`", "`f_32_list`", "`f_64_list`"}
}
func (v Slice) Values() []any {
	return []any{int64(v.ID), encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringList(v.StrList), encoding.MarshalStringList(v.CustomStrList), encoding.MarshalSignedIntList(v.IntList), encoding.MarshalSignedIntList(v.Int8List), encoding.MarshalSignedIntList(v.Int16List), encoding.MarshalSignedIntList(v.Int32List), encoding.MarshalSignedIntList(v.Int64List), encoding.MarshalUnsignedIntList(v.UintList), encoding.MarshalUnsignedIntList(v.Uint8List), encoding.MarshalUnsignedIntList(v.Uint16List), encoding.MarshalUnsignedIntList(v.Uint32List), encoding.MarshalUnsignedIntList(v.Uint64List), encoding.MarshalFloatList(v.F32List), encoding.MarshalFloatList(v.F64List)}
}
func (v *Slice) Addrs() []any {
	return []any{types.Integer(&v.ID), types.BoolList(&v.BoolList), types.StringList(&v.StrList), types.StringList(&v.CustomStrList), types.IntList(&v.IntList), types.IntList(&v.Int8List), types.IntList(&v.Int16List), types.IntList(&v.Int32List), types.IntList(&v.Int64List), types.UintList(&v.UintList), types.UintList(&v.Uint8List), types.UintList(&v.Uint16List), types.UintList(&v.Uint32List), types.UintList(&v.Uint64List), types.FloatList(&v.F32List), types.FloatList(&v.F64List)}
}
func (Slice) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (v Slice) InsertOneStmt() (string, []any) {
	return "INSERT INTO `slice` (`bool_list`,`str_list`,`custom_str_list`,`int_list`,`int_8_list`,`int_16_list`,`int_32_list`,`int_64_list`,`uint_list`,`uint_8_list`,`uint_16_list`,`uint_32_list`,`uint_64_list`,`f_32_list`,`f_64_list`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringList(v.StrList), encoding.MarshalStringList(v.CustomStrList), encoding.MarshalSignedIntList(v.IntList), encoding.MarshalSignedIntList(v.Int8List), encoding.MarshalSignedIntList(v.Int16List), encoding.MarshalSignedIntList(v.Int32List), encoding.MarshalSignedIntList(v.Int64List), encoding.MarshalUnsignedIntList(v.UintList), encoding.MarshalUnsignedIntList(v.Uint8List), encoding.MarshalUnsignedIntList(v.Uint16List), encoding.MarshalUnsignedIntList(v.Uint32List), encoding.MarshalUnsignedIntList(v.Uint64List), encoding.MarshalFloatList(v.F32List), encoding.MarshalFloatList(v.F64List)}
}
func (v Slice) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`bool_list`,`str_list`,`custom_str_list`,`int_list`,`int_8_list`,`int_16_list`,`int_32_list`,`int_64_list`,`uint_list`,`uint_8_list`,`uint_16_list`,`uint_32_list`,`uint_64_list`,`f_32_list`,`f_64_list` FROM `slice` WHERE `id` = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v Slice) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `slice` SET `bool_list` = ?,`str_list` = ?,`custom_str_list` = ?,`int_list` = ?,`int_8_list` = ?,`int_16_list` = ?,`int_32_list` = ?,`int_64_list` = ?,`uint_list` = ?,`uint_8_list` = ?,`uint_16_list` = ?,`uint_32_list` = ?,`uint_64_list` = ?,`f_32_list` = ?,`f_64_list` = ? WHERE `id` = ? LIMIT 1;", []any{encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringList(v.StrList), encoding.MarshalStringList(v.CustomStrList), encoding.MarshalSignedIntList(v.IntList), encoding.MarshalSignedIntList(v.Int8List), encoding.MarshalSignedIntList(v.Int16List), encoding.MarshalSignedIntList(v.Int32List), encoding.MarshalSignedIntList(v.Int64List), encoding.MarshalUnsignedIntList(v.UintList), encoding.MarshalUnsignedIntList(v.Uint8List), encoding.MarshalUnsignedIntList(v.Uint16List), encoding.MarshalUnsignedIntList(v.Uint32List), encoding.MarshalUnsignedIntList(v.Uint64List), encoding.MarshalFloatList(v.F32List), encoding.MarshalFloatList(v.F64List), int64(v.ID)}
}
func (v Slice) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("`id`", v.ID, func(val uint64) driver.Value { return int64(val) })
}
func (v Slice) GetBoolList() sequel.ColumnValuer[[]bool] {
	return sequel.Column("`bool_list`", v.BoolList, func(val []bool) driver.Value { return encoding.MarshalBoolList(val) })
}
func (v Slice) GetStrList() sequel.ColumnValuer[[]string] {
	return sequel.Column("`str_list`", v.StrList, func(val []string) driver.Value { return encoding.MarshalStringList(val) })
}
func (v Slice) GetCustomStrList() sequel.ColumnValuer[[]customStr] {
	return sequel.Column("`custom_str_list`", v.CustomStrList, func(val []customStr) driver.Value { return encoding.MarshalStringList(val) })
}
func (v Slice) GetIntList() sequel.ColumnValuer[[]int] {
	return sequel.Column("`int_list`", v.IntList, func(val []int) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Slice) GetInt8List() sequel.ColumnValuer[[]int8] {
	return sequel.Column("`int_8_list`", v.Int8List, func(val []int8) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Slice) GetInt16List() sequel.ColumnValuer[[]int16] {
	return sequel.Column("`int_16_list`", v.Int16List, func(val []int16) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Slice) GetInt32List() sequel.ColumnValuer[[]int32] {
	return sequel.Column("`int_32_list`", v.Int32List, func(val []int32) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Slice) GetInt64List() sequel.ColumnValuer[[]int64] {
	return sequel.Column("`int_64_list`", v.Int64List, func(val []int64) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Slice) GetUintList() sequel.ColumnValuer[[]uint] {
	return sequel.Column("`uint_list`", v.UintList, func(val []uint) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Slice) GetUint8List() sequel.ColumnValuer[[]uint8] {
	return sequel.Column("`uint_8_list`", v.Uint8List, func(val []uint8) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Slice) GetUint16List() sequel.ColumnValuer[[]uint16] {
	return sequel.Column("`uint_16_list`", v.Uint16List, func(val []uint16) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Slice) GetUint32List() sequel.ColumnValuer[[]uint32] {
	return sequel.Column("`uint_32_list`", v.Uint32List, func(val []uint32) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Slice) GetUint64List() sequel.ColumnValuer[[]uint64] {
	return sequel.Column("`uint_64_list`", v.Uint64List, func(val []uint64) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Slice) GetF32List() sequel.ColumnValuer[[]float32] {
	return sequel.Column("`f_32_list`", v.F32List, func(val []float32) driver.Value { return encoding.MarshalFloatList(val) })
}
func (v Slice) GetF64List() sequel.ColumnValuer[[]float64] {
	return sequel.Column("`f_64_list`", v.F64List, func(val []float64) driver.Value { return encoding.MarshalFloatList(val) })
}
