package imported

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Model) TableName() string {
	return "model"
}
func (Model) Columns() []string {
	return []string{"str", "bool", "raw_bytes", "int_16", "int_32", "int_64", "float_64", "time"} // 8
}
func (v Model) Values() []any {
	return []any{
		v.Str,              // 0 - str
		v.Bool,             // 1 - bool
		string(v.RawBytes), // 2 - raw_bytes
		v.Int16,            // 3 - int_16
		v.Int32,            // 4 - int_32
		v.Int64,            // 5 - int_64
		v.Float64,          // 6 - float_64
		v.Time,             // 7 - time
	}
}
func (v *Model) Addrs() []any {
	return []any{
		&v.Str,  // 0 - str
		&v.Bool, // 1 - bool
		encoding.StringScanner[sql.RawBytes](&v.RawBytes), // 2 - raw_bytes
		&v.Int16,   // 3 - int_16
		&v.Int32,   // 4 - int_32
		&v.Int64,   // 5 - int_64
		&v.Float64, // 6 - float_64
		&v.Time,    // 7 - time
	}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?)" // 8
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO model (str,bool,raw_bytes,int_16,int_32,int_64,float_64,time) VALUES (?,?,?,?,?,?,?,?);", v.Values()
}
func (v Model) StrValue() driver.Value {
	return v.Str
}
func (v Model) BoolValue() driver.Value {
	return v.Bool
}
func (v Model) RawBytesValue() driver.Value {
	return string(v.RawBytes)
}
func (v Model) Int16Value() driver.Value {
	return v.Int16
}
func (v Model) Int32Value() driver.Value {
	return v.Int32
}
func (v Model) Int64Value() driver.Value {
	return v.Int64
}
func (v Model) Float64Value() driver.Value {
	return v.Float64
}
func (v Model) TimeValue() driver.Value {
	return v.Time
}
func (v Model) GetStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("str", v.Str, func(val sql.NullString) driver.Value {
		return val
	})
}
func (v Model) GetBool() sequel.ColumnValuer[sql.NullBool] {
	return sequel.Column("bool", v.Bool, func(val sql.NullBool) driver.Value {
		return val
	})
}
func (v Model) GetRawBytes() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column("raw_bytes", v.RawBytes, func(val sql.RawBytes) driver.Value {
		return string(val)
	})
}
func (v Model) GetInt16() sequel.ColumnValuer[sql.NullInt16] {
	return sequel.Column("int_16", v.Int16, func(val sql.NullInt16) driver.Value {
		return val
	})
}
func (v Model) GetInt32() sequel.ColumnValuer[sql.NullInt32] {
	return sequel.Column("int_32", v.Int32, func(val sql.NullInt32) driver.Value {
		return val
	})
}
func (v Model) GetInt64() sequel.ColumnValuer[sql.NullInt64] {
	return sequel.Column("int_64", v.Int64, func(val sql.NullInt64) driver.Value {
		return val
	})
}
func (v Model) GetFloat64() sequel.ColumnValuer[sql.NullFloat64] {
	return sequel.Column("float_64", v.Float64, func(val sql.NullFloat64) driver.Value {
		return val
	})
}
func (v Model) GetTime() sequel.ColumnValuer[sql.NullTime] {
	return sequel.Column("time", v.Time, func(val sql.NullTime) driver.Value {
		return val
	})
}

func (Some) TableName() string {
	return "some"
}
func (Some) Columns() []string {
	return []string{"id"} // 1
}
func (v Some) Values() []any {
	return []any{
		v.ID, // 0 - id
	}
}
func (v *Some) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
	}
}
func (Some) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v Some) InsertOneStmt() (string, []any) {
	return "INSERT INTO some (id) VALUES (?);", v.Values()
}
func (v Some) IDValue() driver.Value {
	return v.ID
}
func (v Some) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value {
		return val
	})
}
