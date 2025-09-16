package slice

import (
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
	return "id", 0, v.ID
}
func (Slice) Columns() []string {
	return []string{"id", "bool_list", "str_list", "custom_str_list", "int_list", "int_8_list", "int_16_list", "int_32_list", "int_64_list", "uint_list", "uint_8_list", "uint_16_list", "uint_32_list", "uint_64_list", "f_32_list", "f_64_list"} // 16
}
func (v Slice) Values() []any {
	return []any{
		(sqltype.BoolSlice[bool])(v.BoolList),             //  1 - bool_list
		(sqltype.StringSlice[string])(v.StrList),          //  2 - str_list
		(sqltype.StringSlice[customStr])(v.CustomStrList), //  3 - custom_str_list
		(sqltype.IntSlice[int])(v.IntList),                //  4 - int_list
		(sqltype.Int8Slice[int8])(v.Int8List),             //  5 - int_8_list
		(sqltype.Int16Slice[int16])(v.Int16List),          //  6 - int_16_list
		(sqltype.Int32Slice[int32])(v.Int32List),          //  7 - int_32_list
		(sqltype.Int64Slice[int64])(v.Int64List),          //  8 - int_64_list
		(sqltype.UintSlice[uint])(v.UintList),             //  9 - uint_list
		(sqltype.Uint8Slice[uint8])(v.Uint8List),          // 10 - uint_8_list
		(sqltype.Uint16Slice[uint16])(v.Uint16List),       // 11 - uint_16_list
		(sqltype.Uint32Slice[uint32])(v.Uint32List),       // 12 - uint_32_list
		(sqltype.Uint64Slice[uint64])(v.Uint64List),       // 13 - uint_64_list
		(sqltype.Float32Slice[float32])(v.F32List),        // 14 - f_32_list
		(sqltype.Float64Slice[float64])(v.F64List),        // 15 - f_64_list
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
	return "INSERT INTO `slice` (`bool_list`,`str_list`,`custom_str_list`,`int_list`,`int_8_list`,`int_16_list`,`int_32_list`,`int_64_list`,`uint_list`,`uint_8_list`,`uint_16_list`,`uint_32_list`,`uint_64_list`,`f_32_list`,`f_64_list`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{(sqltype.BoolSlice[bool])(v.BoolList), (sqltype.StringSlice[string])(v.StrList), (sqltype.StringSlice[customStr])(v.CustomStrList), (sqltype.IntSlice[int])(v.IntList), (sqltype.Int8Slice[int8])(v.Int8List), (sqltype.Int16Slice[int16])(v.Int16List), (sqltype.Int32Slice[int32])(v.Int32List), (sqltype.Int64Slice[int64])(v.Int64List), (sqltype.UintSlice[uint])(v.UintList), (sqltype.Uint8Slice[uint8])(v.Uint8List), (sqltype.Uint16Slice[uint16])(v.Uint16List), (sqltype.Uint32Slice[uint32])(v.Uint32List), (sqltype.Uint64Slice[uint64])(v.Uint64List), (sqltype.Float32Slice[float32])(v.F32List), (sqltype.Float64Slice[float64])(v.F64List)}
}
func (v Slice) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`bool_list`,`str_list`,`custom_str_list`,`int_list`,`int_8_list`,`int_16_list`,`int_32_list`,`int_64_list`,`uint_list`,`uint_8_list`,`uint_16_list`,`uint_32_list`,`uint_64_list`,`f_32_list`,`f_64_list` FROM `slice` WHERE `id` = ? LIMIT 1;", []any{v.ID}
}
func (v Slice) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `slice` SET `bool_list` = ?,`str_list` = ?,`custom_str_list` = ?,`int_list` = ?,`int_8_list` = ?,`int_16_list` = ?,`int_32_list` = ?,`int_64_list` = ?,`uint_list` = ?,`uint_8_list` = ?,`uint_16_list` = ?,`uint_32_list` = ?,`uint_64_list` = ?,`f_32_list` = ?,`f_64_list` = ? WHERE `id` = ?;", []any{(sqltype.BoolSlice[bool])(v.BoolList), (sqltype.StringSlice[string])(v.StrList), (sqltype.StringSlice[customStr])(v.CustomStrList), (sqltype.IntSlice[int])(v.IntList), (sqltype.Int8Slice[int8])(v.Int8List), (sqltype.Int16Slice[int16])(v.Int16List), (sqltype.Int32Slice[int32])(v.Int32List), (sqltype.Int64Slice[int64])(v.Int64List), (sqltype.UintSlice[uint])(v.UintList), (sqltype.Uint8Slice[uint8])(v.Uint8List), (sqltype.Uint16Slice[uint16])(v.Uint16List), (sqltype.Uint32Slice[uint32])(v.Uint32List), (sqltype.Uint64Slice[uint64])(v.Uint64List), (sqltype.Float32Slice[float32])(v.F32List), (sqltype.Float64Slice[float64])(v.F64List), v.ID}
}
func (v Slice) IDValue() any {
	return v.ID
}
func (v Slice) BoolListValue() any {
	return (sqltype.BoolSlice[bool])(v.BoolList)
}
func (v Slice) StrListValue() any {
	return (sqltype.StringSlice[string])(v.StrList)
}
func (v Slice) CustomStrListValue() any {
	return (sqltype.StringSlice[customStr])(v.CustomStrList)
}
func (v Slice) IntListValue() any {
	return (sqltype.IntSlice[int])(v.IntList)
}
func (v Slice) Int8ListValue() any {
	return (sqltype.Int8Slice[int8])(v.Int8List)
}
func (v Slice) Int16ListValue() any {
	return (sqltype.Int16Slice[int16])(v.Int16List)
}
func (v Slice) Int32ListValue() any {
	return (sqltype.Int32Slice[int32])(v.Int32List)
}
func (v Slice) Int64ListValue() any {
	return (sqltype.Int64Slice[int64])(v.Int64List)
}
func (v Slice) UintListValue() any {
	return (sqltype.UintSlice[uint])(v.UintList)
}
func (v Slice) Uint8ListValue() any {
	return (sqltype.Uint8Slice[uint8])(v.Uint8List)
}
func (v Slice) Uint16ListValue() any {
	return (sqltype.Uint16Slice[uint16])(v.Uint16List)
}
func (v Slice) Uint32ListValue() any {
	return (sqltype.Uint32Slice[uint32])(v.Uint32List)
}
func (v Slice) Uint64ListValue() any {
	return (sqltype.Uint64Slice[uint64])(v.Uint64List)
}
func (v Slice) F32ListValue() any {
	return (sqltype.Float32Slice[float32])(v.F32List)
}
func (v Slice) F64ListValue() any {
	return (sqltype.Float64Slice[float64])(v.F64List)
}
func (v Slice) ColumnID() sequel.ColumnConvertClause[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) any {
		return val
	})
}
func (v Slice) ColumnBoolList() sequel.ColumnConvertClause[[]bool] {
	return sequel.Column("bool_list", v.BoolList, func(val []bool) any {
		return (sqltype.BoolSlice[bool])(val)
	})
}
func (v Slice) ColumnStrList() sequel.ColumnConvertClause[[]string] {
	return sequel.Column("str_list", v.StrList, func(val []string) any {
		return (sqltype.StringSlice[string])(val)
	})
}
func (v Slice) ColumnCustomStrList() sequel.ColumnConvertClause[[]customStr] {
	return sequel.Column("custom_str_list", v.CustomStrList, func(val []customStr) any {
		return (sqltype.StringSlice[customStr])(val)
	})
}
func (v Slice) ColumnIntList() sequel.ColumnConvertClause[[]int] {
	return sequel.Column("int_list", v.IntList, func(val []int) any {
		return (sqltype.IntSlice[int])(val)
	})
}
func (v Slice) ColumnInt8List() sequel.ColumnConvertClause[[]int8] {
	return sequel.Column("int_8_list", v.Int8List, func(val []int8) any {
		return (sqltype.Int8Slice[int8])(val)
	})
}
func (v Slice) ColumnInt16List() sequel.ColumnConvertClause[[]int16] {
	return sequel.Column("int_16_list", v.Int16List, func(val []int16) any {
		return (sqltype.Int16Slice[int16])(val)
	})
}
func (v Slice) ColumnInt32List() sequel.ColumnConvertClause[[]int32] {
	return sequel.Column("int_32_list", v.Int32List, func(val []int32) any {
		return (sqltype.Int32Slice[int32])(val)
	})
}
func (v Slice) ColumnInt64List() sequel.ColumnConvertClause[[]int64] {
	return sequel.Column("int_64_list", v.Int64List, func(val []int64) any {
		return (sqltype.Int64Slice[int64])(val)
	})
}
func (v Slice) ColumnUintList() sequel.ColumnConvertClause[[]uint] {
	return sequel.Column("uint_list", v.UintList, func(val []uint) any {
		return (sqltype.UintSlice[uint])(val)
	})
}
func (v Slice) ColumnUint8List() sequel.ColumnConvertClause[[]uint8] {
	return sequel.Column("uint_8_list", v.Uint8List, func(val []uint8) any {
		return (sqltype.Uint8Slice[uint8])(val)
	})
}
func (v Slice) ColumnUint16List() sequel.ColumnConvertClause[[]uint16] {
	return sequel.Column("uint_16_list", v.Uint16List, func(val []uint16) any {
		return (sqltype.Uint16Slice[uint16])(val)
	})
}
func (v Slice) ColumnUint32List() sequel.ColumnConvertClause[[]uint32] {
	return sequel.Column("uint_32_list", v.Uint32List, func(val []uint32) any {
		return (sqltype.Uint32Slice[uint32])(val)
	})
}
func (v Slice) ColumnUint64List() sequel.ColumnConvertClause[[]uint64] {
	return sequel.Column("uint_64_list", v.Uint64List, func(val []uint64) any {
		return (sqltype.Uint64Slice[uint64])(val)
	})
}
func (v Slice) ColumnF32List() sequel.ColumnConvertClause[[]float32] {
	return sequel.Column("f_32_list", v.F32List, func(val []float32) any {
		return (sqltype.Float32Slice[float32])(val)
	})
}
func (v Slice) ColumnF64List() sequel.ColumnConvertClause[[]float64] {
	return sequel.Column("f_64_list", v.F64List, func(val []float64) any {
		return (sqltype.Float64Slice[float64])(val)
	})
}
