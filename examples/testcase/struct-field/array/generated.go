package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Array) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,`bool_list` JSON NOT NULL,`str_list` JSON NOT NULL,`custom_str_list` JSON NOT NULL,`int_list` JSON NOT NULL,`int_8_list` JSON NOT NULL,`int_16_list` JSON NOT NULL,`int_32_list` JSON NOT NULL,`int_64_list` JSON NOT NULL,`uint_list` JSON NOT NULL,`uint_8_list` JSON NOT NULL,`uint_16_list` JSON NOT NULL,`uint_32_list` JSON NOT NULL,`uint_64_list` JSON NOT NULL,`f_32_list` JSON NOT NULL,`f_64_list` JSON NOT NULL,PRIMARY KEY (`id`));"
}
func (Array) TableName() string {
	return "array"
}
func (Array) InsertOneStmt() string {
	return "INSERT INTO array (bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
}
func (Array) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (Array) Columns() []string {
	return []string{"id", "bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"}
}
func (Array) IsAutoIncr() {}
func (v Array) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, int64(v.ID)
}
func (v Array) FindByPKStmt() string {
	return "SELECT id,bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list FROM array WHERE id = ? LIMIT 1;"
}
func (Array) UpdateByPKStmt() string {
	return "UPDATE array SET bool_list = ?,str_list = ?,custom_str_list = ?,int_list = ?,int_8_list = ?,int_16_list = ?,int_32_list = ?,int_64_list = ?,uint_list = ?,uint_8_list = ?,uint_16_list = ?,uint_32_list = ?,uint_64_list = ?,f_32_list = ?,f_64_list = ? WHERE id = ? LIMIT 1;"
}
func (v Array) Values() []any {
	return []any{int64(v.ID), encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringList(v.StrList), encoding.MarshalStringList(v.CustomStrList), encoding.MarshalSignedIntList(v.IntList), encoding.MarshalSignedIntList(v.Int8List), encoding.MarshalSignedIntList(v.Int16List), encoding.MarshalSignedIntList(v.Int32List), encoding.MarshalSignedIntList(v.Int64List), encoding.MarshalUnsignedIntList(v.UintList), encoding.MarshalUnsignedIntList(v.Uint8List), encoding.MarshalUnsignedIntList(v.Uint16List), encoding.MarshalUnsignedIntList(v.Uint32List), encoding.MarshalUnsignedIntList(v.Uint64List), encoding.MarshalFloatList(v.F32List), encoding.MarshalFloatList(v.F64List)}
}
func (v *Array) Addrs() []any {
	return []any{types.Integer(&v.ID), types.BoolList(&v.BoolList), types.StringList(&v.StrList), types.StringList(&v.CustomStrList), types.IntList(&v.IntList), types.IntList(&v.Int8List), types.IntList(&v.Int16List), types.IntList(&v.Int32List), types.IntList(&v.Int64List), types.UintList(&v.UintList), types.UintList(&v.Uint8List), types.UintList(&v.Uint16List), types.UintList(&v.Uint32List), types.UintList(&v.Uint64List), types.FloatList(&v.F32List), types.FloatList(&v.F64List)}
}
func (v Array) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value { return int64(val) })
}
func (v Array) GetBoolList() sequel.ColumnValuer[[]bool] {
	return sequel.Column("bool_list", v.BoolList, func(val []bool) driver.Value { return encoding.MarshalBoolList(val) })
}
func (v Array) GetStrList() sequel.ColumnValuer[[]string] {
	return sequel.Column("str_list", v.StrList, func(val []string) driver.Value { return encoding.MarshalStringList(val) })
}
func (v Array) GetCustomStrList() sequel.ColumnValuer[[]customStr] {
	return sequel.Column("custom_str_list", v.CustomStrList, func(val []customStr) driver.Value { return encoding.MarshalStringList(val) })
}
func (v Array) GetIntList() sequel.ColumnValuer[[]int] {
	return sequel.Column("int_list", v.IntList, func(val []int) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Array) GetInt8List() sequel.ColumnValuer[[]int8] {
	return sequel.Column("int_8_list", v.Int8List, func(val []int8) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Array) GetInt16List() sequel.ColumnValuer[[]int16] {
	return sequel.Column("int_16_list", v.Int16List, func(val []int16) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Array) GetInt32List() sequel.ColumnValuer[[]int32] {
	return sequel.Column("int_32_list", v.Int32List, func(val []int32) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Array) GetInt64List() sequel.ColumnValuer[[]int64] {
	return sequel.Column("int_64_list", v.Int64List, func(val []int64) driver.Value { return encoding.MarshalSignedIntList(val) })
}
func (v Array) GetUintList() sequel.ColumnValuer[[]uint] {
	return sequel.Column("uint_list", v.UintList, func(val []uint) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Array) GetUint8List() sequel.ColumnValuer[[]uint8] {
	return sequel.Column("uint_8_list", v.Uint8List, func(val []uint8) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Array) GetUint16List() sequel.ColumnValuer[[]uint16] {
	return sequel.Column("uint_16_list", v.Uint16List, func(val []uint16) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Array) GetUint32List() sequel.ColumnValuer[[]uint32] {
	return sequel.Column("uint_32_list", v.Uint32List, func(val []uint32) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Array) GetUint64List() sequel.ColumnValuer[[]uint64] {
	return sequel.Column("uint_64_list", v.Uint64List, func(val []uint64) driver.Value { return encoding.MarshalUnsignedIntList(val) })
}
func (v Array) GetF32List() sequel.ColumnValuer[[]float32] {
	return sequel.Column("f_32_list", v.F32List, func(val []float32) driver.Value { return encoding.MarshalFloatList(val) })
}
func (v Array) GetF64List() sequel.ColumnValuer[[]float64] {
	return sequel.Column("f_64_list", v.F64List, func(val []float64) driver.Value { return encoding.MarshalFloatList(val) })
}
