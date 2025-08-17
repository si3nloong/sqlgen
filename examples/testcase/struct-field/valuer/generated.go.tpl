package valuer

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "value", "ptr_value", "n"} // 4
}
func (v B) Values() []any {
	return []any{
		v.ID,              // 0 - id
		v.Value,           // 1 - value
		v.PtrValueValue(), // 2 - ptr_value
		v.N,               // 3 - n
	}
}
func (v *B) Addrs() []any {
	if v.PtrValue == nil {
		v.PtrValue = new(anyType)
	}
	return []any{
		&v.ID,                             // 0 - id
		encoding.JSONScanner(&v.Value),    // 1 - value
		encoding.JSONScanner(&v.PtrValue), // 2 - ptr_value
		&v.N,                              // 3 - n
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?)" // 4
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,value,ptr_value,n) VALUES (?,?,?,?);", v.Values()
}
func (v B) IDValue() driver.Value {
	return v.ID
}
func (v B) ValueValue() driver.Value {
	return v.Value
}
func (v B) PtrValueValue() driver.Value {
	if v.PtrValue != nil {
		return *v.PtrValue
	}
	return nil
}
func (v B) NValue() driver.Value {
	return v.N
}
func (v B) ColumnID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value {
		return val
	})
}
func (v B) ColumnValue() sequel.ColumnValuer[anyType] {
	return sequel.Column("value", v.Value, func(val anyType) driver.Value {
		return val
	})
}
func (v B) ColumnPtrValue() sequel.ColumnValuer[*anyType] {
	return sequel.Column("ptr_value", v.PtrValue, func(val *anyType) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v B) ColumnN() sequel.ColumnValuer[string] {
	return sequel.Column("n", v.N, func(val string) driver.Value {
		return val
	})
}
