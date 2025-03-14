package pkautoincr

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
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
	return []string{"name", "f", "id", "n"} // 4
}
func (v Model) Values() []any {
	return []any{
		(string)(v.Name), // 0 - name
		(bool)(v.F),      // 1 - f
		v.N,              // 3 - n
	}
}
func (v *Model) Addrs() []any {
	return []any{
		encoding.StringScanner[LongText](&v.Name), // 0 - name
		encoding.BoolScanner[Flag](&v.F),          // 1 - f
		encoding.UintScanner[uint](&v.ID),         // 2 - id
		&v.N,                                      // 3 - n
	}
}
func (Model) InsertColumns() []string {
	return []string{"name", "f", "n"} // 3
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO AutoIncrPK (name,f,n) VALUES (?,?,?);", []any{(string)(v.Name), (bool)(v.F), v.N}
}
func (v Model) FindOneByPKStmt() (string, []any) {
	return "SELECT name,f,id,n FROM AutoIncrPK WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Model) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE AutoIncrPK SET name = ?,f = ?,n = ? WHERE id = ?;", []any{(string)(v.Name), (bool)(v.F), v.N, (int64)(v.ID)}
}
func (v Model) NameValue() driver.Value {
	return (string)(v.Name)
}
func (v Model) FValue() driver.Value {
	return (bool)(v.F)
}
func (v Model) IDValue() driver.Value {
	return (int64)(v.ID)
}
func (v Model) NValue() driver.Value {
	return v.N
}
func (v Model) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column("name", v.Name, func(val LongText) driver.Value {
		return (string)(val)
	})
}
func (v Model) GetF() sequel.ColumnValuer[Flag] {
	return sequel.Column("f", v.F, func(val Flag) driver.Value {
		return (bool)(val)
	})
}
func (v Model) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column("id", v.ID, func(val uint) driver.Value {
		return (int64)(val)
	})
}
func (v Model) GetN() sequel.ColumnValuer[int64] {
	return sequel.Column("n", v.N, func(val int64) driver.Value {
		return val
	})
}
