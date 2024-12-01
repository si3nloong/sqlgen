package main

import (
	"database/sql/driver"
)

func (User) TableName() string {
	return "user"
}
func (User) Columns() []string {
	return []string{"id", "name"} // 2
}
func (v User) Values() []any {
	return []any{
		v.ID,   // 0 - id
		v.Name, // 1 - name
	}
}
func (v *User) Addrs() []any {
	return []any{
		&v.ID,   // 0 - id
		&v.Name, // 1 - name
	}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO user (id,name) VALUES (?,?);", v.Values()
}
func (v User) GetID() driver.Value {
	return v.ID
}
func (v User) GetName() driver.Value {
	return v.Name
}
