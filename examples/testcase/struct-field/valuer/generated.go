package valuer

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (id BIGINT NOT NULL,value VARCHAR(255) NOT NULL,n VARCHAR(255) NOT NULL);"
}
func (B) AlterTableStmt() string {
	return "ALTER TABLE b MODIFY id BIGINT NOT NULL,MODIFY value VARCHAR(255) NOT NULL AFTER id,MODIFY n VARCHAR(255) NOT NULL AFTER value;"
}
func (B) TableName() string {
	return "b"
}
func (v B) InsertOneStmt() string {
	return "INSERT INTO b (id,value,n) VALUES (?,?,?);"
}
func (B) InsertVarQuery() string {
	return "(?,?,?)"
}
func (B) Columns() []string {
	return []string{"id", "value", "n"}
}
func (v B) Values() []any {
	return []any{int64(v.ID), (driver.Valuer)(v.Value), string(v.N)}
}
func (v *B) Addrs() []any {
	return []any{types.Integer(&v.ID), &v.Value, types.String(&v.N)}
}
func (v B) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("id", v.ID, func(vi int64) driver.Value { return int64(vi) })
}
func (v B) GetValue() sequel.ColumnValuer[anyType] {
	return sequel.Column[anyType]("value", v.Value, func(vi anyType) driver.Value { return (driver.Valuer)(vi) })
}
func (v B) GetN() sequel.ColumnValuer[string] {
	return sequel.Column[string]("n", v.N, func(vi string) driver.Value { return string(vi) })
}
