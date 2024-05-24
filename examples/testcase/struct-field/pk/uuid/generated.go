package main

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` VARCHAR(36),`name` VARCHAR(255) NOT NULL);"
}
func (User) TableName() string {
	return "user"
}
func (User) InsertOneStmt() string {
	return "INSERT INTO user (id,name) VALUES (?,?);"
}
func (User) InsertVarQuery() string {
	return "(?,?)"
}
func (User) Columns() []string {
	return []string{"id", "name"}
}
func (v User) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.Name)}
}
func (v *User) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID), types.String(&v.Name)}
}
func (v User) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
func (v User) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value { return string(val) })
}
