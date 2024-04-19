package embedded

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (id BIGINT NOT NULL,name VARCHAR(255) NOT NULL,z TINYINT NOT NULL,created DATETIME NOT NULL,ok TINYINT NOT NULL);"
}
func (B) AlterTableStmt() string {
	return "ALTER TABLE b MODIFY id BIGINT NOT NULL,MODIFY name VARCHAR(255) NOT NULL AFTER id,MODIFY z TINYINT NOT NULL AFTER name,MODIFY created DATETIME NOT NULL AFTER z,MODIFY ok TINYINT NOT NULL AFTER created;"
}
func (B) TableName() string {
	return "b"
}
func (v B) InsertOneStmt() string {
	return "INSERT INTO b (id,name,z,created,ok) VALUES (?,?,?,?,?);"
}
func (B) InsertVarQuery() string {
	return "(?,?,?,?,?)"
}
func (B) Columns() []string {
	return []string{"id", "name", "z", "created", "ok"}
}
func (v B) Values() []any {
	return []any{int64(v.a.ID), string(v.a.Name), bool(v.a.Z), time.Time(v.ts.Created), bool(v.ts.OK)}
}
func (v *B) Addrs() []any {
	return []any{types.Integer(&v.a.ID), types.String(&v.a.Name), types.Bool(&v.a.Z), (*time.Time)(&v.ts.Created), types.Bool(&v.ts.OK)}
}
func (v B) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("id", v.a.ID, func(vi int64) driver.Value { return int64(vi) })
}
func (v B) GetName() sequel.ColumnValuer[string] {
	return sequel.Column[string]("name", v.a.Name, func(vi string) driver.Value { return string(vi) })
}
func (v B) GetZ() sequel.ColumnValuer[bool] {
	return sequel.Column[bool]("z", v.a.Z, func(vi bool) driver.Value { return bool(vi) })
}
func (v B) GetCreated() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("created", v.ts.Created, func(vi time.Time) driver.Value { return time.Time(vi) })
}
func (v B) GetOK() sequel.ColumnValuer[bool] {
	return sequel.Column[bool]("ok", v.ts.OK, func(vi bool) driver.Value { return bool(vi) })
}
