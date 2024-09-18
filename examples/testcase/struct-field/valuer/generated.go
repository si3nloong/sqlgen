package valuer

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "value", "ptr_value", "n"}
}
func (v B) Values() []any {
	return []any{int64(v.ID), (driver.Valuer)(v.Value), (driver.Valuer)(v.PtrValue), string(v.N)}
}
func (v *B) Addrs() []any {
	addrs := make([]any, 4)
	addrs[0] = types.Integer(&v.ID)
	addrs[1] = types.JSONUnmarshaler(&v.Value)
	if v.PtrValue == nil {
		v.PtrValue = new(anyType)
	}
	addrs[2] = types.JSONUnmarshaler(&v.PtrValue)
	addrs[3] = types.String(&v.N)
	return addrs
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,value,ptr_value,n) VALUES (?,?,?,?);", v.Values()
}
func (v B) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v B) GetValue() sequel.ColumnValuer[anyType] {
	return sequel.Column("value", v.Value, func(val anyType) driver.Value { return (driver.Valuer)(val) })
}
func (v B) GetPtrValue() sequel.ColumnValuer[*anyType] {
	return sequel.Column("ptr_value", v.PtrValue, func(val *anyType) driver.Value { return (driver.Valuer)(val) })
}
func (v B) GetN() sequel.ColumnValuer[string] {
	return sequel.Column("n", v.N, func(val string) driver.Value { return string(val) })
}
