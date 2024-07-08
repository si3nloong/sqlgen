package imported

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (ImportedEnum) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`weekday`", Definition: "`weekday` INTEGER NOT NULL DEFAULT 0"},
		},
	}
}
func (ImportedEnum) TableName() string {
	return "`imported_enum`"
}
func (ImportedEnum) ColumnNames() []string {
	return []string{"`weekday`"}
}
func (v ImportedEnum) Values() []any {
	return []any{int64(v.Weekday)}
}
func (v *ImportedEnum) Addrs() []any {
	return []any{types.Integer(&v.Weekday)}
}
func (ImportedEnum) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v ImportedEnum) InsertOneStmt() (string, []any) {
	return "INSERT INTO `imported_enum` (`weekday`) VALUES (?);", v.Values()
}
func (v ImportedEnum) GetWeekday() sequel.ColumnValuer[time.Weekday] {
	return sequel.Column("`weekday`", v.Weekday, func(val time.Weekday) driver.Value { return int64(val) })
}
