package int64

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (Model) TableName() string {
	return "model"
}
func (Model) HasPK()      {}
func (Model) IsAutoIncr() {}
func (v *Model) ScanAutoIncr(val int64) error {
	v.ID = int64(val)
	return nil
}
func (v Model) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (Model) Columns() []string {
	return []string{"id"} // 1
}
func (v *Model) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
	}
}
func (v Model) FindOneByPKStmt() (string, []any) {
	return "SELECT `id` FROM `model` WHERE `id` = ? LIMIT 1;", []any{v.ID}
}
func (v Model) IDValue() any {
	return v.ID
}
func (v Model) ColumnID() sequel.ColumnClause[int64] {
	return sequel.BasicColumn("id", v.ID)
}
