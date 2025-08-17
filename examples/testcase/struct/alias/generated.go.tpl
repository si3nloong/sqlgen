package aliasstruct

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (A) Columns() []string {
	return []string{"date", "time"} // 2
}
func (v A) Values() []any {
	return []any{
		encoding.TextValue(v.Date), // 0 - date
		encoding.TextValue(v.Time), // 1 - time
	}
}
func (v *A) Addrs() []any {
	return []any{
		encoding.TextScanner[civil.Date](&v.Date), // 0 - date
		encoding.TextScanner[civil.Time](&v.Time), // 1 - time
	}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (date,time) VALUES (?,?);", v.Values()
}
func (v A) DateValue() driver.Value {
	return encoding.TextValue(v.Date)
}
func (v A) TimeValue() driver.Value {
	return encoding.TextValue(v.Time)
}
func (v A) ColumnDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("date", v.Date, func(val civil.Date) driver.Value {
		return encoding.TextValue(val)
	})
}
func (v A) ColumnTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column("time", v.Time, func(val civil.Time) driver.Value {
		return encoding.TextValue(val)
	})
}

func (C) TableName() string {
	return "c"
}
func (C) Columns() []string {
	return []string{"string", "valid"} // 2
}
func (v C) Values() []any {
	return []any{
		v.String, // 0 - string
		v.Valid,  // 1 - valid
	}
}
func (v *C) Addrs() []any {
	return []any{
		&v.String, // 0 - string
		&v.Valid,  // 1 - valid
	}
}
func (C) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v C) InsertOneStmt() (string, []any) {
	return "INSERT INTO c (string,valid) VALUES (?,?);", v.Values()
}
func (v C) StringValue() driver.Value {
	return v.String
}
func (v C) ValidValue() driver.Value {
	return v.Valid
}
func (v C) ColumnString() sequel.ColumnValuer[string] {
	return sequel.Column("string", v.String, func(val string) driver.Value {
		return val
	})
}
func (v C) ColumnValid() sequel.ColumnValuer[bool] {
	return sequel.Column("valid", v.Valid, func(val bool) driver.Value {
		return val
	})
}
