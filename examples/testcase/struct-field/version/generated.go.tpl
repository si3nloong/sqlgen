package version

import (
	"database/sql"
	"database/sql/driver"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/si3nloong/sqlgen/sequel"
)

func (v Version) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `version` (`id` VARCHAR(36),PRIMARY KEY (`id`));"
}
func (Version) TableName() string {
	return "version"
}
func (Version) InsertOneStmt() string {
	return "INSERT INTO version (id) VALUES (?);"
}
func (Version) InsertVarQuery() string {
	return "(?)"
}
func (Version) Columns() []string {
	return []string{"id"}
}
func (Version) HasPK() {}
func (v Version) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{0}, []any{(driver.Valuer)(v.ID)}
}
func (v Version) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *Version) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (v Version) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
