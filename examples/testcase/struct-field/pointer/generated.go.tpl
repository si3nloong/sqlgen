package pointer

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Ptr) TableName() string {
	return "ptr"
}
func (Ptr) HasPK()      {}
func (Ptr) IsAutoIncr() {}
func (v Ptr) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (Ptr) Columns() []string {
	return []string{"id", "str", "bytes", "bool", "int", "int_8", "int_16", "int_32", "int_64", "uint", "uint_8", "uint_16", "uint_32", "uint_64", "f_32", "f_64", "time", "nested", "embeded_time"}
}
func (v Ptr) Values() []any {
	values := make([]any, 19)
	values[0] = (int64)(v.ID)
	if v.Str != nil {
		values[1] = (string)(*v.Str)
	}
	if v.Bytes != nil {
		values[2] = string(*v.Bytes)
	}
	if v.Bool != nil {
		values[3] = (bool)(*v.Bool)
	}
	if v.Int != nil {
		values[4] = (int64)(*v.Int)
	}
	if v.Int8 != nil {
		values[5] = (int64)(*v.Int8)
	}
	if v.Int16 != nil {
		values[6] = (int64)(*v.Int16)
	}
	if v.Int32 != nil {
		values[7] = (int64)(*v.Int32)
	}
	if v.Int64 != nil {
		values[8] = (int64)(*v.Int64)
	}
	if v.Uint != nil {
		values[9] = (int64)(*v.Uint)
	}
	if v.Uint8 != nil {
		values[10] = (int64)(*v.Uint8)
	}
	if v.Uint16 != nil {
		values[11] = (int64)(*v.Uint16)
	}
	if v.Uint32 != nil {
		values[12] = (int64)(*v.Uint32)
	}
	if v.Uint64 != nil {
		values[13] = (int64)(*v.Uint64)
	}
	if v.F32 != nil {
		values[14] = (float64)(*v.F32)
	}
	if v.F64 != nil {
		values[15] = (float64)(*v.F64)
	}
	if v.Time != nil {
		values[16] = (time.Time)(*v.Time)
	}
	if v.Nested != nil {
		values[17] = types.JSONMarshaler(*v.Nested)
	}
	if v.embeded.EmbededTime != nil {
		values[18] = (time.Time)(*v.embeded.EmbededTime)
	}
	return values
}
func (v *Ptr) Addrs() []any {
	addrs := make([]any, 19)
	addrs[0] = types.Integer(&v.ID)
	if v.Str == nil {
		v.Str = new(string)
	}
	addrs[1] = types.String(v.Str)
	if v.Bytes == nil {
		v.Bytes = new([]byte)
	}
	addrs[2] = types.String(v.Bytes)
	if v.Bool == nil {
		v.Bool = new(bool)
	}
	addrs[3] = types.Bool(v.Bool)
	if v.Int == nil {
		v.Int = new(int)
	}
	addrs[4] = types.Integer(v.Int)
	if v.Int8 == nil {
		v.Int8 = new(int8)
	}
	addrs[5] = types.Integer(v.Int8)
	if v.Int16 == nil {
		v.Int16 = new(int16)
	}
	addrs[6] = types.Integer(v.Int16)
	if v.Int32 == nil {
		v.Int32 = new(int32)
	}
	addrs[7] = types.Integer(v.Int32)
	if v.Int64 == nil {
		v.Int64 = new(int64)
	}
	addrs[8] = types.Integer(v.Int64)
	if v.Uint == nil {
		v.Uint = new(uint)
	}
	addrs[9] = types.Integer(v.Uint)
	if v.Uint8 == nil {
		v.Uint8 = new(uint8)
	}
	addrs[10] = types.Integer(v.Uint8)
	if v.Uint16 == nil {
		v.Uint16 = new(uint16)
	}
	addrs[11] = types.Integer(v.Uint16)
	if v.Uint32 == nil {
		v.Uint32 = new(uint32)
	}
	addrs[12] = types.Integer(v.Uint32)
	if v.Uint64 == nil {
		v.Uint64 = new(uint64)
	}
	addrs[13] = types.Integer(v.Uint64)
	if v.F32 == nil {
		v.F32 = new(float32)
	}
	addrs[14] = types.Float32(v.F32)
	if v.F64 == nil {
		v.F64 = new(float64)
	}
	addrs[15] = types.Float64(v.F64)
	if v.Time == nil {
		v.Time = new(time.Time)
	}
	addrs[16] = types.Time(v.Time)
	if v.Nested == nil {
		v.Nested = new(nested)
	}
	addrs[17] = types.JSONUnmarshaler(v.Nested)
	if v.embeded.EmbededTime == nil {
		v.embeded.EmbededTime = new(time.Time)
	}
	addrs[18] = types.Time(v.embeded.EmbededTime)
	return addrs
}
func (Ptr) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (v Ptr) InsertOneStmt() (string, []any) {
	return "INSERT INTO ptr (str,bytes,bool,int,int_8,int_16,int_32,int_64,uint,uint_8,uint_16,uint_32,uint_64,f_32,f_64,time,nested,embeded_time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{types.String(v.Str), types.String(v.Bytes), types.Bool(v.Bool), types.Integer(v.Int), types.Integer(v.Int8), types.Integer(v.Int16), types.Integer(v.Int32), types.Integer(v.Int64), types.Integer(v.Uint), types.Integer(v.Uint8), types.Integer(v.Uint16), types.Integer(v.Uint32), types.Integer(v.Uint64), types.Float32(v.F32), types.Float64(v.F64), types.Time(v.Time), types.JSONMarshaler(v.Nested), types.Time(v.embeded.EmbededTime)}
}
func (v Ptr) FindOneByPKStmt() (string, []any) {
	return "SELECT id,str,bytes,bool,int,int_8,int_16,int_32,int_64,uint,uint_8,uint_16,uint_32,uint_64,f_32,f_64,time,nested,embeded_time FROM ptr WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Ptr) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE ptr SET str = ?,bytes = ?,bool = ?,int = ?,int_8 = ?,int_16 = ?,int_32 = ?,int_64 = ?,uint = ?,uint_8 = ?,uint_16 = ?,uint_32 = ?,uint_64 = ?,f_32 = ?,f_64 = ?,time = ?,nested = ?,embeded_time = ? WHERE id = ?;", []any{types.String(v.Str), types.String(v.Bytes), types.Bool(v.Bool), types.Integer(v.Int), types.Integer(v.Int8), types.Integer(v.Int16), types.Integer(v.Int32), types.Integer(v.Int64), types.Integer(v.Uint), types.Integer(v.Uint8), types.Integer(v.Uint16), types.Integer(v.Uint32), types.Integer(v.Uint64), types.Float32(v.F32), types.Float64(v.F64), types.Time(v.Time), types.JSONMarshaler(v.Nested), types.Time(v.embeded.EmbededTime), (int64)(v.ID)}
}
func (v Ptr) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return (int64)(val) })
}
func (v Ptr) GetStr() sequel.ColumnValuer[*string] {
	return sequel.Column("str", v.Str, func(val *string) driver.Value { return types.String(val) })
}
func (v Ptr) GetBytes() sequel.ColumnValuer[*[]byte] {
	return sequel.Column("bytes", v.Bytes, func(val *[]byte) driver.Value { return types.String(val) })
}
func (v Ptr) GetBool() sequel.ColumnValuer[*bool] {
	return sequel.Column("bool", v.Bool, func(val *bool) driver.Value { return types.Bool(val) })
}
func (v Ptr) GetInt() sequel.ColumnValuer[*int] {
	return sequel.Column("int", v.Int, func(val *int) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetInt8() sequel.ColumnValuer[*int8] {
	return sequel.Column("int_8", v.Int8, func(val *int8) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetInt16() sequel.ColumnValuer[*int16] {
	return sequel.Column("int_16", v.Int16, func(val *int16) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetInt32() sequel.ColumnValuer[*int32] {
	return sequel.Column("int_32", v.Int32, func(val *int32) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetInt64() sequel.ColumnValuer[*int64] {
	return sequel.Column("int_64", v.Int64, func(val *int64) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetUint() sequel.ColumnValuer[*uint] {
	return sequel.Column("uint", v.Uint, func(val *uint) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetUint8() sequel.ColumnValuer[*uint8] {
	return sequel.Column("uint_8", v.Uint8, func(val *uint8) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetUint16() sequel.ColumnValuer[*uint16] {
	return sequel.Column("uint_16", v.Uint16, func(val *uint16) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetUint32() sequel.ColumnValuer[*uint32] {
	return sequel.Column("uint_32", v.Uint32, func(val *uint32) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetUint64() sequel.ColumnValuer[*uint64] {
	return sequel.Column("uint_64", v.Uint64, func(val *uint64) driver.Value { return types.Integer(val) })
}
func (v Ptr) GetF32() sequel.ColumnValuer[*float32] {
	return sequel.Column("f_32", v.F32, func(val *float32) driver.Value { return types.Float32(val) })
}
func (v Ptr) GetF64() sequel.ColumnValuer[*float64] {
	return sequel.Column("f_64", v.F64, func(val *float64) driver.Value { return types.Float64(val) })
}
func (v Ptr) GetTime() sequel.ColumnValuer[*time.Time] {
	return sequel.Column("time", v.Time, func(val *time.Time) driver.Value { return types.Time(val) })
}
func (v Ptr) GetNested() sequel.ColumnValuer[*nested] {
	return sequel.Column("nested", v.Nested, func(val *nested) driver.Value { return types.JSONMarshaler(val) })
}
func (v Ptr) GetEmbededTime() sequel.ColumnValuer[*time.Time] {
	return sequel.Column("embeded_time", v.embeded.EmbededTime, func(val *time.Time) driver.Value { return types.Time(val) })
}
