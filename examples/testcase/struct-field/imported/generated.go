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
	return "INSERT INTO `model` (`str`,`bool`,`raw_bytes`,`int_16`,`int_32`,`int_64`,`float_64`,`time`) VALUES (?,?,?,?,?,?,?,?);", v.Values()
}
func (v Model) StrValue() any {
	return v.Str
}
func (v Model) BoolValue() any {
	return v.Bool
}
func (v Model) RawBytesValue() any {
	return string(v.RawBytes)
}
func (v Model) Int16Value() any {
	return v.Int16
}
func (v Model) Int32Value() any {
	return v.Int32
}
func (v Model) Int64Value() any {
	return v.Int64
}
func (v Model) Float64Value() any {
	return v.Float64
}
func (v Model) TimeValue() any {
	return v.Time
}
func (v Model) ColumnStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("str", v.Str, func(val sql.NullString) driver.Value {
		return val
	})
}
func (v Model) ColumnBool() sequel.ColumnValuer[sql.NullBool] {
	return sequel.Column("bool", v.Bool, func(val sql.NullBool) driver.Value {
		return val
	})
}
func (v Model) ColumnRawBytes() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column("raw_bytes", v.RawBytes, func(val sql.RawBytes) driver.Value {
		return string(val)
	})
}
func (v Model) ColumnInt16() sequel.ColumnValuer[sql.NullInt16] {
	return sequel.Column("int_16", v.Int16, func(val sql.NullInt16) driver.Value {
		return val
	})
}
func (v Model) ColumnInt32() sequel.ColumnValuer[sql.NullInt32] {
	return sequel.Column("int_32", v.Int32, func(val sql.NullInt32) driver.Value {
		return val
	})
}
func (v Model) ColumnInt64() sequel.ColumnValuer[sql.NullInt64] {
	return sequel.Column("int_64", v.Int64, func(val sql.NullInt64) driver.Value {
		return val
	})
}
func (v Model) ColumnFloat64() sequel.ColumnValuer[sql.NullFloat64] {
	return sequel.Column("float_64", v.Float64, func(val sql.NullFloat64) driver.Value {
		return val
	})
}
func (v Model) ColumnTime() sequel.ColumnValuer[sql.NullTime] {
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
	return "INSERT INTO `some` (`id`) VALUES (?);", v.Values()
}
func (v Some) IDValue() any {
	return v.ID
}
func (v Some) ColumnID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value {
		return val
	})
}
