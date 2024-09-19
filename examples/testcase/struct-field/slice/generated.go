package slice

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Slice) TableName() string {
	return "slice"
}
func (Slice) HasPK()      {}
func (Slice) IsAutoIncr() {}
func (v Slice) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (Slice) Columns() []string {
	return []string{"id", "bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"}
}
func (v Slice) Values() []any {
	return []any{(int64)(v.ID), encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringSlice(v.StrList), encoding.MarshalStringSlice(v.CustomStrList), encoding.MarshalIntSlice(v.IntList), encoding.MarshalIntSlice(v.Int8List), encoding.MarshalIntSlice(v.Int16List), encoding.MarshalIntSlice(v.Int32List), encoding.MarshalIntSlice(v.Int64List), encoding.MarshalUintSlice(v.UintList), encoding.MarshalUintSlice(v.Uint8List), encoding.MarshalUintSlice(v.Uint16List), encoding.MarshalUintSlice(v.Uint32List), encoding.MarshalUintSlice(v.Uint64List), encoding.MarshalFloatList(v.F32List, -1), encoding.MarshalFloatList(v.F64List, -1)}
}
func (v *Slice) Addrs() []any {
	return []any{types.Integer(&v.ID), types.BoolList(&v.BoolList), types.StringList(&v.StrList), types.StringList(&v.CustomStrList), types.IntList(&v.IntList), types.IntList(&v.Int8List), types.IntList(&v.Int16List), types.IntList(&v.Int32List), types.IntList(&v.Int64List), types.UintList(&v.UintList), types.UintList(&v.Uint8List), types.UintList(&v.Uint16List), types.UintList(&v.Uint32List), types.UintList(&v.Uint64List), types.FloatList(&v.F32List), types.FloatList(&v.F64List)}
}
func (Slice) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (v Slice) InsertOneStmt() (string, []any) {
	return "INSERT INTO slice (bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringSlice(v.StrList), encoding.MarshalStringSlice(v.CustomStrList), encoding.MarshalIntSlice(v.IntList), encoding.MarshalIntSlice(v.Int8List), encoding.MarshalIntSlice(v.Int16List), encoding.MarshalIntSlice(v.Int32List), encoding.MarshalIntSlice(v.Int64List), encoding.MarshalUintSlice(v.UintList), encoding.MarshalUintSlice(v.Uint8List), encoding.MarshalUintSlice(v.Uint16List), encoding.MarshalUintSlice(v.Uint32List), encoding.MarshalUintSlice(v.Uint64List), encoding.MarshalFloatList(v.F32List, -1), encoding.MarshalFloatList(v.F64List, -1)}
}
func (v Slice) FindOneByPKStmt() (string, []any) {
	return "SELECT id,bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list FROM slice WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Slice) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE slice SET bool_list = ?,str_list = ?,custom_str_list = ?,int_list = ?,int_8_list = ?,int_16_list = ?,int_32_list = ?,int_64_list = ?,uint_list = ?,uint_8_list = ?,uint_16_list = ?,uint_32_list = ?,uint_64_list = ?,f_32_list = ?,f_64_list = ? WHERE id = ?;", []any{encoding.MarshalBoolList(v.BoolList), encoding.MarshalStringSlice(v.StrList), encoding.MarshalStringSlice(v.CustomStrList), encoding.MarshalIntSlice(v.IntList), encoding.MarshalIntSlice(v.Int8List), encoding.MarshalIntSlice(v.Int16List), encoding.MarshalIntSlice(v.Int32List), encoding.MarshalIntSlice(v.Int64List), encoding.MarshalUintSlice(v.UintList), encoding.MarshalUintSlice(v.Uint8List), encoding.MarshalUintSlice(v.Uint16List), encoding.MarshalUintSlice(v.Uint32List), encoding.MarshalUintSlice(v.Uint64List), encoding.MarshalFloatList(v.F32List, -1), encoding.MarshalFloatList(v.F64List, -1), (int64)(v.ID)}
}
func (v Slice) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value { return (int64)(val) })
}
func (v Slice) GetBoolList() sequel.ColumnValuer[[]bool] {
	return sequel.Column("bool_list", v.BoolList, func(val []bool) driver.Value { return encoding.MarshalBoolList(val) })
}
func (v Slice) GetStrList() sequel.ColumnValuer[[]string] {
	return sequel.Column("str_list", v.StrList, func(val []string) driver.Value { return encoding.MarshalStringSlice(val) })
}
func (v Slice) GetCustomStrList() sequel.ColumnValuer[[]customStr] {
	return sequel.Column("custom_str_list", v.CustomStrList, func(val []customStr) driver.Value { return encoding.MarshalStringSlice(val) })
}
func (v Slice) GetIntList() sequel.ColumnValuer[[]int] {
	return sequel.Column("int_list", v.IntList, func(val []int) driver.Value { return encoding.MarshalIntSlice(val) })
}
func (v Slice) GetInt8List() sequel.ColumnValuer[[]int8] {
	return sequel.Column("int_8_list", v.Int8List, func(val []int8) driver.Value { return encoding.MarshalIntSlice(val) })
}
func (v Slice) GetInt16List() sequel.ColumnValuer[[]int16] {
	return sequel.Column("int_16_list", v.Int16List, func(val []int16) driver.Value { return encoding.MarshalIntSlice(val) })
}
func (v Slice) GetInt32List() sequel.ColumnValuer[[]int32] {
	return sequel.Column("int_32_list", v.Int32List, func(val []int32) driver.Value { return encoding.MarshalIntSlice(val) })
}
func (v Slice) GetInt64List() sequel.ColumnValuer[[]int64] {
	return sequel.Column("int_64_list", v.Int64List, func(val []int64) driver.Value { return encoding.MarshalIntSlice(val) })
}
func (v Slice) GetUintList() sequel.ColumnValuer[[]uint] {
	return sequel.Column("uint_list", v.UintList, func(val []uint) driver.Value { return encoding.MarshalUintSlice(val) })
}
func (v Slice) GetUint8List() sequel.ColumnValuer[[]uint8] {
	return sequel.Column("uint_8_list", v.Uint8List, func(val []uint8) driver.Value { return encoding.MarshalUintSlice(val) })
}
func (v Slice) GetUint16List() sequel.ColumnValuer[[]uint16] {
	return sequel.Column("uint_16_list", v.Uint16List, func(val []uint16) driver.Value { return encoding.MarshalUintSlice(val) })
}
func (v Slice) GetUint32List() sequel.ColumnValuer[[]uint32] {
	return sequel.Column("uint_32_list", v.Uint32List, func(val []uint32) driver.Value { return encoding.MarshalUintSlice(val) })
}
func (v Slice) GetUint64List() sequel.ColumnValuer[[]uint64] {
	return sequel.Column("uint_64_list", v.Uint64List, func(val []uint64) driver.Value { return encoding.MarshalUintSlice(val) })
}
func (v Slice) GetF32List() sequel.ColumnValuer[[]float32] {
	return sequel.Column("f_32_list", v.F32List, func(val []float32) driver.Value { return encoding.MarshalFloatList(val, -1) })
}
func (v Slice) GetF64List() sequel.ColumnValuer[[]float64] {
	return sequel.Column("f_64_list", v.F64List, func(val []float64) driver.Value { return encoding.MarshalFloatList(val, -1) })
}
