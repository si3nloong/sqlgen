package int16

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Model) TableName() string {
	return "model"
}
func (Model) HasPK()      {}
func (Model) IsAutoIncr() {}
func (v *Model) ScanAutoIncr(val int64) error {
	v.ID = int16(val)
	return nil
}
func (v Model) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (Model) Columns() []string {
	return []string{"id"} // 1
}
func (v *Model) Addrs() []any {
	return []any{
		encoding.Int16Scanner[int16](&v.ID), // 0 - id
	}
}
func (v Model) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM model WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Model) IDValue() driver.Value {
	return (int64)(v.ID)
}
func (v Model) GetID() sequel.ColumnValuer[int16] {
	return sequel.Column("id", v.ID, func(val int16) driver.Value {
		return (int64)(val)
	})
}
