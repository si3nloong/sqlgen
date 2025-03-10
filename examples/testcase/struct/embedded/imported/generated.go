package embedded

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"date", "time"}
}
func (v B) Values() []any {
	return []any{types.TextMarshaler(v.DateTime.Date), types.TextMarshaler(v.DateTime.Time)}
}
func (v *B) Addrs() []any {
	return []any{types.TextUnmarshaler(&v.DateTime.Date), types.TextUnmarshaler(&v.DateTime.Time)}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (date,time) VALUES (?,?);", v.Values()
}
func (v B) GetDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("date", v.DateTime.Date, func(val civil.Date) driver.Value {
		return types.TextMarshaler(val)
	})
}
func (v B) GetTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column("time", v.DateTime.Time, func(val civil.Time) driver.Value {
		return types.TextMarshaler(val)
	})
}
