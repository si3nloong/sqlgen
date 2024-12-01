package embedded

import (
	"database/sql/driver"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "name", "z", "created", "ok"} // 5
}
func (v B) Values() []any {
	return []any{
		v.a.ID,       // 0 - id
		v.a.Name,     // 1 - name
		v.a.Z,        // 2 - z
		v.ts.Created, // 3 - created
		v.ts.OK,      // 4 - ok
	}
}
func (v *B) Addrs() []any {
	return []any{
		&v.a.ID,       // 0 - id
		&v.a.Name,     // 1 - name
		&v.a.Z,        // 2 - z
		&v.ts.Created, // 3 - created
		&v.ts.OK,      // 4 - ok
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?)" // 5
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,name,z,created,ok) VALUES (?,?,?,?,?);", v.Values()
}
func (v B) GetID() driver.Value {
	return v.a.ID
}
func (v B) GetName() driver.Value {
	return v.a.Name
}
func (v B) GetZ() driver.Value {
	return v.a.Z
}
func (v B) GetCreated() driver.Value {
	return v.ts.Created
}
func (v B) GetOK() driver.Value {
	return v.ts.OK
}
