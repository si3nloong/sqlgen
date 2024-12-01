package readonly

import (
	"database/sql/driver"
)

func (Model) TableName() string {
	return "model"
}
func (Model) InsertColumns() []string {
	return []string{"a", "b"} // 2
}
func (Model) Columns() []string {
	return []string{"a", "b", "read_only"} // 3
}
func (v Model) Values() []any {
	return []any{
		v.A, // 0 - a
		v.B, // 1 - b
	}
}
func (v *Model) Addrs() []any {
	return []any{
		&v.A,        // 0 - a
		&v.B,        // 1 - b
		&v.ReadOnly, // 2 - read_only
	}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO model (a,b) VALUES (?,?);", []any{v.A, v.B}
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
