package nopk

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Customer) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (Customer) TableName() string {
	return "`customer`"
}
func (Customer) ColumnNames() []string {
	return []string{"`name`", "`age`", "`married`"}
}
func (v Customer) Values() []any {
	return []any{string(v.Name), int64(v.Age), bool(v.Married)}
}
func (v *Customer) Addrs() []any {
	return []any{types.String(&v.Name), types.Integer(&v.Age), types.Bool(&v.Married)}
}
func (Customer) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Customer) InsertOneStmt() (string, []any) {
	return "INSERT INTO `customer` (`name`,`age`,`married`) VALUES (?,?,?);", v.Values()
}
func (v Customer) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("`name`", v.Name, func(val string) driver.Value { return string(val) })
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("`age`", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v Customer) GetMarried() sequel.ColumnValuer[bool] {
	return sequel.Column("`married`", v.Married, func(val bool) driver.Value { return bool(val) })
}
