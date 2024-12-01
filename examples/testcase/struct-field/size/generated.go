package size

import (
	"database/sql/driver"
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
func (v Size) GetStr() driver.Value {
	return v.Str
}
func (v Size) GetTimestamp() driver.Value {
	return v.Timestamp
}
func (v Size) GetTime() driver.Value {
	return v.Time
}
