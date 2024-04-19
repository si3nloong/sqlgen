package tablename

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v CustomTableName1) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (text VARCHAR(255) NOT NULL);"
}
func (CustomTableName1) AlterTableStmt() string {
	return "ALTER TABLE CustomTableName_1 MODIFY text VARCHAR(255) NOT NULL;"
}
func (CustomTableName1) TableName() string {
	return "CustomTableName_1"
}
func (v CustomTableName1) InsertOneStmt() string {
	return "INSERT INTO CustomTableName_1 (text) VALUES (?);"
}
func (CustomTableName1) InsertVarQuery() string {
	return "(?)"
}
func (CustomTableName1) Columns() []string {
	return []string{"text"}
}
func (v CustomTableName1) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName1) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (v CustomTableName1) GetText() sequel.ColumnValuer[string] {
	return sequel.Column[string]("text", v.Text, func(vi string) driver.Value { return string(vi) })
}

func (v CustomTableName2) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (text VARCHAR(255) NOT NULL);"
}
func (CustomTableName2) AlterTableStmt() string {
	return "ALTER TABLE table_2 MODIFY text VARCHAR(255) NOT NULL;"
}
func (CustomTableName2) TableName() string {
	return "table_2"
}
func (v CustomTableName2) InsertOneStmt() string {
	return "INSERT INTO table_2 (text) VALUES (?);"
}
func (CustomTableName2) InsertVarQuery() string {
	return "(?)"
}
func (CustomTableName2) Columns() []string {
	return []string{"text"}
}
func (v CustomTableName2) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName2) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (v CustomTableName2) GetText() sequel.ColumnValuer[string] {
	return sequel.Column[string]("text", v.Text, func(vi string) driver.Value { return string(vi) })
}

func (v CustomTableName3) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (text VARCHAR(255) NOT NULL);"
}
func (CustomTableName3) AlterTableStmt() string {
	return "ALTER TABLE table_3 MODIFY text VARCHAR(255) NOT NULL;"
}
func (CustomTableName3) TableName() string {
	return "table_3"
}
func (v CustomTableName3) InsertOneStmt() string {
	return "INSERT INTO table_3 (text) VALUES (?);"
}
func (CustomTableName3) InsertVarQuery() string {
	return "(?)"
}
func (CustomTableName3) Columns() []string {
	return []string{"text"}
}
func (v CustomTableName3) Values() []any {
	return []any{string(v.Text)}
}
func (v *CustomTableName3) Addrs() []any {
	return []any{types.String(&v.Text)}
}
func (v CustomTableName3) GetText() sequel.ColumnValuer[string] {
	return sequel.Column[string]("text", v.Text, func(vi string) driver.Value { return string(vi) })
}
