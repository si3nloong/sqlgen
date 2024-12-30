package nopk

import (
	"database/sql/driver"

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
	return "INSERT INTO customer (name,age,married) VALUES (?,?,?);", v.Values()
}
func (v Customer) NameValue() driver.Value {
	return v.Name
}
func (v Customer) AgeValue() driver.Value {
	return (int64)(v.Age)
}
func (v Customer) MarriedValue() driver.Value {
	return v.Married
}
func (v Customer) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return val
	})
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) driver.Value {
		return (int64)(val)
	})
}
func (v Customer) GetMarried() sequel.ColumnValuer[bool] {
	return sequel.Column("married", v.Married, func(val bool) driver.Value {
		return val
	})
}
