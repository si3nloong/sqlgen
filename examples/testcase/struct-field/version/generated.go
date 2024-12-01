package version

import (
	"database/sql/driver"
)

func (Version) TableName() string {
	return "version"
}
func (Version) HasPK() {}
func (v Version) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (Version) Columns() []string {
	return []string{"id"} // 1
}
func (v Version) Values() []any {
	return []any{
		v.ID, // 0 - id
	}
}
func (v *Version) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
	}
}
func (Version) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v Version) InsertOneStmt() (string, []any) {
	return "INSERT INTO version (id) VALUES (?);", v.Values()
}
func (v Version) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM version WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v Version) GetID() driver.Value {
	return v.ID
}
