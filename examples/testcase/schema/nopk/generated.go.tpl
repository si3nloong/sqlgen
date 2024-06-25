package nopk

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Customer) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `customer` (`name` VARCHAR(255) NOT NULL,`age` TINYINT UNSIGNED NOT NULL,`married` BOOL NOT NULL);"
}
func (Customer) TableName() string {
	return "customer"
}
func (Customer) InsertOneStmt() string {
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
	return sequel.Column("name", v.Name, func(val string) driver.Value { return string(val) })
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v Customer) GetMarried() sequel.ColumnValuer[bool] {
	return sequel.Column("married", v.Married, func(val bool) driver.Value { return bool(val) })
}
