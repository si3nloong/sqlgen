package embedded

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (B) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "id", Definition: "id BIGINT NOT NULL DEFAULT 0"},
			{Name: "name", Definition: "name VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "z", Definition: "z BOOL NOT NULL DEFAULT false"},
			{Name: "created", Definition: "created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"},
			{Name: "ok", Definition: "ok BOOL NOT NULL DEFAULT false"},
		},
	}
}
func (B) TableName() string {
	return "b"
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
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,name,z,created,ok) VALUES (?,?,?,?,?);", v.Values()
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
