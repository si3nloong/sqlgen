package readonly

import (
	"database/sql/driver"
)

func (Model) TableName() string {
	return "model"
}
func (Model) Columns() []string {
	return []string{"a", "b", "read_only"}
}
func (v Model) Values() []any {
	return []any{v.A, v.B, v.ReadOnly}
}
func (v *Model) Addrs() []any {
	return []any{&v.A, &v.B, &v.ReadOnly}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO model (a,b,read_only) VALUES (?,?,?);", v.Values()
}
func (v Model) GetA() driver.Value {
	return v.A
}
func (v Model) GetB() driver.Value {
	return v.B
}
func (v Model) GetReadOnly() driver.Value {
	return v.ReadOnly
}
