package array

import (
	"database/sql/driver"

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
	return "INSERT INTO array (tuple,runes,bytes,fixed_size,str) VALUES (?,?,?,?,?);", v.Values()
}
func (v Array) GetTuple() driver.Value {
	return string(v.Tuple[:])
}
func (v Array) GetRunes() driver.Value {
	return string(v.Runes[:])
}
func (v Array) GetBytes() driver.Value {
	return string(v.Bytes[:])
}
func (v Array) GetFixedSize() driver.Value {
	return string(v.FixedSize[:])
}
func (v Array) GetStr() driver.Value {
	return encoding.JSONValue(v.Str)
}
