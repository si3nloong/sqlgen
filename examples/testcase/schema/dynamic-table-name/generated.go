package tabler

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (A) HasPK() {}
func (v A) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (A) Columns() []string {
	return []string{"id", "name"}
}
func (v A) Values() []any {
	return []any{(int64)(v.ID), (string)(v.Name)}
}
func (v *A) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name)}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (id,name) VALUES (?,?);", v.Values()
}
func (v A) FindOneByPKStmt() (string, []any) {
	return "SELECT id,name FROM " + v.TableName() + " WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v A) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE " + v.TableName() + " SET name = ? WHERE id = ?;", []any{(string)(v.Name), (int64)(v.ID)}
}
func (v A) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value {
		return (int64)(val)
	})
}
func (v A) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return (string)(val)
	})
}

func (Model) Columns() []string {
	return []string{"name"}
}
func (v Model) Values() []any {
	return []any{(string)(v.Name)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name)}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (name) VALUES (?);", v.Values()
}
func (v Model) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return (string)(val)
	})
}
