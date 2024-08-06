package array

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Array) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "tuple", Definition: "tuple VARCHAR(255) NOT NULL"},
			{Name: "runes", Definition: "runes VARCHAR(255) NOT NULL"},
			{Name: "bytes", Definition: "bytes VARCHAR(255) NOT NULL"},
			{Name: "fixed_size", Definition: "fixed_size VARCHAR(255) NOT NULL"},
			{Name: "str", Definition: "str VARCHAR(255) NOT NULL"},
		},
	}
}
func (Array) TableName() string {
	return "array"
}
func (Array) Columns() []string {
	return []string{"tuple", "runes", "bytes", "fixed_size", "str"}
}
func (v Array) Values() []any {
	return []any{string(v.Tuple[:]), string(v.Runes[:]), string(v.Bytes[:]), string(v.FixedSize[:]), types.JSONMarshaler(v.Str)}
}
func (v *Array) Addrs() []any {
	return []any{types.FixedSizeBytes(v.Tuple[:], 2), types.FixedSizeRunes(v.Runes[:], 4), types.FixedSizeBytes(v.Bytes[:], 10), types.FixedSizeBytes(v.FixedSize[:], 10), types.JSONUnmarshaler(&v.Str)}
}
func (Array) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?)"
}
func (v Array) InsertOneStmt() (string, []any) {
	return "INSERT INTO array (tuple,runes,bytes,fixed_size,str) VALUES (?,?,?,?,?);", v.Values()
}
func (v Array) GetTuple() sequel.ColumnValuer[[2]byte] {
	return sequel.Column("tuple", v.Tuple, func(val [2]byte) driver.Value { return string(val[:]) })
}
func (v Array) GetRunes() sequel.ColumnValuer[[4]rune] {
	return sequel.Column("runes", v.Runes, func(val [4]rune) driver.Value { return string(val[:]) })
}
func (v Array) GetBytes() sequel.ColumnValuer[[10]byte] {
	return sequel.Column("bytes", v.Bytes, func(val [10]byte) driver.Value { return string(val[:]) })
}
func (v Array) GetFixedSize() sequel.ColumnValuer[[10]byte] {
	return sequel.Column("fixed_size", v.FixedSize, func(val [10]byte) driver.Value { return string(val[:]) })
}
func (v Array) GetStr() sequel.ColumnValuer[[100]Str] {
	return sequel.Column("str", v.Str, func(val [100]Str) driver.Value { return types.JSONMarshaler(val) })
}
