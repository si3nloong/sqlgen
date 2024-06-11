package pkautoincr

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `AutoIncrPK` (`name` VARCHAR(255) NOT NULL,`f` TINYINT NOT NULL,`id` INTEGER NOT NULL AUTO_INCREMENT,`n` BIGINT NOT NULL,PRIMARY KEY (`id`));"
}
func (Model) TableName() string {
	return "AutoIncrPK"
}
func (Model) InsertOneStmt() string {
	return "INSERT INTO AutoIncrPK (name,f,n) VALUES (?,?,?);"
}
func (Model) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Model) Columns() []string {
	return []string{"name", "f", "id", "n"}
}
func (Model) HasPK()      {}
func (Model) IsAutoIncr() {}
func (v Model) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{2}, []any{int64(v.ID)}
}
func (Model) FindByPKStmt() string {
	return "SELECT name,f,id,n FROM AutoIncrPK WHERE id = ? LIMIT 1;"
}
func (Model) UpdateByPKStmt() string {
	return "UPDATE AutoIncrPK SET name = ?,f = ?,n = ? WHERE id = ? LIMIT 1;"
}
func (v Model) Values() []any {
	return []any{string(v.Name), bool(v.F), int64(v.ID), int64(v.N)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name), types.Bool(&v.F), types.Integer(&v.ID), types.Integer(&v.N)}
}
func (v Model) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column("name", v.Name, func(val LongText) driver.Value { return string(val) })
}
func (v Model) GetF() sequel.ColumnValuer[Flag] {
	return sequel.Column("f", v.F, func(val Flag) driver.Value { return bool(val) })
}
func (v Model) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column("id", v.ID, func(val uint) driver.Value { return int64(val) })
}
func (v Model) GetN() sequel.ColumnValuer[int64] {
	return sequel.Column("n", v.N, func(val int64) driver.Value { return int64(val) })
}
