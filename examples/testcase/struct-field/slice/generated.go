package slice

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/sqltype"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Slice) TableName() string {
	return "slice"
}
func (Slice) HasPK()      {}
func (Slice) IsAutoIncr() {}
func (v *Slice) ScanAutoIncr(val int64) error {
	v.ID = uint64(val)
	return nil
}
func (v Slice) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (Slice) InsertColumns() []string {
	return []string{"bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"} // 15
}
func (Slice) Columns() []string {
	return []string{"id", "bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"} // 16
}
func (v Slice) Values() []any {
	return []any{
		encoding.MarshalBoolSlice(v.BoolList),           //  1 - bool_list
		sqltype.StringSlice[string](v.StrList),          //  2 - str_list
		sqltype.StringSlice[customStr](v.CustomStrList), //  3 - custom_str_list
		sqltype.IntSlice[int](v.IntList),                //  4 - int_list
		(sqltype.Int8Slice[int8])(v.Int8List),           //  5 - int_8_list
		(sqltype.Int16Slice[int16])(v.Int16List),        //  6 - int_16_list
		(sqltype.Int32Slice[int32])(v.Int32List),        //  7 - int_32_list
		(sqltype.Int64Slice[int64])(v.Int64List),        //  8 - int_64_list
		encoding.MarshalUintSlice(v.UintList),           //  9 - uint_list
		encoding.MarshalUintSlice(v.Uint8List),          // 10 - uint_8_list
		encoding.MarshalUintSlice(v.Uint16List),         // 11 - uint_16_list
		encoding.MarshalUintSlice(v.Uint32List),         // 12 - uint_32_list
		encoding.MarshalUintSlice(v.Uint64List),         // 13 - uint_64_list
		(sqltype.Float32Slice[float32])(v.F32List),      // 14 - f_32_list
		(sqltype.Float64Slice[float64])(v.F64List),      // 15 - f_64_list
	}
}
func (v *Slice) Addrs() []any {
	return []any{
		encoding.Uint64Scanner[uint64](&v.ID),               //  0 - id
		types.BoolSlice(&v.BoolList),                        //  1 - bool_list
		(*sqltype.StringSlice[string])(&v.StrList),          //  2 - str_list
		(*sqltype.StringSlice[customStr])(&v.CustomStrList), //  3 - custom_str_list
		(*sqltype.IntSlice[int])(&v.IntList),                //  4 - int_list
		(*sqltype.Int8Slice[int8])(&v.Int8List),             //  5 - int_8_list
		(*sqltype.Int16Slice[int16])(&v.Int16List),          //  6 - int_16_list
		(*sqltype.Int32Slice[int32])(&v.Int32List),          //  7 - int_32_list
		(*sqltype.Int64Slice[int64])(&v.Int64List),          //  8 - int_64_list
		types.UintSlice(&v.UintList),                        //  9 - uint_list
		types.UintSlice(&v.Uint8List),                       // 10 - uint_8_list
		types.UintSlice(&v.Uint16List),                      // 11 - uint_16_list
		types.UintSlice(&v.Uint32List),                      // 12 - uint_32_list
		types.UintSlice(&v.Uint64List),                      // 13 - uint_64_list
		(*sqltype.Float32Slice[float32])(&v.F32List),        // 14 - f_32_list
		(*sqltype.Float64Slice[float64])(&v.F64List),        // 15 - f_64_list
	}
}
func (Slice) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" // 15
}
func (v Slice) InsertOneStmt() (string, []any) {
	return "INSERT INTO slice (bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{encoding.MarshalBoolSlice(v.BoolList), sqltype.StringSlice[string](v.StrList), sqltype.StringSlice[customStr](v.CustomStrList), sqltype.IntSlice[int](v.IntList), (sqltype.Int8Slice[int8])(v.Int8List), (sqltype.Int16Slice[int16])(v.Int16List), (sqltype.Int32Slice[int32])(v.Int32List), (sqltype.Int64Slice[int64])(v.Int64List), encoding.MarshalUintSlice(v.UintList), encoding.MarshalUintSlice(v.Uint8List), encoding.MarshalUintSlice(v.Uint16List), encoding.MarshalUintSlice(v.Uint32List), encoding.MarshalUintSlice(v.Uint64List), (sqltype.Float32Slice[float32])(v.F32List), (sqltype.Float64Slice[float64])(v.F64List)}
}
func (v Slice) FindOneByPKStmt() (string, []any) {
	return "SELECT id,bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list FROM slice WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Slice) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE slice SET bool_list = ?,str_list = ?,custom_str_list = ?,int_list = ?,int_8_list = ?,int_16_list = ?,int_32_list = ?,int_64_list = ?,uint_list = ?,uint_8_list = ?,uint_16_list = ?,uint_32_list = ?,uint_64_list = ?,f_32_list = ?,f_64_list = ? WHERE id = ?;", []any{encoding.MarshalBoolSlice(v.BoolList), sqltype.StringSlice[string](v.StrList), sqltype.StringSlice[customStr](v.CustomStrList), sqltype.IntSlice[int](v.IntList), (sqltype.Int8Slice[int8])(v.Int8List), (sqltype.Int16Slice[int16])(v.Int16List), (sqltype.Int32Slice[int32])(v.Int32List), (sqltype.Int64Slice[int64])(v.Int64List), encoding.MarshalUintSlice(v.UintList), encoding.MarshalUintSlice(v.Uint8List), encoding.MarshalUintSlice(v.Uint16List), encoding.MarshalUintSlice(v.Uint32List), encoding.MarshalUintSlice(v.Uint64List), (sqltype.Float32Slice[float32])(v.F32List), (sqltype.Float64Slice[float64])(v.F64List), (int64)(v.ID)}
}
func (v Slice) GetID() driver.Value {
	return (int64)(v.ID)
}
func (v Slice) GetBoolList() driver.Value {
	return encoding.MarshalBoolSlice(v.BoolList)
}
func (v Slice) GetStrList() driver.Value {
	return sqltype.StringSlice[string](v.StrList)
}
func (v Slice) GetCustomStrList() driver.Value {
	return sqltype.StringSlice[customStr](v.CustomStrList)
}
func (v Slice) GetIntList() driver.Value {
	return sqltype.IntSlice[int](v.IntList)
}
func (v Slice) GetInt8List() driver.Value {
	return (sqltype.Int8Slice[int8])(v.Int8List)
}
func (v Slice) GetInt16List() driver.Value {
	return (sqltype.Int16Slice[int16])(v.Int16List)
}
func (v Slice) GetInt32List() driver.Value {
	return (sqltype.Int32Slice[int32])(v.Int32List)
}
func (v Slice) GetInt64List() driver.Value {
	return (sqltype.Int64Slice[int64])(v.Int64List)
}
func (v Slice) GetUintList() driver.Value {
	return encoding.MarshalUintSlice(v.UintList)
}
func (v Slice) GetUint8List() driver.Value {
	return encoding.MarshalUintSlice(v.Uint8List)
}
func (v Slice) GetUint16List() driver.Value {
	return encoding.MarshalUintSlice(v.Uint16List)
}
func (v Slice) GetUint32List() driver.Value {
	return encoding.MarshalUintSlice(v.Uint32List)
}
func (v Slice) GetUint64List() driver.Value {
	return encoding.MarshalUintSlice(v.Uint64List)
}
func (v Slice) GetF32List() driver.Value {
	return (sqltype.Float32Slice[float32])(v.F32List)
}
func (v Slice) GetF64List() driver.Value {
	return (sqltype.Float64Slice[float64])(v.F64List)
}
