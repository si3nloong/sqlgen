package readonly

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
)

func (Model) TableName() string {
	return "model"
}
func (Model) Columns() []string {
	return []string{"a", "b", "read_only"} // 3
}
func (v Model) Values() []any {
	return []any{
		v.A, // 0 - a
		v.B, // 1 - b
	}
}
func (v *Model) Addrs() []any {
	return []any{
		&v.A,        // 0 - a
		&v.B,        // 1 - b
		&v.ReadOnly, // 2 - read_only
	}
}
func (Model) InsertColumns() []string {
	return []string{"a", "b"} // 2
}
func (Model) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO `model` (`a`,`b`) VALUES (?,?);", []any{v.A, v.B}
}
func (v Model) AValue() any {
	return v.A
}
func (v Model) BValue() any {
	return v.B
}
func (v Model) ReadOnlyValue() any {
	return v.ReadOnly
}
func (v Model) ColumnA() sequel.ColumnValuer[string] {
	return sequel.Column("a", v.A, func(val string) driver.Value {
		return val
	})
}
func (v Model) ColumnB() sequel.ColumnValuer[bool] {
	return sequel.Column("b", v.B, func(val bool) driver.Value {
		return val
	})
}
func (v Model) ColumnReadOnly() sequel.ColumnValuer[string] {
	return sequel.Column("read_only", v.ReadOnly, func(val string) driver.Value {
		return val
	})
}
