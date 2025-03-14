package slice

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/sqltype"
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
func (Slice) Columns() []string {
	return []string{"id", "bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"} // 16
}
func (v Slice) Values() []any {
	return []any{
		sqltype.BoolSlice[bool](v.BoolList),             //  1 - bool_list
		sqltype.StringSlice[string](v.StrList),          //  2 - str_list
		sqltype.StringSlice[customStr](v.CustomStrList), //  3 - custom_str_list
		sqltype.IntSlice[int](v.IntList),                //  4 - int_list
		sqltype.Int8Slice[int8](v.Int8List),             //  5 - int_8_list
		sqltype.Int16Slice[int16](v.Int16List),          //  6 - int_16_list
		sqltype.Int32Slice[int32](v.Int32List),          //  7 - int_32_list
		sqltype.Int64Slice[int64](v.Int64List),          //  8 - int_64_list
		sqltype.UintSlice[uint](v.UintList),             //  9 - uint_list
		sqltype.Uint8Slice[uint8](v.Uint8List),          // 10 - uint_8_list
		sqltype.Uint16Slice[uint16](v.Uint16List),       // 11 - uint_16_list
		sqltype.Uint32Slice[uint32](v.Uint32List),       // 12 - uint_32_list
		sqltype.Uint64Slice[uint64](v.Uint64List),       // 13 - uint_64_list
		(sqltype.Float32Slice[float32])(v.F32List),      // 14 - f_32_list
		(sqltype.Float64Slice[float64])(v.F64List),      // 15 - f_64_list
	}
}
func (v *Slice) Addrs() []any {
	return []any{
		encoding.Uint64Scanner[uint64](&v.ID),               //  0 - id
		(*sqltype.BoolSlice[bool])(&v.BoolList),             //  1 - bool_list
		(*sqltype.StringSlice[string])(&v.StrList),          //  2 - str_list
		(*sqltype.StringSlice[customStr])(&v.CustomStrList), //  3 - custom_str_list
		(*sqltype.IntSlice[int])(&v.IntList),                //  4 - int_list
		(*sqltype.Int8Slice[int8])(&v.Int8List),             //  5 - int_8_list
		(*sqltype.Int16Slice[int16])(&v.Int16List),          //  6 - int_16_list
		(*sqltype.Int32Slice[int32])(&v.Int32List),          //  7 - int_32_list
		(*sqltype.Int64Slice[int64])(&v.Int64List),          //  8 - int_64_list
		(*sqltype.UintSlice[uint])(&v.UintList),             //  9 - uint_list
		(*sqltype.Uint8Slice[uint8])(&v.Uint8List),          // 10 - uint_8_list
		(*sqltype.Uint16Slice[uint16])(&v.Uint16List),       // 11 - uint_16_list
		(*sqltype.Uint32Slice[uint32])(&v.Uint32List),       // 12 - uint_32_list
		(*sqltype.Uint64Slice[uint64])(&v.Uint64List),       // 13 - uint_64_list
		(*sqltype.Float32Slice[float32])(&v.F32List),        // 14 - f_32_list
		(*sqltype.Float64Slice[float64])(&v.F64List),        // 15 - f_64_list
	}
}
func (Slice) InsertColumns() []string {
	return []string{"bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"} // 15
}
func (Slice) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" // 15
}
func (v Slice) InsertOneStmt() (string, []any) {
	return "INSERT INTO slice (bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{sqltype.BoolSlice[bool](v.BoolList), sqltype.StringSlice[string](v.StrList), sqltype.StringSlice[customStr](v.CustomStrList), sqltype.IntSlice[int](v.IntList), sqltype.Int8Slice[int8](v.Int8List), sqltype.Int16Slice[int16](v.Int16List), sqltype.Int32Slice[int32](v.Int32List), sqltype.Int64Slice[int64](v.Int64List), sqltype.UintSlice[uint](v.UintList), sqltype.Uint8Slice[uint8](v.Uint8List), sqltype.Uint16Slice[uint16](v.Uint16List), sqltype.Uint32Slice[uint32](v.Uint32List), sqltype.Uint64Slice[uint64](v.Uint64List), (sqltype.Float32Slice[float32])(v.F32List), (sqltype.Float64Slice[float64])(v.F64List)}
}
func (v Slice) FindOneByPKStmt() (string, []any) {
	return "SELECT id,bool_list,str_list,custom_str_list,int_list,int_8_list,int_16_list,int_32_list,int_64_list,uint_list,uint_8_list,uint_16_list,uint_32_list,uint_64_list,f_32_list,f_64_list FROM slice WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Slice) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE slice SET bool_list = ?,str_list = ?,custom_str_list = ?,int_list = ?,int_8_list = ?,int_16_list = ?,int_32_list = ?,int_64_list = ?,uint_list = ?,uint_8_list = ?,uint_16_list = ?,uint_32_list = ?,uint_64_list = ?,f_32_list = ?,f_64_list = ? WHERE id = ?;", []any{sqltype.BoolSlice[bool](v.BoolList), sqltype.StringSlice[string](v.StrList), sqltype.StringSlice[customStr](v.CustomStrList), sqltype.IntSlice[int](v.IntList), sqltype.Int8Slice[int8](v.Int8List), sqltype.Int16Slice[int16](v.Int16List), sqltype.Int32Slice[int32](v.Int32List), sqltype.Int64Slice[int64](v.Int64List), sqltype.UintSlice[uint](v.UintList), sqltype.Uint8Slice[uint8](v.Uint8List), sqltype.Uint16Slice[uint16](v.Uint16List), sqltype.Uint32Slice[uint32](v.Uint32List), sqltype.Uint64Slice[uint64](v.Uint64List), (sqltype.Float32Slice[float32])(v.F32List), (sqltype.Float64Slice[float64])(v.F64List), (int64)(v.ID)}
}
func (v Slice) IDValue() driver.Value {
	return (int64)(v.ID)
}
func (v Slice) BoolListValue() driver.Value {
	return sqltype.BoolSlice[bool](v.BoolList)
}
func (v Slice) StrListValue() driver.Value {
	return sqltype.StringSlice[string](v.StrList)
}
func (v Slice) CustomStrListValue() driver.Value {
	return sqltype.StringSlice[customStr](v.CustomStrList)
}
func (v Slice) IntListValue() driver.Value {
	return sqltype.IntSlice[int](v.IntList)
}
func (v Slice) Int8ListValue() driver.Value {
	return sqltype.Int8Slice[int8](v.Int8List)
}
func (v Slice) Int16ListValue() driver.Value {
	return sqltype.Int16Slice[int16](v.Int16List)
}
func (v Slice) Int32ListValue() driver.Value {
	return sqltype.Int32Slice[int32](v.Int32List)
}
func (v Slice) Int64ListValue() driver.Value {
	return sqltype.Int64Slice[int64](v.Int64List)
}
func (v Slice) UintListValue() driver.Value {
	return sqltype.UintSlice[uint](v.UintList)
}
func (v Slice) Uint8ListValue() driver.Value {
	return sqltype.Uint8Slice[uint8](v.Uint8List)
}
func (v Slice) Uint16ListValue() driver.Value {
	return sqltype.Uint16Slice[uint16](v.Uint16List)
}
func (v Slice) Uint32ListValue() driver.Value {
	return sqltype.Uint32Slice[uint32](v.Uint32List)
}
func (v Slice) Uint64ListValue() driver.Value {
	return sqltype.Uint64Slice[uint64](v.Uint64List)
}
func (v Slice) F32ListValue() driver.Value {
	return (sqltype.Float32Slice[float32])(v.F32List)
}
func (v Slice) F64ListValue() driver.Value {
	return (sqltype.Float64Slice[float64])(v.F64List)
}
func (v Slice) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value {
		return (int64)(val)
	})
}
func (v Slice) GetBoolList() sequel.ColumnValuer[[]bool] {
	return sequel.Column("bool_list", v.BoolList, func(val []bool) driver.Value {
		return sqltype.BoolSlice[bool](val)
	})
}
func (v Slice) GetStrList() sequel.ColumnValuer[[]string] {
	return sequel.Column("str_list", v.StrList, func(val []string) driver.Value {
		return sqltype.StringSlice[string](val)
	})
}
func (v Slice) GetCustomStrList() sequel.ColumnValuer[[]customStr] {
	return sequel.Column("custom_str_list", v.CustomStrList, func(val []customStr) driver.Value {
		return sqltype.StringSlice[customStr](val)
	})
}
func (v Slice) GetIntList() sequel.ColumnValuer[[]int] {
	return sequel.Column("int_list", v.IntList, func(val []int) driver.Value {
		return sqltype.IntSlice[int](val)
	})
}
func (v Slice) GetInt8List() sequel.ColumnValuer[[]int8] {
	return sequel.Column("int_8_list", v.Int8List, func(val []int8) driver.Value {
		return sqltype.Int8Slice[int8](val)
	})
}
func (v Slice) GetInt16List() sequel.ColumnValuer[[]int16] {
	return sequel.Column("int_16_list", v.Int16List, func(val []int16) driver.Value {
		return sqltype.Int16Slice[int16](val)
	})
}
func (v Slice) GetInt32List() sequel.ColumnValuer[[]int32] {
	return sequel.Column("int_32_list", v.Int32List, func(val []int32) driver.Value {
		return sqltype.Int32Slice[int32](val)
	})
}
func (v Slice) GetInt64List() sequel.ColumnValuer[[]int64] {
	return sequel.Column("int_64_list", v.Int64List, func(val []int64) driver.Value {
		return sqltype.Int64Slice[int64](val)
	})
}
func (v Slice) GetUintList() sequel.ColumnValuer[[]uint] {
	return sequel.Column("uint_list", v.UintList, func(val []uint) driver.Value {
		return sqltype.UintSlice[uint](val)
	})
}
func (v Slice) GetUint8List() sequel.ColumnValuer[[]uint8] {
	return sequel.Column("uint_8_list", v.Uint8List, func(val []uint8) driver.Value {
		return sqltype.Uint8Slice[uint8](val)
	})
}
func (v Slice) GetUint16List() sequel.ColumnValuer[[]uint16] {
	return sequel.Column("uint_16_list", v.Uint16List, func(val []uint16) driver.Value {
		return sqltype.Uint16Slice[uint16](val)
	})
}
func (v Slice) GetUint32List() sequel.ColumnValuer[[]uint32] {
	return sequel.Column("uint_32_list", v.Uint32List, func(val []uint32) driver.Value {
		return sqltype.Uint32Slice[uint32](val)
	})
}
func (v Slice) GetUint64List() sequel.ColumnValuer[[]uint64] {
	return sequel.Column("uint_64_list", v.Uint64List, func(val []uint64) driver.Value {
		return sqltype.Uint64Slice[uint64](val)
	})
}
func (v Slice) GetF32List() sequel.ColumnValuer[[]float32] {
	return sequel.Column("f_32_list", v.F32List, func(val []float32) driver.Value {
		return (sqltype.Float32Slice[float32])(val)
	})
}
func (v Slice) GetF64List() sequel.ColumnValuer[[]float64] {
	return sequel.Column("f_64_list", v.F64List, func(val []float64) driver.Value {
		return (sqltype.Float64Slice[float64])(val)
	})
}
