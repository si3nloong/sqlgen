package tablename

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
)

func (CustomTableName1) TableName() string {
	return "CustomTableName_1"
}
func (CustomTableName1) Columns() []string {
	return []string{"text"} // 1
}
func (v CustomTableName1) Values() []any {
	return []any{
		v.Text, // 0 - text
	}
}
func (v *CustomTableName1) Addrs() []any {
	return []any{
		&v.Text, // 0 - text
	}
}
func (CustomTableName1) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v CustomTableName1) InsertOneStmt() (string, []any) {
	return "INSERT INTO `CustomTableName_1` (`text`) VALUES (?);", v.Values()
}
func (v CustomTableName1) TextValue() any {
	return v.Text
}
func (v CustomTableName1) ColumnText() sequel.ColumnValuer[string] {
	return sequel.Column("text", v.Text, func(val string) driver.Value {
		return val
	})
}

func (CustomTableName2) TableName() string {
	return "table_2"
}
func (CustomTableName2) Columns() []string {
	return []string{"text"} // 1
}
func (v CustomTableName2) Values() []any {
	return []any{
		v.Text, // 0 - text
	}
}
func (v *CustomTableName2) Addrs() []any {
	return []any{
		&v.Text, // 0 - text
	}
}
func (CustomTableName2) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v CustomTableName2) InsertOneStmt() (string, []any) {
	return "INSERT INTO `table_2` (`text`) VALUES (?);", v.Values()
}
func (v CustomTableName2) TextValue() any {
	return v.Text
}
func (v CustomTableName2) ColumnText() sequel.ColumnValuer[string] {
	return sequel.Column("text", v.Text, func(val string) driver.Value {
		return val
	})
}

func (CustomTableName3) TableName() string {
	return "table_3"
}
func (CustomTableName3) Columns() []string {
	return []string{"text"} // 1
}
func (v CustomTableName3) Values() []any {
	return []any{
		v.Text, // 0 - text
	}
}
func (v *CustomTableName3) Addrs() []any {
	return []any{
		&v.Text, // 0 - text
	}
}
func (CustomTableName3) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v CustomTableName3) InsertOneStmt() (string, []any) {
	return "INSERT INTO `table_3` (`text`) VALUES (?);", v.Values()
}
func (v CustomTableName3) TextValue() any {
	return v.Text
}
func (v CustomTableName3) ColumnText() sequel.ColumnValuer[string] {
	return sequel.Column("text", v.Text, func(val string) driver.Value {
		return val
	})
}
