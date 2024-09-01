package version

import (
	"database/sql"
	"database/sql/driver"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/si3nloong/sqlgen/sequel"
)

func (Version) TableName() string {
	return "version"
}
func (Version) HasPK() {}
func (v Version) PK() (string, int, any) {
	return "id", 0, (driver.Valuer)(v.ID)
}
func (Version) Columns() []string {
	return []string{"id"}
}
func (v Version) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *Version) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (Version) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v Version) InsertOneStmt() (string, []any) {
	return "INSERT INTO version (id) VALUES (?);", v.Values()
}
func (v Version) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM version WHERE id = ? LIMIT 1;", []any{(driver.Valuer)(v.ID)}
}
func (v Version) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
