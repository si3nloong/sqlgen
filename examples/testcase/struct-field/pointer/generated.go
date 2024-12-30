package pointer

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Ptr) TableName() string {
	return "`ptr`"
}
func (Ptr) HasPK()      {}
func (Ptr) IsAutoIncr() {}
func (v *Ptr) ScanAutoIncr(val int64) error {
	v.ID = int64(val)
	return nil
}
func (v Ptr) PK() (string, int, any) {
	return "`id`", 0, v.ID
}
func (Ptr) Columns() []string {
	return []string{"`id`", "`str`", "`bytes`", "`bool`", "`int`", "`int_8`", "`int_16`", "`int_32`", "`int_64`", "`uint`", "`uint_8`", "`uint_16`", "`uint_32`", "`uint_64`", "`f_32`", "`f_64`", "`time`", "`nested`", "`embedded_time`", "`any_time`"} // 20
}
func (v Ptr) Values() []any {
	return []any{
		v.StrValue(),          //  1 - str
		v.BytesValue(),        //  2 - bytes
		v.BoolValue(),         //  3 - bool
		v.IntValue(),          //  4 - int
		v.Int8Value(),         //  5 - int_8
		v.Int16Value(),        //  6 - int_16
		v.Int32Value(),        //  7 - int_32
		v.Int64Value(),        //  8 - int_64
		v.UintValue(),         //  9 - uint
		v.Uint8Value(),        // 10 - uint_8
		v.Uint16Value(),       // 11 - uint_16
		v.Uint32Value(),       // 12 - uint_32
		v.Uint64Value(),       // 13 - uint_64
		v.F32Value(),          // 14 - f_32
		v.F64Value(),          // 15 - f_64
		v.TimeValue(),         // 16 - time
		v.NestedValue(),       // 17 - nested
		v.EmbeddedTimeValue(), // 18 - embedded_time
		v.AnyTimeValue(),      // 19 - any_time
	}
}
func (v *Ptr) Addrs() []any {
	if v.Str == nil {
		v.Str = new(string)
	}
	if v.Bytes == nil {
		v.Bytes = new([]byte)
	}
	if v.Bool == nil {
		v.Bool = new(bool)
	}
	if v.Int == nil {
		v.Int = new(int)
	}
	if v.Int8 == nil {
		v.Int8 = new(int8)
	}
	if v.Int16 == nil {
		v.Int16 = new(int16)
	}
	if v.Int32 == nil {
		v.Int32 = new(int32)
	}
	if v.Int64 == nil {
		v.Int64 = new(int64)
	}
	if v.Uint == nil {
		v.Uint = new(uint)
	}
	if v.Uint8 == nil {
		v.Uint8 = new(uint8)
	}
	if v.Uint16 == nil {
		v.Uint16 = new(uint16)
	}
	if v.Uint32 == nil {
		v.Uint32 = new(uint32)
	}
	if v.Uint64 == nil {
		v.Uint64 = new(uint64)
	}
	if v.F32 == nil {
		v.F32 = new(float32)
	}
	if v.F64 == nil {
		v.F64 = new(float64)
	}
	if v.Time == nil {
		v.Time = new(time.Time)
	}
	if v.Nested == nil {
		v.Nested = new(nested)
	}
	if v.deepNested == nil {
		v.deepNested = new(deepNested)
	}
	if v.deepNested.embedded == nil {
		v.deepNested.embedded = new(embedded)
	}
	if v.deepNested.embedded.EmbeddedTime == nil {
		v.deepNested.embedded.EmbeddedTime = new(time.Time)
	}
	return []any{
		&v.ID,                                                     //  0 - id
		encoding.StringScanner[string](&v.Str),                    //  1 - str
		encoding.StringScanner[[]byte](&v.Bytes),                  //  2 - bytes
		encoding.BoolScanner[bool](&v.Bool),                       //  3 - bool
		encoding.IntScanner[int](&v.Int),                          //  4 - int
		encoding.Int8Scanner[int8](&v.Int8),                       //  5 - int_8
		encoding.Int16Scanner[int16](&v.Int16),                    //  6 - int_16
		encoding.Int32Scanner[int32](&v.Int32),                    //  7 - int_32
		encoding.Int64Scanner[int64](&v.Int64),                    //  8 - int_64
		encoding.UintScanner[uint](&v.Uint),                       //  9 - uint
		encoding.Uint8Scanner[uint8](&v.Uint8),                    // 10 - uint_8
		encoding.Uint16Scanner[uint16](&v.Uint16),                 // 11 - uint_16
		encoding.Uint32Scanner[uint32](&v.Uint32),                 // 12 - uint_32
		encoding.Uint64Scanner[uint64](&v.Uint64),                 // 13 - uint_64
		encoding.Float32Scanner[float32](&v.F32),                  // 14 - f_32
		encoding.Float64Scanner[float64](&v.F64),                  // 15 - f_64
		encoding.TimeScanner(&v.Time),                             // 16 - time
		encoding.JSONScanner(&v.Nested),                           // 17 - nested
		encoding.TimeScanner(&v.deepNested.embedded.EmbeddedTime), // 18 - embedded_time
		&v.deepNested.embedded.AnyTime,                            // 19 - any_time
	}
}
func (Ptr) InsertColumns() []string {
	return []string{"`str`", "`bytes`", "`bool`", "`int`", "`int_8`", "`int_16`", "`int_32`", "`int_64`", "`uint`", "`uint_8`", "`uint_16`", "`uint_32`", "`uint_64`", "`f_32`", "`f_64`", "`time`", "`nested`", "`embedded_time`", "`any_time`"} // 19
}
func (Ptr) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" // 19
}
func (v Ptr) InsertOneStmt() (string, []any) {
	return "INSERT INTO `ptr` (`str`,`bytes`,`bool`,`int`,`int_8`,`int_16`,`int_32`,`int_64`,`uint`,`uint_8`,`uint_16`,`uint_32`,`uint_64`,`f_32`,`f_64`,`time`,`nested`,`embedded_time`,`any_time`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", []any{v.StrValue(), v.BytesValue(), v.BoolValue(), v.IntValue(), v.Int8Value(), v.Int16Value(), v.Int32Value(), v.Int64Value(), v.UintValue(), v.Uint8Value(), v.Uint16Value(), v.Uint32Value(), v.Uint64Value(), v.F32Value(), v.F64Value(), v.TimeValue(), v.NestedValue(), v.EmbeddedTimeValue(), v.AnyTimeValue()}
}
func (v Ptr) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`str`,`bytes`,`bool`,`int`,`int_8`,`int_16`,`int_32`,`int_64`,`uint`,`uint_8`,`uint_16`,`uint_32`,`uint_64`,`f_32`,`f_64`,`time`,`nested`,`embedded_time`,`any_time` FROM `ptr` WHERE `id` = ? LIMIT 1;", []any{v.ID}
}
func (v Ptr) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `ptr` SET `str` = ?,`bytes` = ?,`bool` = ?,`int` = ?,`int_8` = ?,`int_16` = ?,`int_32` = ?,`int_64` = ?,`uint` = ?,`uint_8` = ?,`uint_16` = ?,`uint_32` = ?,`uint_64` = ?,`f_32` = ?,`f_64` = ?,`time` = ?,`nested` = ?,`embedded_time` = ?,`any_time` = ? WHERE `id` = ?;", []any{v.StrValue(), v.BytesValue(), v.BoolValue(), v.IntValue(), v.Int8Value(), v.Int16Value(), v.Int32Value(), v.Int64Value(), v.UintValue(), v.Uint8Value(), v.Uint16Value(), v.Uint32Value(), v.Uint64Value(), v.F32Value(), v.F64Value(), v.TimeValue(), v.NestedValue(), v.EmbeddedTimeValue(), v.AnyTimeValue(), v.ID}
}
func (v Ptr) IDValue() driver.Value {
	return v.ID
}
func (v Ptr) StrValue() driver.Value {
	if v.Str != nil {
		return *v.Str
	}
	return nil
}
func (v Ptr) BytesValue() driver.Value {
	if v.Bytes != nil {
		return string(*v.Bytes)
	}
	return nil
}
func (v Ptr) BoolValue() driver.Value {
	if v.Bool != nil {
		return *v.Bool
	}
	return nil
}
func (v Ptr) IntValue() driver.Value {
	if v.Int != nil {
		return (int64)(*v.Int)
	}
	return nil
}
func (v Ptr) Int8Value() driver.Value {
	if v.Int8 != nil {
		return (int64)(*v.Int8)
	}
	return nil
}
func (v Ptr) Int16Value() driver.Value {
	if v.Int16 != nil {
		return (int64)(*v.Int16)
	}
	return nil
}
func (v Ptr) Int32Value() driver.Value {
	if v.Int32 != nil {
		return (int64)(*v.Int32)
	}
	return nil
}
func (v Ptr) Int64Value() driver.Value {
	if v.Int64 != nil {
		return *v.Int64
	}
	return nil
}
func (v Ptr) UintValue() driver.Value {
	if v.Uint != nil {
		return (int64)(*v.Uint)
	}
	return nil
}
func (v Ptr) Uint8Value() driver.Value {
	if v.Uint8 != nil {
		return (int64)(*v.Uint8)
	}
	return nil
}
func (v Ptr) Uint16Value() driver.Value {
	if v.Uint16 != nil {
		return (int64)(*v.Uint16)
	}
	return nil
}
func (v Ptr) Uint32Value() driver.Value {
	if v.Uint32 != nil {
		return (int64)(*v.Uint32)
	}
	return nil
}
func (v Ptr) Uint64Value() driver.Value {
	if v.Uint64 != nil {
		return (int64)(*v.Uint64)
	}
	return nil
}
func (v Ptr) F32Value() driver.Value {
	if v.F32 != nil {
		return (float64)(*v.F32)
	}
	return nil
}
func (v Ptr) F64Value() driver.Value {
	if v.F64 != nil {
		return *v.F64
	}
	return nil
}
func (v Ptr) TimeValue() driver.Value {
	if v.Time != nil {
		return *v.Time
	}
	return nil
}
func (v Ptr) NestedValue() driver.Value {
	if v.Nested != nil {
		return encoding.JSONValue(*v.Nested)
	}
	return nil
}
func (v Ptr) EmbeddedTimeValue() driver.Value {
	if v.deepNested != nil {
		if v.deepNested.embedded != nil {
			if v.deepNested.embedded.EmbeddedTime != nil {
				return *v.deepNested.embedded.EmbeddedTime
			}
		}
	}
	return nil
}
func (v Ptr) AnyTimeValue() driver.Value {
	if v.deepNested != nil {
		if v.deepNested.embedded != nil {
			return v.deepNested.embedded.AnyTime
		}
	}
	return nil
}
func (v Ptr) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("`id`", v.ID, func(val int64) driver.Value {
		return val
	})
}
func (v Ptr) GetStr() sequel.ColumnValuer[*string] {
	return sequel.Column("`str`", v.Str, func(val *string) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v Ptr) GetBytes() sequel.ColumnValuer[*[]byte] {
	return sequel.Column("`bytes`", v.Bytes, func(val *[]byte) driver.Value {
		if val != nil {
			return string(*val)
		}
		return nil
	})
}
func (v Ptr) GetBool() sequel.ColumnValuer[*bool] {
	return sequel.Column("`bool`", v.Bool, func(val *bool) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v Ptr) GetInt() sequel.ColumnValuer[*int] {
	return sequel.Column("`int`", v.Int, func(val *int) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetInt8() sequel.ColumnValuer[*int8] {
	return sequel.Column("`int_8`", v.Int8, func(val *int8) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetInt16() sequel.ColumnValuer[*int16] {
	return sequel.Column("`int_16`", v.Int16, func(val *int16) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetInt32() sequel.ColumnValuer[*int32] {
	return sequel.Column("`int_32`", v.Int32, func(val *int32) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetInt64() sequel.ColumnValuer[*int64] {
	return sequel.Column("`int_64`", v.Int64, func(val *int64) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v Ptr) GetUint() sequel.ColumnValuer[*uint] {
	return sequel.Column("`uint`", v.Uint, func(val *uint) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetUint8() sequel.ColumnValuer[*uint8] {
	return sequel.Column("`uint_8`", v.Uint8, func(val *uint8) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetUint16() sequel.ColumnValuer[*uint16] {
	return sequel.Column("`uint_16`", v.Uint16, func(val *uint16) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetUint32() sequel.ColumnValuer[*uint32] {
	return sequel.Column("`uint_32`", v.Uint32, func(val *uint32) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetUint64() sequel.ColumnValuer[*uint64] {
	return sequel.Column("`uint_64`", v.Uint64, func(val *uint64) driver.Value {
		if val != nil {
			return (int64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetF32() sequel.ColumnValuer[*float32] {
	return sequel.Column("`f_32`", v.F32, func(val *float32) driver.Value {
		if val != nil {
			return (float64)(*val)
		}
		return nil
	})
}
func (v Ptr) GetF64() sequel.ColumnValuer[*float64] {
	return sequel.Column("`f_64`", v.F64, func(val *float64) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v Ptr) GetTime() sequel.ColumnValuer[*time.Time] {
	return sequel.Column("`time`", v.Time, func(val *time.Time) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v Ptr) GetNested() sequel.ColumnValuer[*nested] {
	return sequel.Column("`nested`", v.Nested, func(val *nested) driver.Value {
		if val != nil {
			return encoding.JSONValue(*val)
		}
		return nil
	})
}
func (v Ptr) GetEmbeddedTime() sequel.ColumnValuer[*time.Time] {
	return sequel.Column("`embedded_time`", v.deepNested.embedded.EmbeddedTime, func(val *time.Time) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v Ptr) GetAnyTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("`any_time`", v.deepNested.embedded.AnyTime, func(val time.Time) driver.Value {
		return val
	})
}
