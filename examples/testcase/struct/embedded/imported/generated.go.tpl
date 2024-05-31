package embedded

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `b` (`date` DATE NOT NULL,`time` VARCHAR(255) NOT NULL);"
}
func (B) TableName() string {
	return "b"
}
func (B) InsertOneStmt() string {
	return "INSERT INTO b (date,time) VALUES (?,?);"
}
func (B) InsertVarQuery() string {
	return "(?,?)"
}
func (B) Columns() []string {
	return []string{"date", "time"}
}
func (v B) Values() []any {
	return []any{types.TextMarshaler(v.DateTime.Date), types.TextMarshaler(v.DateTime.Time)}
}
func (v *B) Addrs() []any {
	return []any{types.Date(&v.DateTime.Date), types.TextUnmarshaler(&v.DateTime.Time)}
}
func (v B) GetDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("date", v.DateTime.Date, func(val civil.Date) driver.Value { return types.TextMarshaler(val) })
}
func (v B) GetTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column("time", v.DateTime.Time, func(val civil.Time) driver.Value { return types.TextMarshaler(val) })
}
