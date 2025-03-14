package main

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
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
func (v User) IDValue() driver.Value {
	return v.ID
}
func (v User) NameValue() driver.Value {
	return v.Name
}
func (v User) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value {
		return val
	})
}
func (v User) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return val
	})
}
