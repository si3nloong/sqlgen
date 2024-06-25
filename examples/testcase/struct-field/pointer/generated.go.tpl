package pointer

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Ptr) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `ptr` (`id` BIGINT NOT NULL AUTO_INCREMENT,`str` VARCHAR(255),`bytes` BLOB,`bool` BOOL,`int` INTEGER,`int_8` TINYINT,`int_16` SMALLINT,`int_32` MEDIUMINT,`int_64` BIGINT,`uint` INTEGER,`uint_8` TINYINT UNSIGNED,`uint_16` SMALLINT UNSIGNED,`uint_32` MEDIUMINT UNSIGNED,`uint_64` BIGINT UNSIGNED,`f_32` FLOAT,`f_64` FLOAT,`time` DATETIME(6),PRIMARY KEY (`id`));"
}
func (Ptr) TableName() string {
	return "ptr"
}
func (Ptr) InsertOneStmt() string {
	return "INSERT INTO ptr (str,bytes,bool,int,int_8,int_16,int_32,int_64,uint,uint_8,uint_16,uint_32,uint_64,f_32,f_64,time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
}
func (Ptr) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
}
func (Ptr) Columns() []string {
	return []string{"id", "str", "bytes", "bool", "int", "int_8", "int_16", "int_32", "int_64", "uint", "uint_8", "uint_16", "uint_32", "uint_64", "f_32", "f_64", "time"}
}
func (Ptr) HasPK()      {}
func (Ptr) IsAutoIncr() {}
func (v Ptr) PK() (string, int, any) {
	return "id", 0, int64(v.ID)
}
func (Ptr) FindByPKStmt() string {
	return "SELECT id,str,bytes,bool,int,int_8,int_16,int_32,int_64,uint,uint_8,uint_16,uint_32,uint_64,f_32,f_64,time FROM ptr WHERE id = ? LIMIT 1;"
}
func (Ptr) UpdateByPKStmt() string {
	return "UPDATE ptr SET str = ?,bytes = ?,bool = ?,int = ?,int_8 = ?,int_16 = ?,int_32 = ?,int_64 = ?,uint = ?,uint_8 = ?,uint_16 = ?,uint_32 = ?,uint_64 = ?,f_32 = ?,f_64 = ?,time = ? WHERE id = ? LIMIT 1;"
}
func (v Ptr) Values() []any {
	return []any{int64(v.ID), types.String(v.Str), types.String(v.Bytes), types.Bool(v.Bool), types.Integer(v.Int), types.Integer(v.Int8), types.Integer(v.Int16), types.Integer(v.Int32), types.Integer(v.Int64), types.Integer(v.Uint), types.Integer(v.Uint8), types.Integer(v.Uint16), types.Integer(v.Uint32), types.Integer(v.Uint64), types.Float(v.F32), types.Float(v.F64), types.Time(v.Time)}
}
func (v *Ptr) Addrs() []any {
	return []any{types.Integer(&v.ID), types.PtrOfString(&v.Str), types.PtrOfString(&v.Bytes), types.PtrOfBool(&v.Bool), types.PtrOfInt(&v.Int), types.PtrOfInt(&v.Int8), types.PtrOfInt(&v.Int16), types.PtrOfInt(&v.Int32), types.PtrOfInt(&v.Int64), types.PtrOfInt(&v.Uint), types.PtrOfInt(&v.Uint8), types.PtrOfInt(&v.Uint16), types.PtrOfInt(&v.Uint32), types.PtrOfInt(&v.Uint64), types.PtrOfFloat(&v.F32), types.PtrOfFloat(&v.F64), types.PtrOfTime(&v.Time)}
}
func (v Ptr) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
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
	return sequel.Column("f_32", v.F32, func(val *float32) driver.Value { return types.Float(val) })
}
func (v Ptr) GetF64() sequel.ColumnValuer[*float64] {
	return sequel.Column("f_64", v.F64, func(val *float64) driver.Value { return types.Float(val) })
}
func (v Ptr) GetTime() sequel.ColumnValuer[*time.Time] {
	return sequel.Column("time", v.Time, func(val *time.Time) driver.Value { return types.Time(val) })
}
