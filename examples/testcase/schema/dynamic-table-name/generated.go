package tabler

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`name` VARCHAR(255) NOT NULL);"
}
func (Model) InsertVarQuery() string {
	return "(?)"
}
func (Model) Columns() []string {
	return []string{"name"}
}
func (v Model) Values() []any {
	return []any{string(v.Name)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name)}
}
func (v Model) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value { return string(val) })
}

func (v A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BIGINT NOT NULL,`name` VARCHAR(255) NOT NULL,PRIMARY KEY (`id`));"
}
func (A) InsertVarQuery() string {
	return "(?,?)"
}
func (A) Columns() []string {
	return []string{"id", "name"}
}
func (A) HasPK() {}
func (v A) PK() (string, int, any) {
	return "id", 0, int64(v.ID)
}
func (v A) Values() []any {
	return []any{int64(v.ID), string(v.Name)}
}
func (v *A) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name)}
}
func (v A) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v A) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value { return string(val) })
}
