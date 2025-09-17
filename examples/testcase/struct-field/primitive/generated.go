package primitive

import (
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
		v.Bytes,           //  1 - bytes
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
		v.Uint64,          // 12 - uint_64
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
	return "INSERT INTO `primitive` (`str`,`bytes`,`bool`,`int`,`int_8`,`int_16`,`int_32`,`int_64`,`uint`,`uint_8`,`uint_16`,`uint_32`,`uint_64`,`f_32`,`f_64`,`time`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", v.Values()
}
func (v Primitive) StrValue() any {
	return v.Str
}
func (v Primitive) BytesValue() any {
	return v.Bytes
}
func (v Primitive) BoolValue() any {
	return v.Bool
}
func (v Primitive) IntValue() any {
	return (int64)(v.Int)
}
func (v Primitive) Int8Value() any {
	return (int64)(v.Int8)
}
func (v Primitive) Int16Value() any {
	return (int64)(v.Int16)
}
func (v Primitive) Int32Value() any {
	return (int64)(v.Int32)
}
func (v Primitive) Int64Value() any {
	return v.Int64
}
func (v Primitive) UintValue() any {
	return (int64)(v.Uint)
}
func (v Primitive) Uint8Value() any {
	return (int64)(v.Uint8)
}
func (v Primitive) Uint16Value() any {
	return (int64)(v.Uint16)
}
func (v Primitive) Uint32Value() any {
	return (int64)(v.Uint32)
}
func (v Primitive) Uint64Value() any {
	return v.Uint64
}
func (v Primitive) F32Value() any {
	return (float64)(v.F32)
}
func (v Primitive) F64Value() any {
	return v.F64
}
func (v Primitive) TimeValue() any {
	return v.Time
}
func (v Primitive) ColumnStr() sequel.ColumnClause[string] {
	return sequel.BasicColumn("str", v.Str)
}
func (v Primitive) ColumnBytes() sequel.ColumnConvertClause[[]byte] {
	return sequel.Column("bytes", v.Bytes, func(val []byte) any {
		return val
	})
}
func (v Primitive) ColumnBool() sequel.ColumnClause[bool] {
	return sequel.BasicColumn("bool", v.Bool)
}
func (v Primitive) ColumnInt() sequel.ColumnConvertClause[int] {
	return sequel.Column("int", v.Int, func(val int) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt8() sequel.ColumnConvertClause[int8] {
	return sequel.Column("int_8", v.Int8, func(val int8) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt16() sequel.ColumnConvertClause[int16] {
	return sequel.Column("int_16", v.Int16, func(val int16) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt32() sequel.ColumnConvertClause[int32] {
	return sequel.Column("int_32", v.Int32, func(val int32) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnInt64() sequel.ColumnClause[int64] {
	return sequel.BasicColumn("int_64", v.Int64)
}
func (v Primitive) ColumnUint() sequel.ColumnConvertClause[uint] {
	return sequel.Column("uint", v.Uint, func(val uint) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint8() sequel.ColumnConvertClause[uint8] {
	return sequel.Column("uint_8", v.Uint8, func(val uint8) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint16() sequel.ColumnConvertClause[uint16] {
	return sequel.Column("uint_16", v.Uint16, func(val uint16) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint32() sequel.ColumnConvertClause[uint32] {
	return sequel.Column("uint_32", v.Uint32, func(val uint32) any {
		return (int64)(val)
	})
}
func (v Primitive) ColumnUint64() sequel.ColumnConvertClause[uint64] {
	return sequel.Column("uint_64", v.Uint64, func(val uint64) any {
		return val
	})
}
func (v Primitive) ColumnF32() sequel.ColumnConvertClause[float32] {
	return sequel.Column("f_32", v.F32, func(val float32) any {
		return (float64)(val)
	})
}
func (v Primitive) ColumnF64() sequel.ColumnClause[float64] {
	return sequel.BasicColumn("f_64", v.F64)
}
func (v Primitive) ColumnTime() sequel.ColumnClause[time.Time] {
	return sequel.BasicColumn("time", v.Time)
}
