package tabler

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
)

func (A) HasPK() {}
func (v A) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (A) Columns() []string {
	return []string{"id", "name"} // 2
}
func (v A) Values() []any {
	return []any{
		v.ID,   // 0 - id
		v.Name, // 1 - name
	}
}
func (v *A) Addrs() []any {
	return []any{
		&v.ID,   // 0 - id
		&v.Name, // 1 - name
	}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
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
func (v A) IDValue() driver.Value {
	return v.ID
}
func (v A) NameValue() driver.Value {
	return v.Name
}
func (v A) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value {
		return val
	})
}
func (v A) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return val
	})
}

func (Model) Columns() []string {
	return []string{"name"} // 1
}
func (v Model) Values() []any {
	return []any{
		v.Name, // 0 - name
	}
}
func (v *Model) Addrs() []any {
	return []any{
		&v.Name, // 0 - name
	}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (name) VALUES (?);", v.Values()
}
func (v Model) NameValue() driver.Value {
	return v.Name
}
func (v Model) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return val
	})
}
