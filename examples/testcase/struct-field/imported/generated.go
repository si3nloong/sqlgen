package imported

import (
	"database/sql"
	"database/sql/driver"

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
func (v Model) GetStr() driver.Value {
	return v.Str
}
func (v Model) GetBool() driver.Value {
	return v.Bool
}
func (v Model) GetRawBytes() driver.Value {
	return string(v.RawBytes)
}
func (v Model) GetInt16() driver.Value {
	return v.Int16
}
func (v Model) GetInt32() driver.Value {
	return v.Int32
}
func (v Model) GetInt64() driver.Value {
	return v.Int64
}
func (v Model) GetFloat64() driver.Value {
	return v.Float64
}
func (v Model) GetTime() driver.Value {
	return v.Time
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
func (v Some) GetID() driver.Value {
	return v.ID
}
