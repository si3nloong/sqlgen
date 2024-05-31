package size

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Size) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `size` (`str` VARCHAR(25) NOT NULL,`timestamp` DATETIME(6) NOT NULL,`time` DATETIME NOT NULL);"
}
func (Size) TableName() string {
	return "size"
}
func (Size) InsertOneStmt() string {
	return "INSERT INTO size (str,timestamp,time) VALUES (?,?,?);"
}
func (Size) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Size) Columns() []string {
	return []string{"str", "timestamp", "time"}
}
func (v Size) Values() []any {
	return []any{string(v.Str), time.Time(v.Timestamp), time.Time(v.Time)}
}
func (v *Size) Addrs() []any {
	return []any{types.String(&v.Str), (*time.Time)(&v.Timestamp), (*time.Time)(&v.Time)}
}
func (v Size) GetStr() sequel.ColumnValuer[string] {
	return sequel.Column("str", v.Str, func(val string) driver.Value { return string(val) })
}
func (v Size) GetTimestamp() sequel.ColumnValuer[time.Time] {
	return sequel.Column("timestamp", v.Timestamp, func(val time.Time) driver.Value { return time.Time(val) })
}
func (v Size) GetTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("time", v.Time, func(val time.Time) driver.Value { return time.Time(val) })
}
