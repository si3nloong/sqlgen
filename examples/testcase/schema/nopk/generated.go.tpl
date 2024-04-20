package nopk

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Customer) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (name VARCHAR(255) NOT NULL,age TINYINT UNSIGNED NOT NULL,married TINYINT NOT NULL);"
}
func (Customer) AlterTableStmt() string {
	return "ALTER TABLE customer MODIFY name VARCHAR(255) NOT NULL,MODIFY age TINYINT UNSIGNED NOT NULL AFTER name,MODIFY married TINYINT NOT NULL AFTER age;"
}
func (Customer) TableName() string {
	return "customer"
}
func (v Customer) InsertOneStmt() string {
	return "INSERT INTO customer (name,age,married) VALUES (?,?,?);"
}
func (Customer) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Customer) Columns() []string {
	return []string{"name", "age", "married"}
}
func (v Customer) Values() []any {
	return []any{string(v.Name), int64(v.Age), bool(v.Married)}
}
func (v *Customer) Addrs() []any {
	return []any{types.String(&v.Name), types.Integer(&v.Age), types.Bool(&v.Married)}
}
func (v Customer) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(vi string) driver.Value { return string(vi) })
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(vi uint8) driver.Value { return int64(vi) })
}
func (v Customer) GetMarried() sequel.ColumnValuer[bool] {
	return sequel.Column("married", v.Married, func(vi bool) driver.Value { return bool(vi) })
}
