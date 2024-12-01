package tabler

import (
	"database/sql/driver"
)

func (A) HasPK() {}
func (v A) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (A) Columns() []string {
	return []string{"id", "name"}
}
func (v A) Values() []any {
	return []any{v.ID, v.Name}
}
func (v *A) Addrs() []any {
	return []any{&v.ID, &v.Name}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (id,name) VALUES (?,?);", v.Values()
}
func (v A) FindOneByPKStmt() (string, []any) {
	return "SELECT id,name FROM " + v.TableName() + " WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v A) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE " + v.TableName() + " SET name = ? WHERE id = ?;", []any{v.Name, v.ID}
}
func (v A) GetID() driver.Value {
	return v.ID
}
func (v A) GetName() driver.Value {
	return v.Name
}

func (Model) Columns() []string {
	return []string{"name"}
}
func (v Model) Values() []any {
	return []any{v.Name}
}
func (v *Model) Addrs() []any {
	return []any{&v.Name}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (name) VALUES (?);", v.Values()
}
func (v Model) GetName() driver.Value {
	return v.Name
}
