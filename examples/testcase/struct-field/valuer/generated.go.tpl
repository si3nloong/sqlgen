package valuer

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "value", "ptr_value", "n"}
}
func (v B) Values() []any {
	return []any{v.ID, (driver.Valuer)(v.Value), v.GetPtrValue(), v.N}
}
func (v *B) Addrs() []any {
	addrs := make([]any, 4)
	addrs[0] = &v.ID
	addrs[1] = encoding.JSONScanner(&v.Value)
	if v.PtrValue == nil {
		v.PtrValue = new(anyType)
	}
	addrs[2] = encoding.JSONScanner(&v.PtrValue)
	addrs[3] = &v.N
	return addrs
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,value,ptr_value,n) VALUES (?,?,?,?);", v.Values()
}
func (v B) GetID() driver.Value {
	return v.ID
}
func (v B) GetValue() driver.Value {
	return (driver.Valuer)(v.Value)
}
func (v B) GetPtrValue() driver.Value {
	if v.PtrValue != nil {
		return (driver.Valuer)(*v.PtrValue)
	}
	return nil
}
func (v B) GetN() driver.Value {
	return v.N
}
