package version

import (
	"database/sql"
	"database/sql/driver"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/si3nloong/sqlgen/sequel"
)

func (v Version) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (id VARCHAR(36) NOT NULL,PRIMARY KEY (id));"
}
func (Version) AlterTableStmt() string {
	return "ALTER TABLE version MODIFY id VARCHAR(36) NOT NULL;"
}
func (Version) TableName() string {
	return "version"
}
func (v Version) InsertOneStmt() string {
	return "INSERT INTO version (id) VALUES (?);"
}
func (Version) InsertVarQuery() string {
	return "(?)"
}
func (Version) Columns() []string {
	return []string{"id"}
}
func (v Version) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, (driver.Valuer)(v.ID)
}
func (v Version) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *Version) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (v Version) GetID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(vi uuid.UUID) driver.Value { return (driver.Valuer)(vi) })
}
