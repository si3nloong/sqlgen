package embedded

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BIGINT NOT NULL,`name` VARCHAR(255) NOT NULL,`z` TINYINT NOT NULL,`created` DATETIME NOT NULL,`ok` TINYINT NOT NULL);"
}
func (B) TableName() string {
	return "b"
}
func (B) InsertOneStmt() string {
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
	return sequel.Column("id", v.a.ID, func(val int64) driver.Value { return int64(val) })
}
func (v B) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.a.Name, func(val string) driver.Value { return string(val) })
}
func (v B) GetZ() sequel.ColumnValuer[bool] {
	return sequel.Column("z", v.a.Z, func(val bool) driver.Value { return bool(val) })
}
func (v B) GetCreated() sequel.ColumnValuer[time.Time] {
	return sequel.Column("created", v.ts.Created, func(val time.Time) driver.Value { return time.Time(val) })
}
func (v B) GetOK() sequel.ColumnValuer[bool] {
	return sequel.Column("ok", v.ts.OK, func(val bool) driver.Value { return bool(val) })
}
