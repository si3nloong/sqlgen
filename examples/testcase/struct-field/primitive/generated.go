package primitive

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Primitive) TableName() string {
	return "primitive"
}
func (Primitive) Columns() []string {
	return []string{"str", "bytes", "bool", "int", "int_8", "int_16", "int_32", "int_64", "uint", "uint_8", "uint_16", "uint_32", "uint_64", "f_32", "f_64", "time"} // 16
}
func (v Primitive) Values() []any {
	return []any{
		v.Str,             //  0 - str
		string(v.Bytes),   //  1 - bytes
		v.Bool,            //  2 - bool
		(int64)(v.Int),    //  3 - int
		(int64)(v.Int8),   //  4 - int_8
		(int64)(v.Int16),  //  5 - int_16
		(int64)(v.Int32),  //  6 - int_32
		v.Int64,           //  7 - int_64
		(int64)(v.Uint),   //  8 - uint
		(int64)(v.Uint8),  //  9 - uint_8
		(int64)(v.Uint16), // 10 - uint_16
		(int64)(v.Uint32), // 11 - uint_32
		(int64)(v.Uint64), // 12 - uint_64
		(float64)(v.F32),  // 13 - f_32
		v.F64,             // 14 - f_64
		v.Time,            // 15 - time
	}
}
func (v *Primitive) Addrs() []any {
	return []any{
		&v.Str,                                    //  0 - str
		encoding.StringScanner[byte](&v.Bytes),    //  1 - bytes
		&v.Bool,                                   //  2 - bool
		encoding.IntScanner[int](&v.Int),          //  3 - int
		encoding.Int8Scanner[int8](&v.Int8),       //  4 - int_8
		encoding.Int16Scanner[int16](&v.Int16),    //  5 - int_16
		encoding.Int32Scanner[int32](&v.Int32),    //  6 - int_32
		&v.Int64,                                  //  7 - int_64
		encoding.UintScanner[uint](&v.Uint),       //  8 - uint
		encoding.Uint8Scanner[uint8](&v.Uint8),    //  9 - uint_8
		encoding.Uint16Scanner[uint16](&v.Uint16), // 10 - uint_16
		encoding.Uint32Scanner[uint32](&v.Uint32), // 11 - uint_32
		encoding.Uint64Scanner[uint64](&v.Uint64), // 12 - uint_64
		encoding.Float32Scanner[float32](&v.F32),  // 13 - f_32
		&v.F64,                                    // 14 - f_64
		&v.Time,                                   // 15 - time
	}
}
func (Primitive) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" // 16
}
func (v Primitive) InsertOneStmt() (string, []any) {
	return "INSERT INTO primitive (str,bytes,bool,int,int_8,int_16,int_32,int_64,uint,uint_8,uint_16,uint_32,uint_64,f_32,f_64,time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", v.Values()
}
func (v Primitive) GetStr() driver.Value {
	return v.Str
}
func (v Primitive) GetBytes() driver.Value {
	return string(v.Bytes)
}
func (v Primitive) GetBool() driver.Value {
	return v.Bool
}
func (v Primitive) GetInt() driver.Value {
	return (int64)(v.Int)
}
func (v Primitive) GetInt8() driver.Value {
	return (int64)(v.Int8)
}
func (v Primitive) GetInt16() driver.Value {
	return (int64)(v.Int16)
}
func (v Primitive) GetInt32() driver.Value {
	return (int64)(v.Int32)
}
func (v Primitive) GetInt64() driver.Value {
	return v.Int64
}
func (v Primitive) GetUint() driver.Value {
	return (int64)(v.Uint)
}
func (v Primitive) GetUint8() driver.Value {
	return (int64)(v.Uint8)
}
func (v Primitive) GetUint16() driver.Value {
	return (int64)(v.Uint16)
}
func (v Primitive) GetUint32() driver.Value {
	return (int64)(v.Uint32)
}
func (v Primitive) GetUint64() driver.Value {
	return (int64)(v.Uint64)
}
func (v Primitive) GetF32() driver.Value {
	return (float64)(v.F32)
}
func (v Primitive) GetF64() driver.Value {
	return v.F64
}
func (v Primitive) GetTime() driver.Value {
	return v.Time
}
