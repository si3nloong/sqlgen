package size

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Size) TableName() string {
	return "size"
}
func (Size) Columns() []string {
	return []string{"str", "timestamp", "time"}
}
func (v Size) Values() []any {
	return []any{(string)(v.Str), (time.Time)(v.Timestamp), (time.Time)(v.Time)}
}
func (v *Size) Addrs() []any {
	return []any{types.String(&v.Str), (*time.Time)(&v.Timestamp), (*time.Time)(&v.Time)}
}
func (Size) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Size) InsertOneStmt() (string, []any) {
	return "INSERT INTO size (str,timestamp,time) VALUES (?,?,?);", v.Values()
}
func (v Size) GetStr() sequel.ColumnValuer[string] {
	return sequel.Column("str", v.Str, func(val string) driver.Value {
		return (string)(val)
	})
}
func (v Size) GetTimestamp() sequel.ColumnValuer[time.Time] {
	return sequel.Column("timestamp", v.Timestamp, func(val time.Time) driver.Value {
		return (time.Time)(val)
	})
}
func (v Size) GetTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("time", v.Time, func(val time.Time) driver.Value {
		return (time.Time)(val)
	})
}
