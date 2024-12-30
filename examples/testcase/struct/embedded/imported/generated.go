package embedded

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"date", "time"} // 2
}
func (v B) Values() []any {
	return []any{
		encoding.TextValue(v.DateTime.Date), // 0 - date
		encoding.TextValue(v.DateTime.Time), // 1 - time
	}
}
func (v *B) Addrs() []any {
	return []any{
		encoding.TextScanner[civil.Date](&v.DateTime.Date), // 0 - date
		encoding.TextScanner[civil.Time](&v.DateTime.Time), // 1 - time
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (date,time) VALUES (?,?);", v.Values()
}
func (v B) DateValue() driver.Value {
	return encoding.TextValue(v.DateTime.Date)
}
func (v B) TimeValue() driver.Value {
	return encoding.TextValue(v.DateTime.Time)
}
func (v B) GetDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("date", v.DateTime.Date, func(val civil.Date) driver.Value {
		return encoding.TextValue(val)
	})
}
func (v B) GetTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column("time", v.DateTime.Time, func(val civil.Time) driver.Value {
		return encoding.TextValue(val)
	})
}
