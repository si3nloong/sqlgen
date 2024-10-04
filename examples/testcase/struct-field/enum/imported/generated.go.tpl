package imported

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (ImportedEnum) TableName() string {
	return "imported_enum"
}
func (ImportedEnum) Columns() []string {
	return []string{"weekday"}
}
func (v ImportedEnum) Values() []any {
	return []any{(int64)(v.Weekday)}
}
func (v *ImportedEnum) Addrs() []any {
	return []any{types.Integer(&v.Weekday)}
}
func (ImportedEnum) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v ImportedEnum) InsertOneStmt() (string, []any) {
	return "INSERT INTO imported_enum (weekday) VALUES (?);", v.Values()
}
func (v ImportedEnum) GetWeekday() sequel.ColumnValuer[time.Weekday] {
	return sequel.Column("weekday", v.Weekday, func(val time.Weekday) driver.Value { return (int64)(val) })
}
