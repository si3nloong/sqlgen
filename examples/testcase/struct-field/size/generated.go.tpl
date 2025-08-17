package size

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
)

func (Size) TableName() string {
	return "size"
}
func (Size) Columns() []string {
	return []string{"str", "timestamp", "time"} // 3
}
func (v Size) Values() []any {
	return []any{
		v.Str,       // 0 - str
		v.Timestamp, // 1 - timestamp
		v.Time,      // 2 - time
	}
}
func (v *Size) Addrs() []any {
	return []any{
		&v.Str,       // 0 - str
		&v.Timestamp, // 1 - timestamp
		&v.Time,      // 2 - time
	}
}
func (Size) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v Size) InsertOneStmt() (string, []any) {
	return "INSERT INTO size (str,timestamp,time) VALUES (?,?,?);", v.Values()
}
func (v Size) StrValue() driver.Value {
	return v.Str
}
func (v Size) TimestampValue() driver.Value {
	return v.Timestamp
}
func (v Size) TimeValue() driver.Value {
	return v.Time
}
func (v Size) ColumnStr() sequel.ColumnValuer[string] {
	return sequel.Column("str", v.Str, func(val string) driver.Value {
		return val
	})
}
func (v Size) ColumnTimestamp() sequel.ColumnValuer[time.Time] {
	return sequel.Column("timestamp", v.Timestamp, func(val time.Time) driver.Value {
		return val
	})
}
func (v Size) ColumnTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("time", v.Time, func(val time.Time) driver.Value {
		return val
	})
}
