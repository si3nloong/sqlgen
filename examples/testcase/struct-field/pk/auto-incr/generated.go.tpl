package pkautoincr

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Model) TableName() string {
	return "AutoIncrPK"
}
func (Model) HasPK()      {}
func (Model) IsAutoIncr() {}
func (v *Model) ScanAutoIncr(val int64) error {
	v.ID = uint(val)
	return nil
}
func (v Model) PK() (string, int, any) {
	return "id", 2, (int64)(v.ID)
}
func (Model) Columns() []string {
	return []string{"name", "f", "id", "n"}
}
func (v Model) Values() []any {
	return []any{(string)(v.Name), (bool)(v.F), (int64)(v.ID), v.N}
}
func (v *Model) Addrs() []any {
	return []any{encoding.StringScanner[LongText](&v.Name), encoding.BoolScanner[Flag](&v.F), encoding.UintScanner[uint](&v.ID), &v.N}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO AutoIncrPK (name,f,n) VALUES (?,?,?);", []any{v.GetName(), v.GetF(), v.GetN()}
}
func (v Model) FindOneByPKStmt() (string, []any) {
	return "SELECT name,f,id,n FROM AutoIncrPK WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Model) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE AutoIncrPK SET name = ?,f = ?,n = ? WHERE id = ?;", []any{(string)(v.Name), (bool)(v.F), v.N, (int64)(v.ID)}
}
func (v Model) GetName() driver.Value {
	return (string)(v.Name)
}
func (v Model) GetF() driver.Value {
	return (bool)(v.F)
}
func (v Model) GetID() driver.Value {
	return (int64)(v.ID)
}
func (v Model) GetN() driver.Value {
	return v.N
}
