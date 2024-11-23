package int

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Model) TableName() string {
	return "model"
}
func (Model) HasPK()      {}
func (Model) IsAutoIncr() {}
func (v *Model) ScanAutoIncr(val int64) error {
	v.ID = int(val)
	return nil
}
func (v Model) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (Model) Columns() []string {
	return []string{"id"}
}
func (v Model) Values() []any {
	return []any{(int64)(v.ID)}
}
func (v *Model) Addrs() []any {
	return []any{types.Integer(&v.ID)}
}
func (v Model) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM model WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Model) GetID() sequel.ColumnValuer[int] {
	return sequel.Column("id", v.ID, func(val int) driver.Value { return (int64)(val) })
}
