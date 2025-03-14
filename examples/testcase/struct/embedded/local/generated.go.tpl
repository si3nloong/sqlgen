package embedded

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "name", "z", "created", "ok"} // 5
}
func (v B) Values() []any {
	return []any{
		v.a.ID,       // 0 - id
		v.a.Name,     // 1 - name
		v.a.Z,        // 2 - z
		v.ts.Created, // 3 - created
		v.ts.OK,      // 4 - ok
	}
}
func (v *B) Addrs() []any {
	return []any{
		&v.a.ID,       // 0 - id
		&v.a.Name,     // 1 - name
		&v.a.Z,        // 2 - z
		&v.ts.Created, // 3 - created
		&v.ts.OK,      // 4 - ok
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?)" // 5
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,name,z,created,ok) VALUES (?,?,?,?,?);", v.Values()
}
func (v B) IDValue() driver.Value {
	return v.a.ID
}
func (v B) NameValue() driver.Value {
	return v.a.Name
}
func (v B) ZValue() driver.Value {
	return v.a.Z
}
func (v B) CreatedValue() driver.Value {
	return v.ts.Created
}
func (v B) OKValue() driver.Value {
	return v.ts.OK
}
func (v B) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.a.ID, func(val int64) driver.Value {
		return val
	})
}
func (v B) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.a.Name, func(val string) driver.Value {
		return val
	})
}
func (v B) GetZ() sequel.ColumnValuer[bool] {
	return sequel.Column("z", v.a.Z, func(val bool) driver.Value {
		return val
	})
}
func (v B) GetCreated() sequel.ColumnValuer[time.Time] {
	return sequel.Column("created", v.ts.Created, func(val time.Time) driver.Value {
		return val
	})
}
func (v B) GetOK() sequel.ColumnValuer[bool] {
	return sequel.Column("ok", v.ts.OK, func(val bool) driver.Value {
		return val
	})
}
