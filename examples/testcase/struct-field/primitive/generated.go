package primitive

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Primitive) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (Primitive) TableName() string {
	return "`primitive`"
}
func (Primitive) ColumnNames() []string {
	return []string{"`str`", "`bytes`", "`bool`", "`int`", "`int_8`", "`int_16`", "`int_32`", "`int_64`", "`uint`", "`uint_8`", "`uint_16`", "`uint_32`", "`uint_64`", "`f_32`", "`f_64`", "`time`"}
}
func (v Primitive) Values() []any {
	return []any{string(v.Str), string(v.Bytes), bool(v.Bool), int64(v.Int), int64(v.Int8), int64(v.Int16), int64(v.Int32), int64(v.Int64), int64(v.Uint), int64(v.Uint8), int64(v.Uint16), int64(v.Uint32), int64(v.Uint64), float64(v.F32), float64(v.F64), time.Time(v.Time)}
}
func (v *Primitive) Addrs() []any {
	return []any{types.String(&v.Str), types.String(&v.Bytes), types.Bool(&v.Bool), types.Integer(&v.Int), types.Integer(&v.Int8), types.Integer(&v.Int16), types.Integer(&v.Int32), types.Integer(&v.Int64), types.Integer(&v.Uint), types.Integer(&v.Uint8), types.Integer(&v.Uint16), types.Integer(&v.Uint32), types.Integer(&v.Uint64), types.Float(&v.F32), types.Float(&v.F64), (*time.Time)(&v.Time)}
}
func (Primitive) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (v Primitive) InsertOneStmt() (string, []any) {
	return "INSERT INTO `primitive` (`str`,`bytes`,`bool`,`int`,`int_8`,`int_16`,`int_32`,`int_64`,`uint`,`uint_8`,`uint_16`,`uint_32`,`uint_64`,`f_32`,`f_64`,`time`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", v.Values()
}
func (v Primitive) GetStr() sequel.ColumnValuer[string] {
	return sequel.Column("`str`", v.Str, func(val string) driver.Value { return string(val) })
}
func (v Primitive) GetBytes() sequel.ColumnValuer[[]byte] {
	return sequel.Column("`bytes`", v.Bytes, func(val []byte) driver.Value { return string(val) })
}
func (v Primitive) GetBool() sequel.ColumnValuer[bool] {
	return sequel.Column("`bool`", v.Bool, func(val bool) driver.Value { return bool(val) })
}
func (v Primitive) GetInt() sequel.ColumnValuer[int] {
	return sequel.Column("`int`", v.Int, func(val int) driver.Value { return int64(val) })
}
func (v Primitive) GetInt8() sequel.ColumnValuer[int8] {
	return sequel.Column("`int_8`", v.Int8, func(val int8) driver.Value { return int64(val) })
}
func (v Primitive) GetInt16() sequel.ColumnValuer[int16] {
	return sequel.Column("`int_16`", v.Int16, func(val int16) driver.Value { return int64(val) })
}
func (v Primitive) GetInt32() sequel.ColumnValuer[int32] {
	return sequel.Column("`int_32`", v.Int32, func(val int32) driver.Value { return int64(val) })
}
func (v Primitive) GetInt64() sequel.ColumnValuer[int64] {
	return sequel.Column("`int_64`", v.Int64, func(val int64) driver.Value { return int64(val) })
}
func (v Primitive) GetUint() sequel.ColumnValuer[uint] {
	return sequel.Column("`uint`", v.Uint, func(val uint) driver.Value { return int64(val) })
}
func (v Primitive) GetUint8() sequel.ColumnValuer[uint8] {
	return sequel.Column("`uint_8`", v.Uint8, func(val uint8) driver.Value { return int64(val) })
}
func (v Primitive) GetUint16() sequel.ColumnValuer[uint16] {
	return sequel.Column("`uint_16`", v.Uint16, func(val uint16) driver.Value { return int64(val) })
}
func (v Primitive) GetUint32() sequel.ColumnValuer[uint32] {
	return sequel.Column("`uint_32`", v.Uint32, func(val uint32) driver.Value { return int64(val) })
}
func (v Primitive) GetUint64() sequel.ColumnValuer[uint64] {
	return sequel.Column("`uint_64`", v.Uint64, func(val uint64) driver.Value { return int64(val) })
}
func (v Primitive) GetF32() sequel.ColumnValuer[float32] {
	return sequel.Column("`f_32`", v.F32, func(val float32) driver.Value { return float64(val) })
}
func (v Primitive) GetF64() sequel.ColumnValuer[float64] {
	return sequel.Column("`f_64`", v.F64, func(val float64) driver.Value { return float64(val) })
}
func (v Primitive) GetTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("`time`", v.Time, func(val time.Time) driver.Value { return time.Time(val) })
}
