package imported

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Model) TableName() string {
	return "model"
}
func (Model) Columns() []string {
	return []string{"str", "bool", "raw_bytes", "int_16", "int_32", "int_64", "float_64", "time"}
}
func (v Model) Values() []any {
	return []any{(driver.Valuer)(v.Str), (driver.Valuer)(v.Bool), string(v.RawBytes), (driver.Valuer)(v.Int16), (driver.Valuer)(v.Int32), (driver.Valuer)(v.Int64), (driver.Valuer)(v.Float64), (driver.Valuer)(v.Time)}
}
func (v *Model) Addrs() []any {
	return []any{(sql.Scanner)(&v.Str), (sql.Scanner)(&v.Bool), types.String(&v.RawBytes), (sql.Scanner)(&v.Int16), (sql.Scanner)(&v.Int32), (sql.Scanner)(&v.Int64), (sql.Scanner)(&v.Float64), (sql.Scanner)(&v.Time)}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO model (str,bool,raw_bytes,int_16,int_32,int_64,float_64,time) VALUES (?,?,?,?,?,?,?,?);", v.Values()
}
func (v Model) GetStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("str", v.Str, func(val sql.NullString) driver.Value {
		return (driver.Valuer)(val)
	})
}
func (v Model) GetBool() sequel.ColumnValuer[sql.NullBool] {
	return sequel.Column("bool", v.Bool, func(val sql.NullBool) driver.Value {
		return (driver.Valuer)(val)
	})
}
func (v Model) GetRawBytes() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column("raw_bytes", v.RawBytes, func(val sql.RawBytes) driver.Value {
		return string(val)
	})
}
func (v Model) GetInt16() sequel.ColumnValuer[sql.NullInt16] {
	return sequel.Column("int_16", v.Int16, func(val sql.NullInt16) driver.Value {
		return (driver.Valuer)(val)
	})
}
func (v Model) GetInt32() sequel.ColumnValuer[sql.NullInt32] {
	return sequel.Column("int_32", v.Int32, func(val sql.NullInt32) driver.Value {
		return (driver.Valuer)(val)
	})
}
func (v Model) GetInt64() sequel.ColumnValuer[sql.NullInt64] {
	return sequel.Column("int_64", v.Int64, func(val sql.NullInt64) driver.Value {
		return (driver.Valuer)(val)
	})
}
func (v Model) GetFloat64() sequel.ColumnValuer[sql.NullFloat64] {
	return sequel.Column("float_64", v.Float64, func(val sql.NullFloat64) driver.Value {
		return (driver.Valuer)(val)
	})
}
func (v Model) GetTime() sequel.ColumnValuer[sql.NullTime] {
	return sequel.Column("time", v.Time, func(val sql.NullTime) driver.Value {
		return (driver.Valuer)(val)
	})
}

func (Some) TableName() string {
	return "some"
}
func (Some) Columns() []string {
	return []string{"id"}
}
func (v Some) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *Some) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (Some) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v Some) InsertOneStmt() (string, []any) {
	return "INSERT INTO some (id) VALUES (?);", v.Values()
}
func (v Some) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value {
		return (driver.Valuer)(val)
	})
}
