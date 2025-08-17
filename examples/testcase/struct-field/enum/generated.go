package enum

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Custom) TableName() string {
	return "custom"
}
func (Custom) Columns() []string {
	return []string{"text", "e", "num"} // 3
}
func (v Custom) Values() []any {
	return []any{
		(string)(v.Str), // 0 - text
		(int64)(v.Enum), // 1 - e
		(int64)(v.Num),  // 2 - num
	}
}
func (v *Custom) Addrs() []any {
	return []any{
		encoding.StringScanner[longText](&v.Str), // 0 - text
		encoding.IntScanner[Enum](&v.Enum),       // 1 - e
		encoding.Uint16Scanner[uint16](&v.Num),   // 2 - num
	}
}
func (Custom) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v Custom) InsertOneStmt() (string, []any) {
	return "INSERT INTO custom (text,e,num) VALUES (?,?,?);", v.Values()
}
func (v Custom) StrValue() driver.Value {
	return (string)(v.Str)
}
func (v Custom) EnumValue() driver.Value {
	return (int64)(v.Enum)
}
func (v Custom) NumValue() driver.Value {
	return (int64)(v.Num)
}
func (v Custom) ColumnStr() sequel.ColumnValuer[longText] {
	return sequel.Column("text", v.Str, func(val longText) driver.Value {
		return (string)(val)
	})
}
func (v Custom) ColumnEnum() sequel.ColumnValuer[Enum] {
	return sequel.Column("e", v.Enum, func(val Enum) driver.Value {
		return (int64)(val)
	})
}
func (v Custom) ColumnNum() sequel.ColumnValuer[uint16] {
	return sequel.Column("num", v.Num, func(val uint16) driver.Value {
		return (int64)(val)
	})
}
