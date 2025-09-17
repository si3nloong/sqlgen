package nopk

import (
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Customer) TableName() string {
	return "customer"
}
func (Customer) Columns() []string {
	return []string{"name", "age", "married"} // 3
}
func (v Customer) Values() []any {
	return []any{
		v.Name,         // 0 - name
		(int64)(v.Age), // 1 - age
		v.Married,      // 2 - married
	}
}
func (v *Customer) Addrs() []any {
	return []any{
		&v.Name,                              // 0 - name
		encoding.Uint8Scanner[uint8](&v.Age), // 1 - age
		&v.Married,                           // 2 - married
	}
}
func (Customer) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v Customer) InsertOneStmt() (string, []any) {
	return "INSERT INTO `customer` (`name`,`age`,`married`) VALUES (?,?,?);", v.Values()
}
func (v Customer) NameValue() any {
	return v.Name
}
func (v Customer) AgeValue() any {
	return (int64)(v.Age)
}
func (v Customer) MarriedValue() any {
	return v.Married
}
func (v Customer) ColumnName() sequel.ColumnClause[string] {
	return sequel.BasicColumn("name", v.Name)
}
func (v Customer) ColumnAge() sequel.ColumnConvertClause[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) any {
		return (int64)(val)
	})
}
func (v Customer) ColumnMarried() sequel.ColumnClause[bool] {
	return sequel.BasicColumn("married", v.Married)
}
