package tablename

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (CustomTableName1) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (CustomTableName1) TableName() string {
	return "`CustomTableName_1`"
}
func (CustomTableName1) ColumnNames() []string {
	return []string{"`text`"}
}
func (v CustomTableName1) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName1) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (CustomTableName1) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v CustomTableName1) InsertOneStmt() (string, []any) {
	return "INSERT INTO `CustomTableName_1` (`text`) VALUES (?);", v.Values()
}
func (v CustomTableName1) GetText() sequel.ColumnValuer[string] {
	return sequel.Column("`text`", v.Text, func(val string) driver.Value { return string(val) })
}

func (CustomTableName2) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (CustomTableName2) TableName() string {
	return "`table_2`"
}
func (CustomTableName2) ColumnNames() []string {
	return []string{"`text`"}
}
func (v CustomTableName2) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName2) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (CustomTableName2) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v CustomTableName2) InsertOneStmt() (string, []any) {
	return "INSERT INTO `table_2` (`text`) VALUES (?);", v.Values()
}
func (v CustomTableName2) GetText() sequel.ColumnValuer[string] {
	return sequel.Column("`text`", v.Text, func(val string) driver.Value { return string(val) })
}

func (CustomTableName3) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (CustomTableName3) TableName() string {
	return "`table_3`"
}
func (CustomTableName3) ColumnNames() []string {
	return []string{"`text`"}
}
func (v CustomTableName3) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName3) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (CustomTableName3) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v CustomTableName3) InsertOneStmt() (string, []any) {
	return "INSERT INTO `table_3` (`text`) VALUES (?);", v.Values()
}
func (v CustomTableName3) GetText() sequel.ColumnValuer[string] {
	return sequel.Column("`text`", v.Text, func(val string) driver.Value { return string(val) })
}
