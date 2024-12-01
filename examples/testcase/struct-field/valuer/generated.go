package valuer

import (
	"database/sql/driver"

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
		v.ID,            // 0 - id
		v.Value,         // 1 - value
		v.GetPtrValue(), // 2 - ptr_value
		v.N,             // 3 - n
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
func (v B) GetID() driver.Value {
	return v.ID
}
func (v B) GetValue() driver.Value {
	return v.Value
}
func (v B) GetPtrValue() driver.Value {
	if v.PtrValue != nil {
		return *v.PtrValue
	}
	return nil
}
func (v B) GetN() driver.Value {
	return v.N
}
