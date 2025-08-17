package version

import (
	"database/sql/driver"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/si3nloong/sqlgen/sequel"
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
func (v Version) IDValue() driver.Value {
	return v.ID
}
func (v Version) ColumnID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value {
		return val
	})
}
