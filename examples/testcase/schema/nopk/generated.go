package nopk

import (
	"database/sql/driver"

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
func (v Customer) GetName() driver.Value {
	return v.Name
}
func (v Customer) GetAge() driver.Value {
	return (int64)(v.Age)
}
func (v Customer) GetMarried() driver.Value {
	return v.Married
}
