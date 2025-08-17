package primitive

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
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
		encoding.StringScanner[[]byte](&v.Bytes),  //  1 - bytes
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
func (v Primitive) StrValue() driver.Value {
	return v.Str
}
func (v Primitive) BytesValue() driver.Value {
	return string(v.Bytes)
}
func (v Primitive) BoolValue() driver.Value {
	return v.Bool
}
func (v Primitive) IntValue() driver.Value {
	return (int64)(v.Int)
}
func (v Primitive) Int8Value() driver.Value {
	return (int64)(v.Int8)
}
func (v Primitive) Int16Value() driver.Value {
	return (int64)(v.Int16)
}
func (v Primitive) Int32Value() driver.Value {
	return (int64)(v.Int32)
}
func (v Primitive) Int64Value() driver.Value {
	return v.Int64
}
func (v Primitive) UintValue() driver.Value {
	return (int64)(v.Uint)
}
func (v Primitive) Uint8Value() driver.Value {
	return (int64)(v.Uint8)
}
func (v Primitive) Uint16Value() driver.Value {
	return (int64)(v.Uint16)
}
func (v Primitive) Uint32Value() driver.Value {
	return (int64)(v.Uint32)
}
func (v Primitive) Uint64Value() driver.Value {
	return (int64)(v.Uint64)
}
func (v Primitive) F32Value() driver.Value {
	return (float64)(v.F32)
}
func (v Primitive) F64Value() driver.Value {
	return v.F64
}
func (v Primitive) TimeValue() driver.Value {
	return v.Time
}
func (v Primitive) ColumnStr() sequel.ColumnValuer[string] {
	return sequel.Column("str", v.Str, func(val string) driver.Value {
		return val
	})
}
func (v Primitive) ColumnBytes() sequel.ColumnValuer[[]byte] {
	return sequel.Column("bytes", v.Bytes, func(val []byte) driver.Value {
		return string(val)
	})
}
func (v Primitive) ColumnBool() sequel.ColumnValuer[bool] {
	return sequel.Column("bool", v.Bool, func(val bool) driver.Value {
		return val
	})
}
func (v Primitive) ColumnInt() sequel.ColumnValuer[int] {
	return sequel.Column("int", v.Int, func(val int) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt8() sequel.ColumnValuer[int8] {
	return sequel.Column("int_8", v.Int8, func(val int8) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt16() sequel.ColumnValuer[int16] {
	return sequel.Column("int_16", v.Int16, func(val int16) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt32() sequel.ColumnValuer[int32] {
	return sequel.Column("int_32", v.Int32, func(val int32) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt64() sequel.ColumnValuer[int64] {
	return sequel.Column("int_64", v.Int64, func(val int64) driver.Value {
		return val
	})
}
func (v Primitive) ColumnUint() sequel.ColumnValuer[uint] {
	return sequel.Column("uint", v.Uint, func(val uint) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint8() sequel.ColumnValuer[uint8] {
	return sequel.Column("uint_8", v.Uint8, func(val uint8) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint16() sequel.ColumnValuer[uint16] {
	return sequel.Column("uint_16", v.Uint16, func(val uint16) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint32() sequel.ColumnValuer[uint32] {
	return sequel.Column("uint_32", v.Uint32, func(val uint32) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint64() sequel.ColumnValuer[uint64] {
	return sequel.Column("uint_64", v.Uint64, func(val uint64) driver.Value {
		return (int64)(val)
	})
}
func (v Primitive) ColumnF32() sequel.ColumnValuer[float32] {
	return sequel.Column("f_32", v.F32, func(val float32) driver.Value {
		return (float64)(val)
	})
}
func (v Primitive) ColumnF64() sequel.ColumnValuer[float64] {
	return sequel.Column("f_64", v.F64, func(val float64) driver.Value {
		return val
	})
}
func (v Primitive) ColumnTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("time", v.Time, func(val time.Time) driver.Value {
		return val
	})
}
