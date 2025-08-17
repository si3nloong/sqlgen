package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Array) TableName() string {
	return "array"
}
func (Array) Columns() []string {
	return []string{"tuple", "runes", "bytes", "fixed_size", "str"} // 5
}
func (v Array) Values() []any {
	return []any{
		string(v.Tuple[:]),        // 0 - tuple
		string(v.Runes[:]),        // 1 - runes
		string(v.Bytes[:]),        // 2 - bytes
		string(v.FixedSize[:]),    // 3 - fixed_size
		encoding.JSONValue(v.Str), // 4 - str
	}
}
func (v *Array) Addrs() []any {
	return []any{
		encoding.ByteArrayScanner(v.Tuple[:], 2),      // 0 - tuple
		encoding.RuneArrayScanner(v.Runes[:], 4),      // 1 - runes
		encoding.ByteArrayScanner(v.Bytes[:], 10),     // 2 - bytes
		encoding.ByteArrayScanner(v.FixedSize[:], 10), // 3 - fixed_size
		encoding.JSONScanner(&v.Str),                  // 4 - str
	}
}
func (Array) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?)" // 5
}
func (v Array) InsertOneStmt() (string, []any) {
	return "INSERT INTO `array` (`tuple`,`runes`,`bytes`,`fixed_size`,`str`) VALUES (?,?,?,?,?);", v.Values()
}
func (v Array) TupleValue() any {
	return string(v.Tuple[:])
}
func (v Array) RunesValue() any {
	return string(v.Runes[:])
}
func (v Array) BytesValue() any {
	return string(v.Bytes[:])
}
func (v Array) FixedSizeValue() any {
	return string(v.FixedSize[:])
}
func (v Array) StrValue() any {
	return encoding.JSONValue(v.Str)
}
func (v Array) ColumnTuple() sequel.ColumnValuer[[2]byte] {
	return sequel.Column("tuple", v.Tuple, func(val [2]byte) driver.Value {
		return string(val[:])
	})
}
func (v Array) ColumnRunes() sequel.ColumnValuer[[4]rune] {
	return sequel.Column("runes", v.Runes, func(val [4]rune) driver.Value {
		return string(val[:])
	})
}
func (v Array) ColumnBytes() sequel.ColumnValuer[[10]byte] {
	return sequel.Column("bytes", v.Bytes, func(val [10]byte) driver.Value {
		return string(val[:])
	})
}
func (v Array) ColumnFixedSize() sequel.ColumnValuer[[10]byte] {
	return sequel.Column("fixed_size", v.FixedSize, func(val [10]byte) driver.Value {
		return string(val[:])
	})
}
func (v Array) ColumnStr() sequel.ColumnValuer[[100]Str] {
	return sequel.Column("str", v.Str, func(val [100]Str) driver.Value {
		return encoding.JSONValue(val)
	})
}
